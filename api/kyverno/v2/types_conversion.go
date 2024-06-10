package v2

import (
	"encoding/json"
	unsafe "unsafe"

	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	conversionutils "github.com/kyverno/kyverno/pkg/utils/conversion"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiconversion "k8s.io/apimachinery/pkg/conversion"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

func (src *ClusterPolicy) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*kyvernov1.ClusterPolicy)

	if err := Convert_v2_ClusterPolicy_To_v1_ClusterPolicy(src, dst, nil); err != nil {
		return err
	}

	// Manually restore data.
	restored := &kyvernov1.ClusterPolicy{}
	if ok, err := conversionutils.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}

	dst.Spec.GenerateExistingOnPolicyUpdate = restored.Spec.GenerateExistingOnPolicyUpdate
	dst.Spec.SchemaValidation = restored.Spec.SchemaValidation

	for r, rule := range restored.Spec.Rules {
		for i, image := range rule.VerifyImages {
			if image.Image != "" {
				dst.Spec.Rules[r].VerifyImages[i].Image = image.Image
			}
			if image.Key != "" {
				dst.Spec.Rules[r].VerifyImages[i].Key = image.Key
			}
			if image.Roots != "" {
				dst.Spec.Rules[r].VerifyImages[i].Roots = image.Roots
			}
			if image.Subject != "" {
				dst.Spec.Rules[r].VerifyImages[i].Subject = image.Subject
			}
			if image.Issuer != "" {
				dst.Spec.Rules[r].VerifyImages[i].Issuer = image.Issuer
			}
			if len(image.AdditionalExtensions) != 0 {
				dst.Spec.Rules[r].VerifyImages[i].AdditionalExtensions = image.AdditionalExtensions
			}
			if len(image.Annotations) != 0 {
				dst.Spec.Rules[r].VerifyImages[i].Annotations = image.Annotations
			}

			for a, attestation := range image.Attestations {
				if attestation.PredicateType != "" {
					dst.Spec.Rules[r].VerifyImages[i].Attestations[a].PredicateType = attestation.PredicateType
				}
			}
		}
	}

	return nil
}

func (dst *ClusterPolicy) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*kyvernov1.ClusterPolicy)

	if err := Convert_v1_ClusterPolicy_To_v2_ClusterPolicy(src, dst, nil); err != nil {
		return err
	}

	// Preserve Hub data on down-conversion except for metadata
	if err := conversionutils.MarshalData(src, dst); err != nil {
		return err
	}

	return nil
}

func (src *ClusterPolicyList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*kyvernov1.ClusterPolicyList)

	return Convert_v2_ClusterPolicyList_To_v1_ClusterPolicyList(src, dst, nil)
}

func (dst *ClusterPolicyList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*kyvernov1.ClusterPolicyList)

	return Convert_v1_ClusterPolicyList_To_v2_ClusterPolicyList(src, dst, nil)
}

func (src *Policy) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*kyvernov1.Policy)

	return Convert_v2_Policy_To_v1_Policy(src, dst, nil)
}

func (dst *Policy) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*kyvernov1.Policy)

	return Convert_v1_Policy_To_v2_Policy(src, dst, nil)
}

func (src *PolicyList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*kyvernov1.PolicyList)

	return Convert_v2_PolicyList_To_v1_PolicyList(src, dst, nil)
}

func (dst *PolicyList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*kyvernov1.PolicyList)

	return Convert_v1_PolicyList_To_v2_PolicyList(src, dst, nil)
}

