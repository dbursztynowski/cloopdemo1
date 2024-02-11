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

// Monitoringv2Spec defines the desired state of Monitoringv2
type Monitoringv2Spec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// This is where you define what Spec you want for your CR Monitoringv2
	RequestedPod       bool               `json:"requestedpod"`
	Affix              string             `json:"affix,omitempty"`
	Time               string             `json:"time,omitempty"`
	DecisionKind       string             `json:"decisionkind"`
	MonitoringPolicies MonitoringPolicies `json:"monitoringpolicies"`
	Data               map[string]string  `json:"data,omitempty"`
}

// Monitoringv2Status defines the observed state of Monitoringv2
type Monitoringv2Status struct {
	Affix string `json:"affix,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Monitoringv2 is the Schema for the monitoringv2s API
type Monitoringv2 struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   Monitoringv2Spec   `json:"spec,omitempty"`
	Status Monitoringv2Status `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// Monitoringv2List contains a list of Monitoringv2
type Monitoringv2List struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Monitoringv2 `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Monitoringv2{}, &Monitoringv2List{})
}
