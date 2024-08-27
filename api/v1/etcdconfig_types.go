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

// EtcdConfigSpec defines the desired state of EtcdConfig
type EtcdConfigSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// field of EtcdConfig. Edit etcdconfig_types.go to remove/update
	Items []EtcdConfigItem `json:"items,omitempty"`
}

// EtcdConfigItem represents a key-value pair in the spec
type EtcdConfigItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// EtcdConfigStatus defines the observed state of EtcdConfig
type EtcdConfigStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	LastSyncedTime metav1.Time `json:"lastSyncedTime,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// EtcdConfig is the Schema for the etcdconfigs API
type EtcdConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EtcdConfigSpec   `json:"spec,omitempty"`
	Status EtcdConfigStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// EtcdConfigList contains a list of EtcdConfig
type EtcdConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EtcdConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EtcdConfig{}, &EtcdConfigList{})
}
