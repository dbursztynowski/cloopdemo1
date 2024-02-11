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

// ClosedLoopSpec defines the desired state of ClosedLoop
type ClosedLoopSpec struct {
	// Important: Run "make" to regenerate code after modifying this file

	// This is where you define the Spec you want the ClosedLoop to have
	Name       string           `json:"name,omitempty"`
	Monitoring MonitoringObject `json:"monitoring"` // You can add Structure as Spec field as i did with MonitoringObject and DecisionObject
	Decision   DecisionObject   `json:"decision"`
	Execution  Execution        `json:"execution"`
}

// ClosedLoopStatus defines the observed state of ClosedLoop
type ClosedLoopStatus struct {
	// In the ClosedLoop Controllers We don't use Status field but it's possible, here it's a example of how to define a status field
	Name string `json:"name,omitempty"`
	IncreaseRank string `json:"increaserank,omitempty"`
	IncreaseTime string `json:"increasetime,omitempty"`

	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ClosedLoop is the Schema for the closedloops API
type ClosedLoop struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClosedLoopSpec   `json:"spec,omitempty"`
	Status ClosedLoopStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ClosedLoopList contains a list of ClosedLoop
type ClosedLoopList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClosedLoop `json:"items"`
}

type MonitoringObject struct {
	Affix              string             `json:"affix,omitempty"`
	MonitoringKind     MonitoringKind     `json:"monitoringkind"`
	MonitoringPolicies MonitoringPolicies `json:"monitorinpolicies"`
}

type DecisionObject struct {
	Affix            string           `json:"affix,omitempty"`
	DecisionKind     DecisionKind     `json:"decisionkind"`
	DecisionPolicies DecisionPolicies `json:"decisionpolicies"`
}

type MonitoringKind struct {
	MonitoringKindName string `json:"monitoringkindname"`
	Source             Source `json:"source,omitempty"`
	RequestedPod       bool   `json:"requestedpod,omitempty"`
}

type DecisionKind struct {
	DecisionKindName string `json:"decisionkindname"`
}

func init() {
	SchemeBuilder.Register(&ClosedLoop{}, &ClosedLoopList{})
}
