/*
Copyright 2023.

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

package v1

import (
	// appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// corev1 "k8s.io/api/core/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MicroDevSpec defines the desired state of MicroDev
type MicroDevSpec struct {
	//+kubebuilder:validation:MinLength=0
	QueryUrl string `json:"queryUrl"`

	Selector *metav1.LabelSelector `json:"selector"`

	//+kubebuilder:validation:Minimum=1
	// +optional
	QueryInterval *int32 `json:"queryInterval,omitempty"`
}

// MicroDevStatus defines the observed state of MicroDev
type MicroDevStatus struct {
	// +optional
	Replicas *int32 `json:"replicas"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// MicroDev is the Schema for the microdevs API
type MicroDev struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MicroDevSpec   `json:"spec,omitempty"`
	Status MicroDevStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MicroDevList contains a list of MicroDev
type MicroDevList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MicroDev `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MicroDev{}, &MicroDevList{})
}
