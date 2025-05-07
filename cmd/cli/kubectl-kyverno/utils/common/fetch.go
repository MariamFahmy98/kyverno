package common

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/go-git/go-billy/v5"
	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	"github.com/kyverno/kyverno/cmd/cli/kubectl-kyverno/apis/v1alpha1"
	"github.com/kyverno/kyverno/cmd/cli/kubectl-kyverno/log"
	"github.com/kyverno/kyverno/cmd/cli/kubectl-kyverno/resource"
	"github.com/kyverno/kyverno/pkg/admissionpolicy"
	"github.com/kyverno/kyverno/pkg/autogen"
	"github.com/kyverno/kyverno/pkg/clients/dclient"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	kubeutils "github.com/kyverno/kyverno/pkg/utils/kube"
	utils "github.com/kyverno/kyverno/pkg/utils/restmapper"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// GetResources gets matched resources by the given policies
// The resources are fetched from:
// - local paths, if given
// - the k8s cluster, if given
func GetResources(
	out io.Writer,
	policies []engineapi.GenericPolicy,
	resourcePaths []string,
	dClient dclient.Interface,
	cluster bool,
	namespace string,
	policyReport bool,
	clusterWideResources bool,
) ([]*unstructured.Unstructured, error) {
	resources := make([]*unstructured.Unstructured, 0)
	var err error

	if cluster && dClient != nil {
		resources, err = fetchResourcesFromCluster(out, policies, resourcePaths, dClient, namespace, policyReport, clusterWideResources)
		if err != nil {
			return resources, err
		}
	} else if len(resourcePaths) > 0 {
		resources, err = fetchResourcesFromLocalFiles(out, resourcePaths, policyReport)
		if err != nil {
			return resources, err
		}
	}
	return resources, err
}

func fetchResourcesFromLocalFiles(out io.Writer, resourcePaths []string, policyReport bool) ([]*unstructured.Unstructured, error) {
	resources := make([]*unstructured.Unstructured, 0)
	for _, resourcePath := range resourcePaths {
		resourceBytes, err := resource.GetFileBytes(resourcePath)
		if err != nil {
			if policyReport {
				log.Log.V(3).Info(fmt.Sprintf("failed to load resources: %s.", resourcePath), "error", err)
			} else {
				fmt.Fprintf(out, "\n----------------------------------------------------------------------\nfailed to load resources: %s. \nerror: %s\n----------------------------------------------------------------------\n", resourcePath, err)
			}
			continue
		}

		getResources, err := resource.GetUnstructuredResources(resourceBytes)
		if err != nil {
			return nil, err
		}

		resources = append(resources, getResources...)
	}
	return resources, nil
}

type ResourceTypeInfo struct {
	gvkMap         map[schema.GroupVersionKind]bool
	subresourceMap map[schema.GroupVersionKind]v1alpha1.Subresource
}

// fetchResourcesFromCluster fetches resources from the cluster.
// It will first extract the matched resources from the policies.
// Then it will fetch the resources from the cluster.
func fetchResourcesFromCluster(
	out io.Writer,
	policies []engineapi.GenericPolicy,
	resourcePaths []string,
	dClient dclient.Interface,
	namespace string,
	policyReport bool,
	clusterWideResources bool,
) ([]*unstructured.Unstructured, error) {
	resources := make([]*unstructured.Unstructured, 0)
	info := &ResourceTypeInfo{
		gvkMap:         make(map[schema.GroupVersionKind]bool),
		subresourceMap: make(map[schema.GroupVersionKind]v1alpha1.Subresource),
	}
	// extract the matched resources from the policies.
	extractResourcesFromPolicies(policies, info, dClient, clusterWideResources)
	// fetch the resources from the cluster.
	resourceMap, err := listResources(info, dClient, namespace)
	if err != nil {
		return nil, err
	}
	if len(resourcePaths) == 0 {
		for _, rr := range resourceMap {
			resources = append(resources, rr)
		}
	} else {
		for _, resourcePath := range resourcePaths {
			lenOfResource := len(resources)
			for rn, rr := range resourceMap {
				s := strings.Split(rn, "-")
				if s[2] == resourcePath {
					resources = append(resources, rr)
				}
			}
			if lenOfResource >= len(resources) {
				if policyReport {
					log.Log.V(3).Info(fmt.Sprintf("%s not found in cluster", resourcePath))
				} else {
					fmt.Fprintf(out, "\n----------------------------------------------------------------------\nresource %s not found in cluster\n----------------------------------------------------------------------\n", resourcePath)
				}
				return nil, fmt.Errorf("%s not found in cluster", resourcePath)
			}
		}
	}
	return resources, nil
}

func extractResourcesFromPolicies(
	policies []engineapi.GenericPolicy,
	info *ResourceTypeInfo,
	dClient dclient.Interface,
	clusterWideResources bool,
) {
	for _, policy := range policies {
		if kpol := policy.AsKyvernoPolicy(); kpol != nil {
			for _, rule := range autogen.Default.ComputeRules(kpol, "") {
				getKindsFromRule(rule, info, dClient, clusterWideResources)
			}
		} else if vap := policy.AsValidatingAdmissionPolicy(); vap != nil {
			getKindsFromValidatingAdmissionPolicy(*vap.GetDefinition(), info, dClient, clusterWideResources)
		}
	}
}

