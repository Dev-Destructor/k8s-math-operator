/*
Copyright 2024.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ArithmeticSpec defines the desired state of Arithmetic
type ArithmeticSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Arithmetic. Edit arithmetic_types.go to remove/update
	Expression string `json:"expression,omitempty"`
}

// ArithmeticStatus defines the observed state of Arithmetic
type ArithmeticStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Answer string `json:"answer,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Arithmetic is the Schema for the arithmetics API
type Arithmetic struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ArithmeticSpec   `json:"spec,omitempty"`
	Status ArithmeticStatus `json:"status,omitempty"`
}

// DeepCopyObject implements client.Object.
// func (*Arithmetic) DeepCopyObject() runtime.Object {
// 	panic("unimplemented")
// }

//+kubebuilder:object:root=true

// ArithmeticList contains a list of Arithmetic
type ArithmeticList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Arithmetic `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Arithmetic{}, &ArithmeticList{})
}