func autoConvert_v2_Spec_To_v1_Spec(in *Spec, out *kyvernov1.Spec, s apiconversion.Scope) error {
	for _, rule := range in.Rules {
		var v1Rule kyvernov1.Rule
		if err := Convert_v2_Rule_To_v1_Rule(&rule, &v1Rule, s); err != nil {
			return err
		}
		out.Rules = append(out.Rules, v1Rule)
		if rule.HasValidate() {
			out.ValidationFailureAction = rule.Validation.ValidationFailureAction
		} else if rule.HasMutate() {
			out.MutateExistingOnPolicyUpdate = rule.Mutation.MutateExistingOnPolicyUpdate
		} else if rule.HasGenerate() {
			out.GenerateExisting = rule.Generation.GenerateExisting
		}
	}
	out.ApplyRules = (*kyvernov1.ApplyRulesType)(unsafe.Pointer(in.ApplyRules))
	out.FailurePolicy = (*kyvernov1.FailurePolicyType)(unsafe.Pointer(in.FailurePolicy))
	out.Admission = (*bool)(unsafe.Pointer(in.Admission))
	out.Background = (*bool)(unsafe.Pointer(in.Background))
	out.WebhookTimeoutSeconds = (*int32)(unsafe.Pointer(in.WebhookTimeoutSeconds))
	out.UseServerSideApply = in.UseServerSideApply
	out.WebhookConfiguration = (*kyvernov1.WebhookConfiguration)(unsafe.Pointer(in.WebhookConfiguration))
	return nil
}

func Convert_v2_Spec_To_v1_Spec(in *Spec, out *kyvernov1.Spec, s apiconversion.Scope) error {
	return autoConvert_v2_Spec_To_v1_Spec(in, out, s)
}

func autoConvert_v1_Spec_To_v2_Spec(in *kyvernov1.Spec, out *Spec, s apiconversion.Scope) error {
	for _, rule := range in.Rules {
		var v2Rule Rule
		if err := Convert_v1_Rule_To_v2_Rule(&rule, &v2Rule, s); err != nil {
			return err
		}

		if rule.HasValidate() {
			v2Rule.Validation.ValidationFailureAction = in.ValidationFailureAction
		} else if rule.HasMutate() {
			v2Rule.Mutation.MutateExistingOnPolicyUpdate = in.MutateExistingOnPolicyUpdate
		} else if rule.HasGenerate() {
			// TODO: what if the user specified the GenerateExistingOnPolicyUpdate field instead?
			v2Rule.Generation.GenerateExisting = in.GenerateExisting
		}
		out.Rules = append(out.Rules, v2Rule)
	}
	out.ApplyRules = (*kyvernov1.ApplyRulesType)(unsafe.Pointer(in.ApplyRules))
	out.FailurePolicy = (*kyvernov1.FailurePolicyType)(unsafe.Pointer(in.FailurePolicy))
	out.Admission = (*bool)(unsafe.Pointer(in.Admission))
	out.Background = (*bool)(unsafe.Pointer(in.Background))
	out.WebhookTimeoutSeconds = (*int32)(unsafe.Pointer(in.WebhookTimeoutSeconds))
	out.UseServerSideApply = in.UseServerSideApply
	out.WebhookConfiguration = (*WebhookConfiguration)(unsafe.Pointer(in.WebhookConfiguration))
	return nil
}

func Convert_v1_Spec_To_v2_Spec(in *kyvernov1.Spec, out *Spec, s apiconversion.Scope) error {
	return autoConvert_v1_Spec_To_v2_Spec(in, out, s)
}

func autoConvert_v1_ImageVerification_To_v2_ImageVerification(in *kyvernov1.ImageVerification, out *ImageVerification, s apiconversion.Scope) error {
	out.Type = kyvernov1.ImageVerificationType(in.Type)
	out.ImageReferences = *(*[]string)(unsafe.Pointer(&in.ImageReferences))
	out.ImageReferences = append(out.ImageReferences, in.Image)
	out.SkipImageReferences = *(*[]string)(unsafe.Pointer(&in.SkipImageReferences))
	// WARNING: in.Key requires manual conversion: does not exist in peer-type
	// WARNING: in.Roots requires manual conversion: does not exist in peer-type
	// WARNING: in.Subject requires manual conversion: does not exist in peer-type
	// WARNING: in.Issuer requires manual conversion: does not exist in peer-type
	// WARNING: in.AdditionalExtensions requires manual conversion: does not exist in peer-type
	out.Attestors = *(*[]kyvernov1.AttestorSet)(unsafe.Pointer(&in.Attestors))
	for _, attestation := range in.Attestations {
		var v2Attestation Attestation
		if err := Convert_v1_Attestation_To_v2_Attestation(&attestation, &v2Attestation, s); err != nil {
			return err
		}
		for _, attestor := range v2Attestation.Attestors {
			for _, entry := range attestor.Entries {
				if in.Key != "" {
					entry.Keys = &kyvernov1.StaticKeyAttestor{
						PublicKeys: in.Key,
					}
				}
				var keyless *kyvernov1.KeylessAttestor
				if in.Roots != "" {
					keyless.Roots = in.Roots
				}
				if in.Subject != "" {
					keyless.Subject = in.Subject
				}
				if in.Issuer != "" {
					keyless.Issuer = in.Issuer
				}
				if len(in.AdditionalExtensions) != 0 {
					keyless.AdditionalExtensions = in.AdditionalExtensions
				}
				entry.Keyless = keyless
			}
		}
		out.Attestations = append(out.Attestations, v2Attestation)
	}
	// WARNING: in.Annotations requires manual conversion: does not exist in peer-type
	out.Repository = in.Repository
	out.MutateDigest = in.MutateDigest
	out.VerifyDigest = in.VerifyDigest
	out.Required = in.Required
	out.ImageRegistryCredentials = (*kyvernov1.ImageRegistryCredentials)(unsafe.Pointer(in.ImageRegistryCredentials))
	out.UseCache = in.UseCache
	return nil
}

