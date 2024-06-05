package v2

import (
	unsafe "unsafe"

	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	v1alpha1 "k8s.io/api/admissionregistration/v1alpha1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiconversion "k8s.io/apimachinery/pkg/conversion"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

func (src *ClusterPolicy) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*kyvernov1.ClusterPolicy)

	return Convert_v2_ClusterPolicy_To_v1_ClusterPolicy(src, dst, nil)
}

func (dst *ClusterPolicy) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*kyvernov1.ClusterPolicy)

	return Convert_v1_ClusterPolicy_To_v2_ClusterPolicy(src, dst, nil)
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

func Convert_v2_Spec_To_v1_Spec(in *Spec, out *kyvernov1.Spec, s apiconversion.Scope) error {
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
	out.ApplyRules = in.ApplyRules
	out.FailurePolicy = in.FailurePolicy
	out.Background = in.Background
	out.Admission = in.Admission
	out.WebhookTimeoutSeconds = in.WebhookTimeoutSeconds
	out.UseServerSideApply = in.UseServerSideApply
	// TODO: check if we need to move this field for v2
	out.WebhookConfiguration = (*kyvernov1.WebhookConfiguration)(in.WebhookConfiguration)
	// TODO: handle deprecated fields like SchemaValidation, GenerateExistingOnPolicyUpdate

	// NOTE: GENERATED CODE
	// 	if in.Rules != nil {
	// 		in, out := &in.Rules, &out.Rules
	// 		*out = make([]v1.Rule, len(*in))
	// 		for i := range *in {
	// 			if err := Convert_v2_Rule_To_v1_Rule(&(*in)[i], &(*out)[i], s); err != nil {
	// 				return err
	// 			}
	// 		}
	// 	} else {
	// 		out.Rules = nil
	// 	}
	// 	out.ApplyRules = (*v1.ApplyRulesType)(unsafe.Pointer(in.ApplyRules))
	// 	out.FailurePolicy = (*v1.FailurePolicyType)(unsafe.Pointer(in.FailurePolicy))
	// 	out.Admission = (*bool)(unsafe.Pointer(in.Admission))
	// 	out.Background = (*bool)(unsafe.Pointer(in.Background))
	// 	out.WebhookTimeoutSeconds = (*int32)(unsafe.Pointer(in.WebhookTimeoutSeconds))
	// 	out.UseServerSideApply = in.UseServerSideApply
	// 	out.WebhookConfiguration = (*v1.WebhookConfiguration)(unsafe.Pointer(in.WebhookConfiguration))
	return nil
}

func Convert_v1_Spec_To_v2_Spec(in *kyvernov1.Spec, out *Spec, s apiconversion.Scope) error {
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
	out.ApplyRules = in.ApplyRules
	out.FailurePolicy = in.FailurePolicy
	out.Background = in.Background
	out.Admission = in.Admission
	out.WebhookTimeoutSeconds = in.WebhookTimeoutSeconds
	out.UseServerSideApply = in.UseServerSideApply
	out.WebhookConfiguration = (*WebhookConfiguration)(in.WebhookConfiguration)

	// NOTE: GENERATED CODE
	// 	if in.Rules != nil {
	// 		in, out := &in.Rules, &out.Rules
	// 		*out = make([]Rule, len(*in))
	// 		for i := range *in {
	// 			if err := Convert_v1_Rule_To_v2_Rule(&(*in)[i], &(*out)[i], s); err != nil {
	// 				return err
	// 			}
	// 		}
	// 	} else {
	// 		out.Rules = nil
	// 	}
	// 	out.ApplyRules = (*v1.ApplyRulesType)(unsafe.Pointer(in.ApplyRules))
	// 	out.FailurePolicy = (*v1.FailurePolicyType)(unsafe.Pointer(in.FailurePolicy))
	// 	// WARNING: in.ValidationFailureAction requires manual conversion: does not exist in peer-type
	// 	// WARNING: in.ValidationFailureActionOverrides requires manual conversion: does not exist in peer-type
	// 	out.Admission = (*bool)(unsafe.Pointer(in.Admission))
	// 	out.Background = (*bool)(unsafe.Pointer(in.Background))
	// 	// WARNING: in.SchemaValidation requires manual conversion: does not exist in peer-type
	// 	out.WebhookTimeoutSeconds = (*int32)(unsafe.Pointer(in.WebhookTimeoutSeconds))
	// 	// WARNING: in.MutateExistingOnPolicyUpdate requires manual conversion: does not exist in peer-type
	// 	// WARNING: in.GenerateExistingOnPolicyUpdate requires manual conversion: does not exist in peer-type
	// 	// WARNING: in.GenerateExisting requires manual conversion: does not exist in peer-type
	// 	out.UseServerSideApply = in.UseServerSideApply
	// 	out.WebhookConfiguration = (*WebhookConfiguration)(unsafe.Pointer(in.WebhookConfiguration))
	return nil
}

