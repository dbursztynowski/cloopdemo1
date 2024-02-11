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

// ClosedLoopDSpec defines the desired state of ClosedLoopD
type ClosedLoopDSpec struct {
	// Important: Run "make" to regenerate code after modifying this file

	// This is where you define the Spec you want the ClosedLoop to have
	Name       string            `json:"name,omitempty"`
	Monitoring MonitoringDObject `json:"monitoring"` // You can add Structure as Spec field as i did with MonitoringObject and DecisionObject
	Decision   DecisionDObject   `json:"decision"`
	Execution  ExecutionD        `json:"execution"`
}

// ClosedLoopDStatus defines the observed state of ClosedLoopD
type ClosedLoopDStatus struct {
	// In the ClosedLoop Controllers We don't use Status field but it's possible, here it's a example of how to define a status field
	Name string `json:"name,omitempty"`
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ClosedLoopD is the Schema for the closedloopds API
type ClosedLoopD struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClosedLoopDSpec   `json:"spec,omitempty"`
	Status ClosedLoopDStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ClosedLoopDList contains a list of ClosedLoopD
type ClosedLoopDList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClosedLoopD `json:"items"`
}

type MonitoringDObject struct {
	Affix              string              `json:"affix,omitempty"`
	MonitoringKind     MonitoringDKind     `json:"monitoringkind"`
}

type DecisionDObject struct {
	Affix            string            `json:"affix,omitempty"`
	DecisionKind     DecisionDKind     `json:"decisionkind"`
	DecisionPolicies DecisionDPolicies `json:"decisionpolicies"`
}

type MonitoringDKind struct {
	MonitoringKindName string `json:"monitoringkindname"`
	Source             Source `json:"source,omitempty"`
	RequestedPod       bool   `json:"requestedpod,omitempty"`
}

type DecisionDKind struct {
	DecisionKindName string `json:"decisionkindname"`
}

func init() {
	SchemeBuilder.Register(&ClosedLoopD{}, &ClosedLoopDList{})
}
