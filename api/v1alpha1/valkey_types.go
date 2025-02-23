/*
Copyright 2025.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ValkeySpec defines the desired state of Valkey
type ValkeySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Image of Valkey to deploy
	// +kubebuilder:validation:Required
	Image string `json:"image"`

	// Replicas count
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=5
	Replicas int32 `json:"replicas"`

	// User that will be admin
	// +kubebuilder:validation:Required
	User string `json:"user"`

	// Password for admin
	// +kubebuilder:validation:Required
	Password string `json:"password"`

	// UsePersistentVolume for Valkey
	// +kubebuilder:validation:Required
	Volume Volume `json:"volume"`

	// Resource requirements
	// +kubebuilder:validation:Required
	Resource Resource `json:"resource"`
}

type Volume struct {
	// Enabled means that persistent storage should be added
	Enabled bool `json:"enabled"`

	// Storage requirements (e.g., "200Mi", "1Gi", "10Gi", "1Ti")
	// +kubebuilder:validation:Pattern=^[0-9]+[MGT]i$
	Storage string `json:"storage"`
}

type Resource struct {
	// Memory requirements (e.g., "512Mi", "1Gi")
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=^[0-9]+[KMG]i$
	Memory string `json:"memory"`

	// CPU requirements (e.g., "100m", "1", "2.5")
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=^[0-9]+m?$
	CPU string `json:"cpu"`

	// Storage requirements (e.g., "200Mi", "1Gi", "10Gi", "1Ti")
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=^[0-9]+[MGT]i$
	Storage string `json:"storage"`
}

// ValkeyStatus defines the observed state of Valkey
type ValkeyStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Status could be 'running', 'failed', 'stopped'
	Status string `json:"status,omitempty"`
	// ReadyReplicas is a number of working replicas
	ReadyReplicas int32 `json:"ready_replicas"`
	// LastReconcileAt contains timestamp of the last reconcile
	// only if something was changed
	LastReconcileAt *metav1.Time `json:"last_reconcile_at,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Image",type="string",JSONPath=".spec.image"
//+kubebuilder:printcolumn:name="CPU",type="string",JSONPath=".spec.resource.cpu"
//+kubebuilder:printcolumn:name="Memory",type="string",JSONPath=".spec.resource.memory"
//+kubebuilder:printcolumn:name="Has volume",type="boolean",JSONPath=".spec.volume.enabled"
//+kubebuilder:printcolumn:name="Volume size",type="boolean",JSONPath=".spec.volume.enabled"
//+kubebuilder:printcolumn:name="Replicas",type="integer",JSONPath=".spec.replicas"
//+kubebuilder:printcolumn:name="Ready Replicas",type="integer",JSONPath=".status.ready_replicas"
//+kubebuilder:printcolumn:name="Last reconcile",type="date",JSONPath=".status.last_reconcile_at"
//+kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.status"

// Valkey is the Schema for the valkeys API
type Valkey struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ValkeySpec   `json:"spec,omitempty"`
	Status ValkeyStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ValkeyList contains a list of Valkey
type ValkeyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Valkey `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Valkey{}, &ValkeyList{})
}
