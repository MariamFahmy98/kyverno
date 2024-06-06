package v2beta1

import (
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

func (src *ClusterPolicy) ConvertTo(dstRaw conversion.Hub) error {
	return nil
}

func (dst *ClusterPolicy) ConvertFrom(srcRaw conversion.Hub) error {
	return nil
}

func (src *ClusterPolicyList) ConvertTo(dstRaw conversion.Hub) error {
	return nil
}

func (dst *ClusterPolicyList) ConvertFrom(srcRaw conversion.Hub) error {
	return nil
}

func (src *Policy) ConvertTo(dstRaw conversion.Hub) error {
	return nil
}

func (dst *Policy) ConvertFrom(srcRaw conversion.Hub) error {
	return nil
}

func (src *PolicyList) ConvertTo(dstRaw conversion.Hub) error {
	return nil
}

func (dst *PolicyList) ConvertFrom(srcRaw conversion.Hub) error {
	return nil
}