// getKindsFromRule will return the kinds from policy match block
func getKindsFromRule(
	rule kyvernov1.Rule,
	info *ResourceTypeInfo,
	client dclient.Interface,
	clusterWideResources bool,
) {
	for _, kind := range rule.MatchResources.Kinds {
		addToResourceTypeInfo(kind, info, client, clusterWideResources)
	}
	if rule.MatchResources.Any != nil {
		for _, resFilter := range rule.MatchResources.Any {
			for _, kind := range resFilter.ResourceDescription.Kinds {
				addToResourceTypeInfo(kind, info, client, clusterWideResources)
			}
		}
	}
	if rule.MatchResources.All != nil {
		for _, resFilter := range rule.MatchResources.All {
			for _, kind := range resFilter.ResourceDescription.Kinds {
				addToResourceTypeInfo(kind, info, client, clusterWideResources)
			}
		}
	}
}

// getKindsFromValidatingAdmissionPolicy will return the kinds from policy match block
func getKindsFromValidatingAdmissionPolicy(
	policy admissionregistrationv1.ValidatingAdmissionPolicy,
	info *ResourceTypeInfo,
	client dclient.Interface,
	clusterWideResources bool,
) {
	restMapper, err := utils.GetRESTMapper(client, false)
	if err != nil {
		log.Log.V(3).Info("failed to get rest mapper", "error", err)
		return
	}
	kinds, err := admissionpolicy.GetKinds(policy.Spec.MatchConstraints, restMapper)
	if err != nil {
		log.Log.V(3).Info("failed to get kinds from validating admission policy", "error", err)
		return
	}
	for _, kind := range kinds {
		addToResourceTypeInfo(kind, info, client, clusterWideResources)
	}
}

func addToResourceTypeInfo(
	kind string,
	info *ResourceTypeInfo,
	client dclient.Interface,
	clusterWideResources bool,
) {
	group, version, kind, subresource := kubeutils.ParseKindSelector(kind)
	resourceDefs, err := client.Discovery().FindResources(group, version, kind, subresource)
	if err != nil {
		log.Log.V(2).Info("Failed to find resource", "kind", kind, "error", err)
		return
	}

	for parent, child := range resourceDefs {
		if clusterWideResources && child.Namespaced {
			continue
		}

		if parent.SubResource == "" {
			info.gvkMap[parent.GroupVersionKind()] = true
		} else {
			subGVK := schema.GroupVersionKind{
				Group:   child.Group,
				Version: child.Version,
				Kind:    child.Kind,
			}
			info.subresourceMap[subGVK] = v1alpha1.Subresource{
				Subresource: child,
				ParentResource: metav1.APIResource{
					Group:   parent.Group,
					Version: parent.Version,
					Kind:    parent.Kind,
					Name:    parent.Resource,
				},
			}
		}
	}
}

func listResources(
	info *ResourceTypeInfo,
	dClient dclient.Interface,
	namespace string,
) (map[string]*unstructured.Unstructured, error) {
	result := make(map[string]*unstructured.Unstructured)

	// list standard resources
	for gvk := range info.gvkMap {
		resourceList, err := dClient.ListResource(
			context.TODO(),
			gvk.GroupVersion().String(),
			gvk.Kind,
			namespace,
			nil,
		)
		if err != nil {
			log.Log.V(3).Info("failed to list resource", "gvk", gvk, "error", err)
			continue
		}
		for _, resource := range resourceList.Items {
			key := fmt.Sprintf("%s-%s-%s", gvk.Kind, resource.GetNamespace(), resource.GetName())
			resource.SetGroupVersionKind(gvk)
			result[key] = resource.DeepCopy()
		}
	}

	// list subresources
	for subGVK, subresource := range info.subresourceMap {
		parentGV := schema.GroupVersion{
			Group:   subresource.ParentResource.Group,
			Version: subresource.ParentResource.Version,
		}
		resourceList, err := dClient.ListResource(
			context.TODO(),
			parentGV.String(),
			subresource.ParentResource.Kind,
			namespace,
			nil,
		)
		if err != nil {
			log.Log.V(3).Info("failed to list parent resource", "gv", parentGV, "kind", subresource.ParentResource.Kind, "error", err)
			continue
		}

		for _, parent := range resourceList.Items {
			subresourceName := strings.Split(subresource.Subresource.Name, "/")[1]
			resource, err := dClient.GetResource(
				context.TODO(),
				parentGV.String(),
				subresource.ParentResource.Kind,
				namespace,
				parent.GetName(),
				subresourceName,
			)
			if err != nil {
				log.Log.V(3).Info("failed to get subresource", "parent", parent.GetName(), "subresource", subresourceName, "error", err)
				continue
			}

			key := fmt.Sprintf("%s-%s-%s", subGVK.Kind, resource.GetNamespace(), resource.GetName())
			resource.SetGroupVersionKind(subGVK)
			result[key] = resource.DeepCopy()
		}
	}
	return result, nil
}

// GetResourcesWithTest with gets matched resources by the given policies
func GetResourcesWithTest(out io.Writer, fs billy.Filesystem, resourcePaths []string, policyResourcePath string) ([]*unstructured.Unstructured, error) {
	resources := make([]*unstructured.Unstructured, 0)
	if len(resourcePaths) > 0 {
		for _, resourcePath := range resourcePaths {
			var resourceBytes []byte
			var err error
			if fs != nil {
				filep, err := fs.Open(filepath.Join(policyResourcePath, resourcePath))
				if err != nil {
					fmt.Fprintf(out, "Unable to open resource file: %s. error: %s", resourcePath, err)
					continue
				}
				resourceBytes, _ = io.ReadAll(filep)
			} else {
				resourceBytes, err = resource.GetFileBytes(resourcePath)
			}
			if err != nil {
				fmt.Fprintf(out, "\n----------------------------------------------------------------------\nfailed to load resources: %s. \nerror: %s\n----------------------------------------------------------------------\n", resourcePath, err)
				continue
			}

			getResources, err := resource.GetUnstructuredResources(resourceBytes)
			if err != nil {
				return nil, err
			}

			resources = append(resources, getResources...)
		}
	}
	return resources, nil
}
