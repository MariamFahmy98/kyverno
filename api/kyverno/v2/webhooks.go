package v2

import (
	ctrl "sigs.k8s.io/controller-runtime"
)

// SetupWebhookWithManager sets up ClusterPolicy conversion webhook.
func (r *ClusterPolicy) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// SetupWebhookWithManager sets up ClusterPolicyList conversion webhook.
func (r *ClusterPolicyList) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// SetupWebhookWithManager sets up Policy conversion webhook.
func (r *Policy) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// SetupWebhookWithManager sets up PolicyList conversion webhook.
func (r *PolicyList) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}
