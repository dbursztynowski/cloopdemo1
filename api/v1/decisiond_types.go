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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DecisionDSpec defines the desired state of DecisionD
type DecisionDSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of DecisionD. Edit decisiond_types.go to remove/update
	Affix            string            `json:"affix,omitempty"`
	Time             string            `json:"time,omitempty"`
	Data		     map[string]string `json:"data,omitempty"`
}

// DecisionDStatus defines the observed state of DecisionD
type DecisionDStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Affix string `json:"affix,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// DecisionD is the Schema for the decisionds API
type DecisionD struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DecisionDSpec   `json:"spec,omitempty"`
	Status DecisionDStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DecisionDList contains a list of DecisionD
type DecisionDList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DecisionD `json:"items"`
}

type DecisionDPolicies struct {
	DecisionType string        `json:"decisiontype"`
	PrioritySpec PriorityDSpec `json:"priorityspec,omitempty"`
}

type PriorityDSpec struct {
	PriorityType string            `json:"prioritytype,omitempty"`
	PriorityRank map[string]string `json:"priorityrank,omitempty"`
}

func init() {
	SchemeBuilder.Register(&DecisionD{}, &DecisionDList{})
}