func Convert_v1_ResourceDescription_To_v2_ResourceDescription(in *kyvernov1.ResourceDescription, out *ResourceDescription, s apiconversion.Scope) error {
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

func Convert_v1_ImageVerification_To_v2_ImageVerification(in *kyvernov1.ImageVerification, out *ImageVerification, s apiconversion.Scope) error {
	out.Type = kyvernov1.ImageVerificationType(in.Type)
	// WARNING: in.Image requires manual conversion: does not exist in peer-type
	out.ImageReferences = *(*[]string)(unsafe.Pointer(&in.ImageReferences))
	out.SkipImageReferences = *(*[]string)(unsafe.Pointer(&in.SkipImageReferences))
	// WARNING: in.Key requires manual conversion: does not exist in peer-type
	// WARNING: in.Roots requires manual conversion: does not exist in peer-type
	// WARNING: in.Subject requires manual conversion: does not exist in peer-type
	// WARNING: in.Issuer requires manual conversion: does not exist in peer-type
	// WARNING: in.AdditionalExtensions requires manual conversion: does not exist in peer-type
	out.Attestors = *(*[]kyvernov1.AttestorSet)(unsafe.Pointer(&in.Attestors))
	if in.Attestations != nil {
		in, out := &in.Attestations, &out.Attestations
		*out = make([]Attestation, len(*in))
		for i := range *in {
			if err := Convert_v1_Attestation_To_v2_Attestation(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Attestations = nil
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

func Convert_v1_MatchResources_To_v2_MatchResources(in *kyvernov1.MatchResources, out *MatchResources, s apiconversion.Scope) error {
	if in.Any != nil {
		in, out := &in.Any, &out.Any
		*out = make(ResourceFilters, len(*in))
		for i := range *in {
			if err := Convert_v1_ResourceFilter_To_v2_ResourceFilter(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Any = nil
	}
	if in.All != nil {
		in, out := &in.All, &out.All
		*out = make(ResourceFilters, len(*in))
		for i := range *in {
			if err := Convert_v1_ResourceFilter_To_v2_ResourceFilter(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.All = nil
	}
	// WARNING: in.UserInfo requires manual conversion: does not exist in peer-type
	// WARNING: in.ResourceDescription requires manual conversion: does not exist in peer-type
	return nil
}

func Convert_v2_Validation_To_v1_Validation(in *Validation, out *kyvernov1.Validation, s apiconversion.Scope) error {
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

func Convert_v2_Mutation_To_v1_Mutation(in *Mutation, out *kyvernov1.Mutation, s apiconversion.Scope) error {
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

func Convert_v2_Generation_To_v1_Generation(in *Generation, out *kyvernov1.Generation, s apiconversion.Scope) error {
	out.ResourceSpec = in.ResourceSpec
	// WARNING: in.GenerateExisting requires manual conversion: does not exist in peer-type
	out.Synchronize = in.Synchronize
	out.OrphanDownstreamOnPolicyDelete = in.OrphanDownstreamOnPolicyDelete
	out.RawData = (*apiextensionsv1.JSON)(unsafe.Pointer(in.RawData))
	out.Clone = in.Clone
	out.CloneList = in.CloneList
	return nil
}

func Convert_v1_Attestation_To_v2_Attestation(in *kyvernov1.Attestation, out *Attestation, s apiconversion.Scope) error {
	// WARNING: in.PredicateType requires manual conversion: does not exist in peer-type
	out.Type = in.Type
	out.Attestors = *(*[]kyvernov1.AttestorSet)(unsafe.Pointer(&in.Attestors))
	out.Conditions = *(*[]AnyAllConditions)(unsafe.Pointer(&in.Conditions))
	return nil
}

func autoConvert_v1_TargetResourceSpec_To_v2_TargetResourceSpec(in *kyvernov1.TargetResourceSpec, out *TargetResourceSpec, _ apiconversion.Scope) error {
	out.ResourceSpec = in.ResourceSpec
	out.Context = *(*[]kyvernov1.ContextEntry)(unsafe.Pointer(&in.Context))
	// if in.RawAnyAllConditions != nil {
	// 	in, out := &in.RawAnyAllConditions, &out.RawAnyAllConditions
	// 	*out = new(AnyAllConditions)
	// 	if err := Convert_v1_JSON_To_v2_AnyAllConditions(*in, *out, s); err != nil {
	// 		return err
	// 	}
	// } else {
	// 	out.RawAnyAllConditions = nil
	// }
	return nil
}

func Convert_v1_TargetResourceSpec_To_v2_TargetResourceSpec(in *kyvernov1.TargetResourceSpec, out *TargetResourceSpec, s apiconversion.Scope) error {
	return autoConvert_v1_TargetResourceSpec_To_v2_TargetResourceSpec(in, out, s)

}

func autoConvert_v2_TargetResourceSpec_To_v1_TargetResourceSpec(in *TargetResourceSpec, out *kyvernov1.TargetResourceSpec, _ apiconversion.Scope) error {
	out.ResourceSpec = in.ResourceSpec
	out.Context = *(*[]kyvernov1.ContextEntry)(unsafe.Pointer(&in.Context))
	// if in.RawAnyAllConditions != nil {
	// 	in, out := &in.RawAnyAllConditions, &out.RawAnyAllConditions
	// 	*out = new(apiextensionsv1.JSON)
	// 	if err := Convert_v2_AnyAllConditions_To_v1_JSON(*in, *out, s); err != nil {
	// 		return err
	// 	}
	// } else {
	// 	out.RawAnyAllConditions = nil
	// }
	return nil
}

func Convert_v2_TargetResourceSpec_To_v1_TargetResourceSpec(in *TargetResourceSpec, out *kyvernov1.TargetResourceSpec, s apiconversion.Scope) error {
	return autoConvert_v2_TargetResourceSpec_To_v1_TargetResourceSpec(in, out, s)
}

func autoConvert_v1_Rule_To_v2_Rule(in *kyvernov1.Rule, out *Rule, s apiconversion.Scope) error {
	out.Name = in.Name
	out.Context = *(*[]kyvernov1.ContextEntry)(unsafe.Pointer(&in.Context))
	if err := Convert_v1_MatchResources_To_v2_MatchResources(&in.MatchResources, &out.MatchResources, s); err != nil {
		return err
	}
	if err := Convert_v1_MatchResources_To_v2_MatchResources(&in.ExcludeResources, &out.ExcludeResources, s); err != nil {
		return err
	}
	out.ImageExtractors = *(*kyvernov1.ImageExtractorConfigs)(unsafe.Pointer(&in.ImageExtractors))
	// if in.RawAnyAllConditions != nil {
	// 	in, out := &in.RawAnyAllConditions, &out.RawAnyAllConditions
	// 	*out = new(AnyAllConditions)
	// 	if err := Convert_v1_JSON_To_v2_AnyAllConditions(*in, *out, s); err != nil {
	// 		return err
	// 	}
	// } else {
	// 	out.RawAnyAllConditions = nil
	// }
	out.CELPreconditions = *(*[]admissionregistrationv1.MatchCondition)(unsafe.Pointer(&in.CELPreconditions))
	if err := Convert_v1_Mutation_To_v2_Mutation(&in.Mutation, &out.Mutation, s); err != nil {
		return err
	}
	if err := Convert_v1_Validation_To_v2_Validation(&in.Validation, &out.Validation, s); err != nil {
		return err
	}
	if err := Convert_v1_Generation_To_v2_Generation(&in.Generation, &out.Generation, s); err != nil {
		return err
	}
	if in.VerifyImages != nil {
		in, out := &in.VerifyImages, &out.VerifyImages
		*out = make([]ImageVerification, len(*in))
		for i := range *in {
			if err := Convert_v1_ImageVerification_To_v2_ImageVerification(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.VerifyImages = nil
	}
	out.SkipBackgroundRequests = in.SkipBackgroundRequests
	return nil
}

func Convert_v1_Rule_To_v2_Rule(in *kyvernov1.Rule, out *Rule, s apiconversion.Scope) error {
	return autoConvert_v1_Rule_To_v2_Rule(in, out, s)
}

func autoConvert_v2_Rule_To_v1_Rule(in *Rule, out *kyvernov1.Rule, s apiconversion.Scope) error {
	out.Name = in.Name
	out.Context = *(*[]kyvernov1.ContextEntry)(unsafe.Pointer(&in.Context))
	if err := Convert_v2_MatchResources_To_v1_MatchResources(&in.MatchResources, &out.MatchResources, s); err != nil {
		return err
	}
	if err := Convert_v2_MatchResources_To_v1_MatchResources(&in.ExcludeResources, &out.ExcludeResources, s); err != nil {
		return err
	}
	// out.ImageExtractors = *(*kyvernov1.ImageExtractorConfigs)(unsafe.Pointer(&in.ImageExtractors))
	// if in.RawAnyAllConditions != nil {
	// 	in, out := &in.RawAnyAllConditions, &out.RawAnyAllConditions
	// 	*out = new(apiextensionsv1.JSON)
	// 	if err := Convert_v2_AnyAllConditions_To_v1_JSON(*in, *out, s); err != nil {
	// 		return err
	// 	}
	// } else {
	// 	out.RawAnyAllConditions = nil
	// }
	out.CELPreconditions = *(*[]v1alpha1.MatchCondition)(unsafe.Pointer(&in.CELPreconditions))
	if err := Convert_v2_Mutation_To_v1_Mutation(&in.Mutation, &out.Mutation, s); err != nil {
		return err
	}
	if err := Convert_v2_Validation_To_v1_Validation(&in.Validation, &out.Validation, s); err != nil {
		return err
	}
	if err := Convert_v2_Generation_To_v1_Generation(&in.Generation, &out.Generation, s); err != nil {
		return err
	}
	if in.VerifyImages != nil {
		in, out := &in.VerifyImages, &out.VerifyImages
		*out = make([]kyvernov1.ImageVerification, len(*in))
		for i := range *in {
			if err := Convert_v2_ImageVerification_To_v1_ImageVerification(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.VerifyImages = nil
	}
	out.SkipBackgroundRequests = in.SkipBackgroundRequests
	return nil
}

func Convert_v2_Rule_To_v1_Rule(in *Rule, out *kyvernov1.Rule, s apiconversion.Scope) error {
	return autoConvert_v2_Rule_To_v1_Rule(in, out, s)

}

func autoConvert_v2_Deny_To_v1_Deny(in *Deny, out *kyvernov1.Deny, s apiconversion.Scope) error {
	// if in.RawAnyAllConditions != nil {
	// 	var conditions []kyvernov1.Condition
	// 	for _, any := range in.RawAnyAllConditions.AnyConditions {
	// 		var condition kyvernov1.Condition
	// 		if err := Convert_v2_Condition_To_v1_Condition(&any, &condition, s); err != nil {
	// 			return err
	// 		}
	// 		conditions = append(conditions, condition)
	// 	}
	// 	for _, all := range in.RawAnyAllConditions.AllConditions {
	// 		var condition kyvernov1.Condition
	// 		if err := Convert_v2_Condition_To_v1_Condition(&all, &condition, s); err != nil {
	// 			return err
	// 		}
	// 		conditions = append(conditions, condition)
	// 	}
	// 	jsonBytes, err := json.Marshal(conditions)
	// 	if err != nil {
	// 		return fmt.Errorf("error occurred while marshalling: %+v", err)
	// 	}
	// 	out.SetAnyAllConditions(jsonBytes)
	// }
	return nil
	// if in.RawAnyAllConditions != nil {
	// 	// Marshal the entire AnyAllConditions struct into JSON
	// 	jsonBytes, err := json.Marshal(in.RawAnyAllConditions)
	// 	if err != nil {
	// 		return fmt.Errorf("error occurred while marshalling: %+v", err)
	// 	}
	// 	out.SetAnyAllConditions(jsonBytes)
	// }
	//return nil
}

func Convert_v2_Deny_To_v1_Deny(in *Deny, out *kyvernov1.Deny, s apiconversion.Scope) error {
	return autoConvert_v2_Deny_To_v1_Deny(in, out, s)
}

func autoConvert_v1_Deny_To_v2_Deny(in *kyvernov1.Deny, out *Deny, s apiconversion.Scope) error {
	// if in.RawAnyAllConditions != nil {
	// 	var anyConditions []Condition
	// 	conditions := in.GetAnyAllConditions()
	// 	// marshalling the abstract apiextensions.JSON back to JSON form
	// 	jsonByte, err := json.Marshal(conditions)
	// 	if err != nil {
	// 		return fmt.Errorf("error occurred while marshalling: %+v", err)
	// 	}

	// 	var v1conditions []kyvernov1.Condition
	// 	err = json.Unmarshal(jsonByte, &v1conditions)
	// 	if err != nil {
	// 		return fmt.Errorf("error occurred while unmarshalling: %+v", err)
	// 	}

	// 	for _, cond := range v1conditions {
	// 		var v2condition Condition
	// 		if err := Convert_v1_Condition_To_v2_Condition(&cond, &v2condition, s); err != nil {
	// 			return err
	// 		}
	// 		anyConditions = append(anyConditions, v2condition)
	// 	}
	// 	anyAllConditions := &AnyAllConditions{
	// 		AnyConditions: anyConditions,
	// 	}
	// 	out.RawAnyAllConditions = anyAllConditions
	// }
	return nil
}

func Convert_v1_Deny_To_v2_Deny(in *kyvernov1.Deny, out *Deny, s apiconversion.Scope) error {
	return autoConvert_v1_Deny_To_v2_Deny(in, out, s)
}
