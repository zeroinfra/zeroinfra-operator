/*
Copyright 2025 zeroinfra.

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

// MongoReplicaSetSpec defines the desired state of MongoReplicaSet.
type MongoReplicaSetSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of MongoReplicaSet. Edit mongoreplicaset_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// MongoReplicaSetStatus defines the observed state of MongoReplicaSet.
type MongoReplicaSetStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MongoReplicaSet is the Schema for the mongoreplicasets API.
type MongoReplicaSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MongoReplicaSetSpec   `json:"spec,omitempty"`
	Status MongoReplicaSetStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MongoReplicaSetList contains a list of MongoReplicaSet.
type MongoReplicaSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MongoReplicaSet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MongoReplicaSet{}, &MongoReplicaSetList{})
}
