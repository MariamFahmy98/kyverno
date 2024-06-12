package main

import (
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"os"
	"sync"
	"time"

	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	kyvernov2 "github.com/kyverno/kyverno/api/kyverno/v2"
	"github.com/kyverno/kyverno/cmd/internal"
	"github.com/kyverno/kyverno/pkg/client/clientset/versioned/scheme"
	kyvernoinformer "github.com/kyverno/kyverno/pkg/client/informers/externalversions"
	"github.com/kyverno/kyverno/pkg/config"
	"github.com/kyverno/kyverno/pkg/controllers/certmanager"
	conversionwebhook "github.com/kyverno/kyverno/pkg/controllers/crd-conversion-webhook"
	"github.com/kyverno/kyverno/pkg/informers"
	"github.com/kyverno/kyverno/pkg/leaderelection"
	tlsutils "github.com/kyverno/kyverno/pkg/tls"
	kubeutils "github.com/kyverno/kyverno/pkg/utils/kube"
	corev1 "k8s.io/api/core/v1"
	apiserver "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	kubeinformers "k8s.io/client-go/informers"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

const (
	resyncPeriod = 15 * time.Minute
)

var (
	caSecretName  string
	tlsSecretName string
)

func sanityChecks(apiserverClient apiserver.Interface) error {
	return kubeutils.CRDsInstalled(apiserverClient, "policies.kyverno.io", "clusterpolicies.kyverno.io")
}

func main() {
	var (
		serverIP          string
		servicePort       int
		webhookServerPort int
		renewBefore       time.Duration
	)
	flagset := flag.NewFlagSet("crd-conversion-controller", flag.ExitOnError)
	flagset.StringVar(&serverIP, "serverIP", "", "IP address where Kyverno controller runs. Only required if out-of-cluster.")
	flagset.IntVar(&servicePort, "servicePort", 443, "Port used by the Kyverno Service resource and for webhook configurations.")
	flagset.IntVar(&webhookServerPort, "webhookServerPort", 9445, "Port used by the webhook server.")
	flagset.StringVar(&caSecretName, "caSecretName", "", "Name of the secret containing CA.")
	flagset.StringVar(&tlsSecretName, "tlsSecretName", "", "Name of the secret containing TLS pair.")
	flagset.DurationVar(&renewBefore, "renewBefore", 15*24*time.Hour, "The certificate renewal time before expiration")

	// config
	appConfig := internal.NewConfiguration(
		internal.WithProfiling(),
		internal.WithMetrics(),
		internal.WithTracing(),
		internal.WithKubeconfig(),
		internal.WithLeaderElection(),
		internal.WithKyvernoClient(),
		internal.WithKyvernoDynamicClient(),
		internal.WithEventsClient(),
		internal.WithConfigMapCaching(),
		internal.WithDeferredLoading(),
		internal.WithMetadataClient(),
		internal.WithApiServerClient(),
		internal.WithFlagSets(flagset),
	)
	// parse flags
	internal.ParseFlags(appConfig)
	// setup
	ctx, setup, sdown := internal.Setup(appConfig, "kyverno-crd-conversion-controller", false)
	defer sdown()
	if caSecretName == "" {
		setup.Logger.Error(errors.New("exiting... caSecretName is a required flag"), "exiting... caSecretName is a required flag")
		os.Exit(1)
	}
	if tlsSecretName == "" {
		setup.Logger.Error(errors.New("exiting... tlsSecretName is a required flag"), "exiting... tlsSecretName is a required flag")
		os.Exit(1)
	}
	if err := sanityChecks(setup.ApiServerClient); err != nil {
		setup.Logger.Error(err, "sanity checks failed")
		os.Exit(1)
	}
	// certificates informers
	caSecret := informers.NewSecretInformer(setup.KubeClient, config.KyvernoNamespace(), caSecretName, resyncPeriod)
	tlsSecret := informers.NewSecretInformer(setup.KubeClient, config.KyvernoNamespace(), tlsSecretName, resyncPeriod)
	if !informers.StartInformersAndWaitForCacheSync(ctx, setup.Logger, caSecret, tlsSecret) {
		setup.Logger.Error(errors.New("failed to wait for cache sync"), "failed to wait for cache sync")
		os.Exit(1)
	}
	// informer factories
	kubeInformer := kubeinformers.NewSharedInformerFactoryWithOptions(setup.KubeClient, resyncPeriod)
	kyvernoInformer := kyvernoinformer.NewSharedInformerFactory(setup.KyvernoClient, resyncPeriod)
	var wg sync.WaitGroup
	// start informers and wait for cache sync
	if !internal.StartInformersAndWaitForCacheSync(ctx, setup.Logger, kubeInformer, kyvernoInformer) {
		os.Exit(1)
	}
	// setup leader election
	le, err := leaderelection.New(
		setup.Logger.WithName("leader-election"),
		"kyverno-crd-conversion-controller",
		config.KyvernoNamespace(),
		setup.LeaderElectionClient,
		config.KyvernoPodName(),
		internal.LeaderElectionRetryPeriod(),
		func(ctx context.Context) {
			logger := setup.Logger.WithName("leader")
			// informer factories
			kubeInformer := kubeinformers.NewSharedInformerFactoryWithOptions(setup.KubeClient, resyncPeriod)
			kyvernoInformer := kyvernoinformer.NewSharedInformerFactory(setup.KyvernoClient, resyncPeriod)

			// controllers
			renewer := tlsutils.NewCertRenewer(
				setup.KubeClient.CoreV1().Secrets(config.KyvernoNamespace()),
				tlsutils.CertRenewalInterval,
				tlsutils.CAValidityDuration,
				tlsutils.TLSValidityDuration,
				renewBefore,
				serverIP,
				config.KyvernoServiceName(),
				config.DnsNames(config.KyvernoServiceName(), config.KyvernoNamespace()),
				config.KyvernoNamespace(),
				caSecretName,
				tlsSecretName,
			)
			certController := internal.NewController(
				certmanager.ControllerName,
				certmanager.NewController(
					caSecret,
					tlsSecret,
					renewer,
					caSecretName,
					tlsSecretName,
					config.KyvernoNamespace(),
				),
				certmanager.Workers,
			)
			crdConversionWebhookController := internal.NewController(
				conversionwebhook.ControllerName,
				conversionwebhook.NewController(
					setup.ApiServerClient,
					kubeInformer.Core().V1().Secrets(),
					caSecretName,
					serverIP,
					int32(servicePort),
					config.CrdConversionWebhookServicePath,
				),
				conversionwebhook.Workers,
			)
			// start informers and wait for cache sync
			if !internal.StartInformersAndWaitForCacheSync(ctx, logger, kyvernoInformer, kubeInformer) {
				logger.Error(errors.New("failed to wait for cache sync"), "failed to wait for cache sync")
				os.Exit(1)
			}
			// start leader controllers
			var wg sync.WaitGroup
			certController.Run(ctx, logger, &wg)
			crdConversionWebhookController.Run(ctx, logger, &wg)
			wg.Wait()
		},
		nil,
	)
	if err != nil {
		setup.Logger.Error(err, "failed to initialize leader election")
		os.Exit(1)
	}
	// create server
	var tlsOptions []func(config *tls.Config)
	tlsOptions = append(tlsOptions, func(cfg *tls.Config) {
		cfg.GetCertificate = func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
			secret, err := tlsSecret.Lister().Secrets(config.KyvernoNamespace()).Get(tlsSecretName)
			if err != nil {
				return nil, err
			}

			certPem := secret.Data[corev1.TLSCertKey]
			keyPem := secret.Data[corev1.TLSPrivateKeyKey]
			pair, err := tls.X509KeyPair(certPem, keyPem)
			if err != nil {
				return nil, err
			}
			return &pair, nil
		}
		cfg.MinVersion = tls.VersionTLS12
		cfg.CipherSuites =
			[]uint16{
				// AEADs w/ ECDHE
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			}
	})

	webhookServer := webhook.NewServer(webhook.Options{
		Port: webhookServerPort,
		// CertDir: "/tmp/k8s-webhook-server/serving-certs",
		TLSOpts: tlsOptions,
	})
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:        scheme.Scheme,
		WebhookServer: webhookServer,
	})
	if err != nil {
		setup.Logger.Error(err, "unable to start manager")
		os.Exit(1)
	}
	if err = (&kyvernov1.ClusterPolicy{}).SetupWebhookWithManager(mgr); err != nil {
		setup.Logger.Error(err, "unable to create webhook", "webhook", "ClusterPolicy")
		os.Exit(1)
	}
	if err = (&kyvernov1.Policy{}).SetupWebhookWithManager(mgr); err != nil {
		setup.Logger.Error(err, "unable to create webhook", "webhook", "Policy")
		os.Exit(1)
	}
	if err = (&kyvernov1.ClusterPolicyList{}).SetupWebhookWithManager(mgr); err != nil {
		setup.Logger.Error(err, "unable to create webhook", "webhook", "ClusterPolicy")
		os.Exit(1)
	}
	if err = (&kyvernov1.PolicyList{}).SetupWebhookWithManager(mgr); err != nil {
		setup.Logger.Error(err, "unable to create webhook", "webhook", "Policy")
		os.Exit(1)
	}
	if err = (&kyvernov2.ClusterPolicy{}).SetupWebhookWithManager(mgr); err != nil {
		setup.Logger.Error(err, "unable to create webhook", "webhook", "ClusterPolicy")
		os.Exit(1)
	}
	if err = (&kyvernov2.Policy{}).SetupWebhookWithManager(mgr); err != nil {
		setup.Logger.Error(err, "unable to create webhook", "webhook", "Policy")
		os.Exit(1)
	}
	if err = (&kyvernov2.ClusterPolicyList{}).SetupWebhookWithManager(mgr); err != nil {
		setup.Logger.Error(err, "unable to create webhook", "webhook", "ClusterPolicy")
		os.Exit(1)
	}
	if err = (&kyvernov2.PolicyList{}).SetupWebhookWithManager(mgr); err != nil {
		setup.Logger.Error(err, "unable to create webhook", "webhook", "Policy")
		os.Exit(1)
	}
	setup.Logger.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setup.Logger.Error(err, "problem running manager")
		os.Exit(1)
	}
	// start leader election
	le.Run(ctx)
	// wait for everything to shut down and exit
	wg.Wait()
}
