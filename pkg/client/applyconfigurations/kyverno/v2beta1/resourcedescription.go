/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v2beta1

import (
	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ResourceDescriptionApplyConfiguration represents an declarative configuration of the ResourceDescription type for use
// with apply.
type ResourceDescriptionApplyConfiguration struct {
	Kinds             []string                       `json:"kinds,omitempty"`
	Names             []string                       `json:"names,omitempty"`
	Namespaces        []string                       `json:"namespaces,omitempty"`
	Annotations       map[string]string              `json:"annotations,omitempty"`
	Selector          *v1.LabelSelector              `json:"selector,omitempty"`
	NamespaceSelector *v1.LabelSelector              `json:"namespaceSelector,omitempty"`
	Operations        []kyvernov1.AdmissionOperation `json:"operations,omitempty"`
}

// ResourceDescriptionApplyConfiguration constructs an declarative configuration of the ResourceDescription type for use with
// apply.
func ResourceDescription() *ResourceDescriptionApplyConfiguration {
	return &ResourceDescriptionApplyConfiguration{}
}

// WithKinds adds the given value to the Kinds field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Kinds field.
func (b *ResourceDescriptionApplyConfiguration) WithKinds(values ...string) *ResourceDescriptionApplyConfiguration {
	for i := range values {
		b.Kinds = append(b.Kinds, values[i])
	}
	return b
}

// WithNames adds the given value to the Names field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Names field.
func (b *ResourceDescriptionApplyConfiguration) WithNames(values ...string) *ResourceDescriptionApplyConfiguration {
	for i := range values {
		b.Names = append(b.Names, values[i])
	}
	return b
}

// WithNamespaces adds the given value to the Namespaces field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Namespaces field.
func (b *ResourceDescriptionApplyConfiguration) WithNamespaces(values ...string) *ResourceDescriptionApplyConfiguration {
	for i := range values {
		b.Namespaces = append(b.Namespaces, values[i])
	}
	return b
}

// WithAnnotations puts the entries into the Annotations field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the entries provided by each call will be put on the Annotations field,
// overwriting an existing map entries in Annotations field with the same key.
func (b *ResourceDescriptionApplyConfiguration) WithAnnotations(entries map[string]string) *ResourceDescriptionApplyConfiguration {
	if b.Annotations == nil && len(entries) > 0 {
		b.Annotations = make(map[string]string, len(entries))
	}
	for k, v := range entries {
		b.Annotations[k] = v
	}
	return b
}

// WithSelector sets the Selector field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Selector field is set to the value of the last call.
func (b *ResourceDescriptionApplyConfiguration) WithSelector(value v1.LabelSelector) *ResourceDescriptionApplyConfiguration {
	b.Selector = &value
	return b
}

// WithNamespaceSelector sets the NamespaceSelector field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the NamespaceSelector field is set to the value of the last call.
func (b *ResourceDescriptionApplyConfiguration) WithNamespaceSelector(value v1.LabelSelector) *ResourceDescriptionApplyConfiguration {
	b.NamespaceSelector = &value
	return b
}

// WithOperations adds the given value to the Operations field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Operations field.
func (b *ResourceDescriptionApplyConfiguration) WithOperations(values ...kyvernov1.AdmissionOperation) *ResourceDescriptionApplyConfiguration {
	for i := range values {
		b.Operations = append(b.Operations, values[i])
	}
	return b
}
