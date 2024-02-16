package v2

import (
	"testing"

	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	kyvernov2beta1 "github.com/kyverno/kyverno/api/kyverno/v2beta1"
	"gotest.tools/assert"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func Test_MatchResources(t *testing.T) {
	testCases := []struct {
		name       string
		namespaced bool
		subject    MatchResources
		errors     []string
	}{{
		name:       "valid",
		namespaced: true,
		subject: MatchResources{
			Any: kyvernov2beta1.ResourceFilters{{
				UserInfo: kyvernov1.UserInfo{
					Subjects: []rbacv1.Subject{{
						Kind:      "ServiceAccount",
						Namespace: "ns",
						Name:      "sa-1",
					}},
				},
			}},
		},
	}, {
		name:       "any-all",
		namespaced: true,
		subject: MatchResources{
			Any: kyvernov2beta1.ResourceFilters{{
				UserInfo: kyvernov1.UserInfo{
					Subjects: []rbacv1.Subject{{
						Kind:      "ServiceAccount",
						Namespace: "ns",
						Name:      "sa-1",
					}},
				},
			}},
			All: kyvernov2beta1.ResourceFilters{{
				UserInfo: kyvernov1.UserInfo{
					Subjects: []rbacv1.Subject{{
						Kind:      "ServiceAccount",
						Namespace: "ns",
						Name:      "sa-1",
					}},
				},
			}},
		},
		errors: []string{
			`dummy: Invalid value: v2.MatchResources{Any:v2beta1.ResourceFilters{v2beta1.ResourceFilter{UserInfo:v1.UserInfo{Roles:[]string(nil), ClusterRoles:[]string(nil), Subjects:[]v1.Subject{v1.Subject{Kind:"ServiceAccount", APIGroup:"", Name:"sa-1", Namespace:"ns"}}}, ResourceDescription:v2beta1.ResourceDescription{Kinds:[]string(nil), Names:[]string(nil), Namespaces:[]string(nil), Annotations:map[string]string(nil), Selector:(*v1.LabelSelector)(nil), NamespaceSelector:(*v1.LabelSelector)(nil), Operations:[]v1.AdmissionOperation(nil)}}}, All:v2beta1.ResourceFilters{v2beta1.ResourceFilter{UserInfo:v1.UserInfo{Roles:[]string(nil), ClusterRoles:[]string(nil), Subjects:[]v1.Subject{v1.Subject{Kind:"ServiceAccount", APIGroup:"", Name:"sa-1", Namespace:"ns"}}}, ResourceDescription:v2beta1.ResourceDescription{Kinds:[]string(nil), Names:[]string(nil), Namespaces:[]string(nil), Annotations:map[string]string(nil), Selector:(*v1.LabelSelector)(nil), NamespaceSelector:(*v1.LabelSelector)(nil), Operations:[]v1.AdmissionOperation(nil)}}}}: Can't specify any and all together`,
		},
	}}

	path := field.NewPath("dummy")
	for _, testCase := range testCases {
		errs := testCase.subject.Validate(path, testCase.namespaced, nil)
		assert.Equal(t, len(errs), len(testCase.errors))
		for i, err := range errs {
			assert.Equal(t, err.Error(), testCase.errors[i])
		}
	}
}