func Convert_v1_ImageVerification_To_v2_ImageVerification(in *kyvernov1.ImageVerification, out *ImageVerification, s apiconversion.Scope) error {
	return autoConvert_v1_ImageVerification_To_v2_ImageVerification(in, out, s)
}

func autoConvert_v2_ImageVerification_To_v1_ImageVerification(in *ImageVerification, out *kyvernov1.ImageVerification, s apiconversion.Scope) error {
	out.Type = kyvernov1.ImageVerificationType(in.Type)
	out.ImageReferences = *(*[]string)(unsafe.Pointer(&in.ImageReferences))
	out.SkipImageReferences = *(*[]string)(unsafe.Pointer(&in.SkipImageReferences))
	out.Attestors = *(*[]kyvernov1.AttestorSet)(unsafe.Pointer(&in.Attestors))
	if in.Attestations != nil {
		in, out := &in.Attestations, &out.Attestations
		*out = make([]kyvernov1.Attestation, len(*in))
		for i := range *in {
			if err := Convert_v2_Attestation_To_v1_Attestation(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Attestations = nil
	}
	out.Repository = in.Repository
	out.MutateDigest = in.MutateDigest
	out.VerifyDigest = in.VerifyDigest
	out.Required = in.Required
	out.ImageRegistryCredentials = (*kyvernov1.ImageRegistryCredentials)(unsafe.Pointer(in.ImageRegistryCredentials))
	out.UseCache = in.UseCache
	return nil
}

func Convert_v2_ImageVerification_To_v1_ImageVerification(in *ImageVerification, out *kyvernov1.ImageVerification, s apiconversion.Scope) error {
	return autoConvert_v2_ImageVerification_To_v1_ImageVerification(in, out, s)
}

func autoConvert_v1_MatchResources_To_v2_MatchResources(in *kyvernov1.MatchResources, out *MatchResources, s apiconversion.Scope) error {
	if in.Any != nil {
		in, out := &in.Any, &out.Any
		*out = make(ResourceFilters, len(*in))
		for i := range *in {
			if err := Convert_v1_ResourceFilter_To_v2_ResourceFilter(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	}

	if in.All != nil {
		in, out := &in.All, &out.All
		*out = make(ResourceFilters, len(*in))
		for i := range *in {
			if err := Convert_v1_ResourceFilter_To_v2_ResourceFilter(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	}

	if !in.UserInfo.IsEmpty() {
		out.Any = append(out.Any, ResourceFilter{
			UserInfo: in.UserInfo,
		})
	}
	if !in.ResourceDescription.IsEmpty() {
		var resource ResourceDescription
		if err := Convert_v1_ResourceDescription_To_v2_ResourceDescription(&in.ResourceDescription, &resource, s); err != nil {
			return err
		}
		out.Any = append(out.Any, ResourceFilter{
			ResourceDescription: resource,
		})
	}
	return nil
}

func Convert_v1_MatchResources_To_v2_MatchResources(in *kyvernov1.MatchResources, out *MatchResources, s apiconversion.Scope) error {
	return autoConvert_v1_MatchResources_To_v2_MatchResources(in, out, s)
}

func autoConvert_v2_MatchResources_To_v1_MatchResources(in *MatchResources, out *kyvernov1.MatchResources, s apiconversion.Scope) error {
	if in.Any != nil {
		in, out := &in.Any, &out.Any
		*out = make(kyvernov1.ResourceFilters, len(*in))
		for i := range *in {
			if err := Convert_v2_ResourceFilter_To_v1_ResourceFilter(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Any = nil
	}
	if in.All != nil {
		in, out := &in.All, &out.All
		*out = make(kyvernov1.ResourceFilters, len(*in))
		for i := range *in {
			if err := Convert_v2_ResourceFilter_To_v1_ResourceFilter(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.All = nil
	}
	return nil
}

// Convert_v2_MatchResources_To_v1_MatchResources is an autogenerated conversion function.
func Convert_v2_MatchResources_To_v1_MatchResources(in *MatchResources, out *kyvernov1.MatchResources, s apiconversion.Scope) error {
	return autoConvert_v2_MatchResources_To_v1_MatchResources(in, out, s)
}

func autoConvert_v1_ResourceDescription_To_v2_ResourceDescription(in *kyvernov1.ResourceDescription, out *ResourceDescription, _ apiconversion.Scope) error {
	out.Kinds = *(*[]string)(unsafe.Pointer(&in.Kinds))
	out.Names = *(*[]string)(unsafe.Pointer(&in.Names))
	if in.Name != "" {
		out.Names = append(out.Names, in.Name)
	}
	out.Namespaces = *(*[]string)(unsafe.Pointer(&in.Namespaces))
	out.Annotations = *(*map[string]string)(unsafe.Pointer(&in.Annotations))
	out.Selector = (*metav1.LabelSelector)(unsafe.Pointer(in.Selector))
	out.NamespaceSelector = (*metav1.LabelSelector)(unsafe.Pointer(in.NamespaceSelector))
	out.Operations = *(*[]kyvernov1.AdmissionOperation)(unsafe.Pointer(&in.Operations))
	return nil
}

func Convert_v1_ResourceDescription_To_v2_ResourceDescription(in *kyvernov1.ResourceDescription, out *ResourceDescription, s apiconversion.Scope) error {
	return autoConvert_v1_ResourceDescription_To_v2_ResourceDescription(in, out, s)
}

func autoConvert_v2_ResourceDescription_To_v1_ResourceDescription(in *ResourceDescription, out *kyvernov1.ResourceDescription, _ apiconversion.Scope) error {
	out.Kinds = *(*[]string)(unsafe.Pointer(&in.Kinds))
	out.Names = *(*[]string)(unsafe.Pointer(&in.Names))
	out.Namespaces = *(*[]string)(unsafe.Pointer(&in.Namespaces))
	out.Annotations = *(*map[string]string)(unsafe.Pointer(&in.Annotations))
	out.Selector = (*metav1.LabelSelector)(unsafe.Pointer(in.Selector))
	out.NamespaceSelector = (*metav1.LabelSelector)(unsafe.Pointer(in.NamespaceSelector))
	out.Operations = *(*[]kyvernov1.AdmissionOperation)(unsafe.Pointer(&in.Operations))
	return nil
}

// Convert_v2_ResourceDescription_To_v1_ResourceDescription is an autogenerated conversion function.
func Convert_v2_ResourceDescription_To_v1_ResourceDescription(in *ResourceDescription, out *kyvernov1.ResourceDescription, s apiconversion.Scope) error {
	return autoConvert_v2_ResourceDescription_To_v1_ResourceDescription(in, out, s)
}

func autoConvert_v1_Validation_To_v2_Validation(in *kyvernov1.Validation, out *Validation, s apiconversion.Scope) error {
	out.Message = in.Message
	out.Manifests = (*kyvernov1.Manifests)(unsafe.Pointer(in.Manifests))
	out.ForEachValidation = *(*[]kyvernov1.ForEachValidation)(unsafe.Pointer(&in.ForEachValidation))
	out.RawPattern = (*apiextensionsv1.JSON)(unsafe.Pointer(in.RawPattern))
	out.RawAnyPattern = (*apiextensionsv1.JSON)(unsafe.Pointer(in.RawAnyPattern))
	if in.Deny != nil {
		in, out := &in.Deny, &out.Deny
		*out = new(Deny)
		if err := Convert_v1_Deny_To_v2_Deny(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.Deny = nil
	}
	out.PodSecurity = (*kyvernov1.PodSecurity)(unsafe.Pointer(in.PodSecurity))
	out.CEL = (*kyvernov1.CEL)(unsafe.Pointer(in.CEL))
	return nil
}

// Convert_v1_Validation_To_v2_Validation is an autogenerated conversion function.
func Convert_v1_Validation_To_v2_Validation(in *kyvernov1.Validation, out *Validation, s apiconversion.Scope) error {
	return autoConvert_v1_Validation_To_v2_Validation(in, out, s)
}

func autoConvert_v2_Validation_To_v1_Validation(in *Validation, out *kyvernov1.Validation, s apiconversion.Scope) error {
	// WARNING: in.ValidationFailureAction requires manual conversion: does not exist in peer-type
	// WARNING: in.ValidationFailureActionOverrides requires manual conversion: does not exist in peer-type
	out.Message = in.Message
	out.Manifests = (*kyvernov1.Manifests)(unsafe.Pointer(in.Manifests))
	out.ForEachValidation = *(*[]kyvernov1.ForEachValidation)(unsafe.Pointer(&in.ForEachValidation))
	out.RawPattern = (*apiextensionsv1.JSON)(unsafe.Pointer(in.RawPattern))
	out.RawAnyPattern = (*apiextensionsv1.JSON)(unsafe.Pointer(in.RawAnyPattern))
	if in.Deny != nil {
		in, out := &in.Deny, &out.Deny
		*out = new(kyvernov1.Deny)
		if err := Convert_v2_Deny_To_v1_Deny(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.Deny = nil
	}
	out.PodSecurity = (*kyvernov1.PodSecurity)(unsafe.Pointer(in.PodSecurity))
	out.CEL = (*kyvernov1.CEL)(unsafe.Pointer(in.CEL))
	return nil
}

func Convert_v2_Validation_To_v1_Validation(in *Validation, out *kyvernov1.Validation, s apiconversion.Scope) error {
	return autoConvert_v2_Validation_To_v1_Validation(in, out, s)

}

func autoConvert_v1_Mutation_To_v2_Mutation(in *kyvernov1.Mutation, out *Mutation, s apiconversion.Scope) error {
	if in.Targets != nil {
		in, out := &in.Targets, &out.Targets
		*out = make([]TargetResourceSpec, len(*in))
		for i := range *in {
			if err := Convert_v1_TargetResourceSpec_To_v2_TargetResourceSpec(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Targets = nil
	}
	out.RawPatchStrategicMerge = (*apiextensionsv1.JSON)(unsafe.Pointer(in.RawPatchStrategicMerge))
	out.PatchesJSON6902 = in.PatchesJSON6902
	out.ForEachMutation = *(*[]kyvernov1.ForEachMutation)(unsafe.Pointer(&in.ForEachMutation))
	return nil
}

func Convert_v1_Mutation_To_v2_Mutation(in *kyvernov1.Mutation, out *Mutation, s apiconversion.Scope) error {
	return autoConvert_v1_Mutation_To_v2_Mutation(in, out, s)
}

func autoConvert_v2_Mutation_To_v1_Mutation(in *Mutation, out *kyvernov1.Mutation, s apiconversion.Scope) error {
	// WARNING: in.MutateExistingOnPolicyUpdate requires manual conversion: does not exist in peer-type
	if in.Targets != nil {
		in, out := &in.Targets, &out.Targets
		*out = make([]kyvernov1.TargetResourceSpec, len(*in))
		for i := range *in {
			if err := Convert_v2_TargetResourceSpec_To_v1_TargetResourceSpec(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Targets = nil
	}
	out.RawPatchStrategicMerge = (*apiextensionsv1.JSON)(unsafe.Pointer(in.RawPatchStrategicMerge))
	out.PatchesJSON6902 = in.PatchesJSON6902
	out.ForEachMutation = *(*[]kyvernov1.ForEachMutation)(unsafe.Pointer(&in.ForEachMutation))
	return nil
}

func Convert_v2_Mutation_To_v1_Mutation(in *Mutation, out *kyvernov1.Mutation, s apiconversion.Scope) error {
	return autoConvert_v2_Mutation_To_v1_Mutation(in, out, s)
}

func autoConvert_v1_Generation_To_v2_Generation(in *kyvernov1.Generation, out *Generation, _ apiconversion.Scope) error {
	out.ResourceSpec = in.ResourceSpec
	out.Synchronize = in.Synchronize
	out.OrphanDownstreamOnPolicyDelete = in.OrphanDownstreamOnPolicyDelete
	out.RawData = (*apiextensionsv1.JSON)(unsafe.Pointer(in.RawData))
	out.Clone = in.Clone
	out.CloneList = in.CloneList
	return nil
}

// Convert_v1_Generation_To_v2_Generation is an autogenerated conversion function.
func Convert_v1_Generation_To_v2_Generation(in *kyvernov1.Generation, out *Generation, s apiconversion.Scope) error {
	return autoConvert_v1_Generation_To_v2_Generation(in, out, s)
}

func autoConvert_v2_Generation_To_v1_Generation(in *Generation, out *kyvernov1.Generation, _ apiconversion.Scope) error {
	out.ResourceSpec = in.ResourceSpec
	// WARNING: in.GenerateExisting requires manual conversion: does not exist in peer-type
	out.Synchronize = in.Synchronize
	out.OrphanDownstreamOnPolicyDelete = in.OrphanDownstreamOnPolicyDelete
	out.RawData = (*apiextensionsv1.JSON)(unsafe.Pointer(in.RawData))
	out.Clone = in.Clone
	out.CloneList = in.CloneList
	return nil
}

func Convert_v2_Generation_To_v1_Generation(in *Generation, out *kyvernov1.Generation, s apiconversion.Scope) error {
	return autoConvert_v2_Generation_To_v1_Generation(in, out, s)
}

func autoConvert_v1_Attestation_To_v2_Attestation(in *kyvernov1.Attestation, out *Attestation, _ apiconversion.Scope) error {
	if in.PredicateType != "" {
		out.Type = in.PredicateType
	} else {
		out.Type = in.Type
	}
	out.Attestors = *(*[]kyvernov1.AttestorSet)(unsafe.Pointer(&in.Attestors))
	out.Conditions = *(*[]AnyAllConditions)(unsafe.Pointer(&in.Conditions))
	return nil
}

func Convert_v1_Attestation_To_v2_Attestation(in *kyvernov1.Attestation, out *Attestation, s apiconversion.Scope) error {
	return autoConvert_v1_Attestation_To_v2_Attestation(in, out, s)
}

func autoConvert_v2_Attestation_To_v1_Attestation(in *Attestation, out *kyvernov1.Attestation, _ apiconversion.Scope) error {
	out.Type = in.Type
	out.Attestors = *(*[]kyvernov1.AttestorSet)(unsafe.Pointer(&in.Attestors))
	out.Conditions = *(*[]kyvernov1.AnyAllConditions)(unsafe.Pointer(&in.Conditions))
	return nil
}

// Convert_v2_Attestation_To_v1_Attestation is an autogenerated conversion function.
func Convert_v2_Attestation_To_v1_Attestation(in *Attestation, out *kyvernov1.Attestation, s apiconversion.Scope) error {
	return autoConvert_v2_Attestation_To_v1_Attestation(in, out, s)
}

func Convert_v1_JSON_To_v2_AnyAllConditions(in *apiextensionsv1.JSON, out *AnyAllConditions, s apiconversion.Scope) error {
	// In case conditions are specified under any/all resource filters
	if err := json.Unmarshal(in.Raw, out); err == nil {
		return nil
	}

	// In case conditions aren't specified under any/all resource filters
	var conditions []Condition
	if err := json.Unmarshal(in.Raw, &conditions); err == nil {
		out.AnyConditions = conditions
		return nil
	}
	return nil
}

func Convert_v2_AnyAllConditions_To_v1_JSON(in *AnyAllConditions, out *apiextensionsv1.JSON, s apiconversion.Scope) error {
	raw, err := json.Marshal(in)
	if err != nil {
		return err
	}
	out.Raw = raw
	return nil
}
