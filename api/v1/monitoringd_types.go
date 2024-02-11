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

// MonitoringDSpec defines the desired state of MonitoringD
type MonitoringDSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of MonitoringD. Edit monitoringd_types.go to remove/update
	Source             SourceD             `json:"source"`
	Affix              string              `json:"affix,omitempty"`
	DecisionKind       string              `json:"decisionkind"`
	MonitoringPolicies MonitoringDPolicies `json:"monitoringpolicies"`
}

// MonitoringDStatus defines the observed state of MonitoringD
type MonitoringDStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Affix string `json:"affix,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// MonitoringD is the Schema for the monitoringds API
type MonitoringD struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MonitoringDSpec   `json:"spec,omitempty"`
	Status MonitoringDStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MonitoringDList contains a list of MonitoringD
type MonitoringDList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MonitoringD `json:"items"`
}

//-------Additional Structure ---------//

type MonitoringDPolicies struct {
	Data          map[string]string `json:"data"`
	TresholdKind  map[string]string `json:"tresholdkind"`
	TresholdValue map[string]string `json:"tresholdvalue"`
}

type SourceD struct {
	Addresse string `json:"addresse"`
	Port     int32  `json:"port"`
	Interval int32  `json:"interval"`
}

func init() {
	SchemeBuilder.Register(&MonitoringD{}, &MonitoringDList{})
}
