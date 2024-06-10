package v2

import (
	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// +k8s:conversion-gen=false
// ImageVerification validates that images that match the specified pattern
// are signed with the supplied public key. Once the image is verified it is
// mutated to include the SHA digest retrieved during the registration.
type ImageVerification struct {
	// Type specifies the method of signature validation. The allowed options
	// are Cosign and Notary. By default Cosign is used if a type is not specified.
	// +kubebuilder:validation:Optional
	Type kyvernov1.ImageVerificationType `json:"type,omitempty" yaml:"type,omitempty"`

	// ImageReferences is a list of matching image reference patterns. At least one pattern in the
	// list must match the image for the rule to apply. Each image reference consists of a registry
	// address (defaults to docker.io), repository, image, and tag (defaults to latest).
	// Wildcards ('*' and '?') are allowed. See: https://kubernetes.io/docs/concepts/containers/images.
	// +kubebuilder:validation:Optional
	ImageReferences []string `json:"imageReferences,omitempty" yaml:"imageReferences,omitempty"`

	// SkipImageReferences is a list of matching image reference patterns that should be skipped.
	// At least one pattern in the list must match the image for the rule to be skipped. Each image reference
	// consists of a registry address (defaults to docker.io), repository, image, and tag (defaults to latest).
	// Wildcards ('*' and '?') are allowed. See: https://kubernetes.io/docs/concepts/containers/images.
	// +kubebuilder:validation:Optional
	SkipImageReferences []string `json:"skipImageReferences,omitempty" yaml:"skipImageReferences,omitempty"`

	// Attestors specified the required attestors (i.e. authorities)
	// +kubebuilder:validation:Optional
	Attestors []kyvernov1.AttestorSet `json:"attestors,omitempty" yaml:"attestors,omitempty"`

	// Attestations are optional checks for signed in-toto Statements used to verify the image.
	// See https://github.com/in-toto/attestation. Kyverno fetches signed attestations from the
	// OCI registry and decodes them into a list of Statement declarations.
	Attestations []Attestation `json:"attestations,omitempty" yaml:"attestations,omitempty"`

	// Repository is an optional alternate OCI repository to use for image signatures and attestations that match this rule.
	// If specified Repository will override the default OCI image repository configured for the installation.
	// The repository can also be overridden per Attestor or Attestation.
	Repository string `json:"repository,omitempty" yaml:"repository,omitempty"`

	// MutateDigest enables replacement of image tags with digests.
	// Defaults to true.
	// +kubebuilder:default=true
	// +kubebuilder:validation:Optional
	MutateDigest bool `json:"mutateDigest" yaml:"mutateDigest"`

	// VerifyDigest validates that images have a digest.
	// +kubebuilder:default=true
	// +kubebuilder:validation:Optional
	VerifyDigest bool `json:"verifyDigest" yaml:"verifyDigest"`

	// Required validates that images are verified i.e. have matched passed a signature or attestation check.
	// +kubebuilder:default=true
	// +kubebuilder:validation:Optional
	Required bool `json:"required" yaml:"required"`

	// ImageRegistryCredentials provides credentials that will be used for authentication with registry
	// +kubebuilder:validation:Optional
	ImageRegistryCredentials *kyvernov1.ImageRegistryCredentials `json:"imageRegistryCredentials,omitempty" yaml:"imageRegistryCredentials,omitempty"`

	// UseCache enables caching of image verify responses for this rule
	// +kubebuilder:default=true
	// +kubebuilder:validation:Optional
	UseCache bool `json:"useCache" yaml:"useCache"`
}

// Validate implements programmatic validation
func (iv *ImageVerification) Validate(isAuditFailureAction bool, path *field.Path) (errs field.ErrorList) {
	copy := iv

	if isAuditFailureAction && iv.MutateDigest {
		errs = append(errs, field.Invalid(path.Child("mutateDigest"), iv.MutateDigest, "mutateDigest must be set to false for ‘Audit’ failure action"))
	}

	if len(copy.ImageReferences) == 0 {
		errs = append(errs, field.Invalid(path, iv, "An image reference is required"))
	}

	asPath := path.Child("attestations")
	for i, attestation := range copy.Attestations {
		attestationErrors := attestation.Validate(asPath.Index(i))
		errs = append(errs, attestationErrors...)
	}

	attestorsPath := path.Child("attestors")
	for i, as := range copy.Attestors {
		attestorErrors := as.Validate(attestorsPath.Index(i))
		errs = append(errs, attestorErrors...)
	}

	if iv.Type == kyvernov1.Notary {
		for _, attestorSet := range iv.Attestors {
			for _, attestor := range attestorSet.Entries {
				if attestor.Keyless != nil {
					errs = append(errs, field.Invalid(attestorsPath, iv, "Keyless field is not allowed for type notary"))
				}
				if attestor.Keys != nil {
					errs = append(errs, field.Invalid(attestorsPath, iv, "Keys field is not allowed for type notary"))
				}
			}
		}
	}

	return errs
}

// +k8s:conversion-gen=false
// Attestation are checks for signed in-toto Statements that are used to verify the image.
// See https://github.com/in-toto/attestation. Kyverno fetches signed attestations from the
// OCI registry and decodes them into a list of Statements.
type Attestation struct {
	// Type defines the type of attestation contained within the Statement.
	// +kubebuilder:validation:Optional
	Type string `json:"type" yaml:"type"`

	// Attestors specify the required attestors (i.e. authorities).
	// +kubebuilder:validation:Optional
	Attestors []kyvernov1.AttestorSet `json:"attestors" yaml:"attestors"`

	// Conditions are used to verify attributes within a Predicate. If no Conditions are specified
	// the attestation check is satisfied as long there are predicates that match the predicate type.
	// +kubebuilder:validation:Optional
	Conditions []AnyAllConditions `json:"conditions,omitempty" yaml:"conditions,omitempty"`
}

func (a *Attestation) Validate(path *field.Path) (errs field.ErrorList) {
	if len(a.Attestors) == 0 {
		return
	}

	attestorsPath := path.Child("attestors")
	for i, as := range a.Attestors {
		attestorErrors := as.Validate(attestorsPath.Index(i))
		errs = append(errs, attestorErrors...)
	}
	return errs
}
