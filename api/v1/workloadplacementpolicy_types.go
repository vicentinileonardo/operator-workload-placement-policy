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

// WorkloadPlacementPolicySpec defines the desired state of WorkloadPlacementPolicy.
type WorkloadPlacementPolicySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	OriginRegion  Region `json:"originRegion"`
	MaxLatency    int    `json:"maxLatency"`
	CloudProvider string `json:"cloudProvider"`
}

type Region struct {
	CloudProviderRegion   string `json:"cloudProviderRegion"`
	ISOCountryCodeA2      string `json:"isoCountryCodeA2"`      // +kubebuilder:validation:Optional
	PhysicalLocation      string `json:"physicalLocation"`      // +kubebuilder:validation:Optional
	ElectricityMapsRegion string `json:"electricityMapsRegion"` // +kubebuilder:validation:Optional
	//Coordinates           Coordinates `json:"coordinates"`           // +kubebuilder:validation:Optional
}

type Coordinates struct {
	Latitude  int `json:"latitude"`
	Longitude int `json:"longitude"`
}

// WorkloadPlacementPolicyStatus defines the observed state of WorkloadPlacementPolicy.
type WorkloadPlacementPolicyStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	EligibleRegions []Region `json:"eligibleRegions"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// WorkloadPlacementPolicy is the Schema for the workloadplacementpolicies API.
type WorkloadPlacementPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WorkloadPlacementPolicySpec   `json:"spec,omitempty"`
	Status WorkloadPlacementPolicyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// WorkloadPlacementPolicyList contains a list of WorkloadPlacementPolicy.
type WorkloadPlacementPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []WorkloadPlacementPolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&WorkloadPlacementPolicy{}, &WorkloadPlacementPolicyList{})
}
