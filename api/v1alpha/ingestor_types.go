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

package v1alpha

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// IngestorSpec defines the desired state of Ingestor.
type IngestorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	SourceNamespace string          `json:"sourceNamespace,omitempty"`
	Destination     Destination     `json:"destination,omitempty"`
	Interval        metav1.Duration `json:"interval,omitempty"`
}

type Destination struct {
	Endpoint     string `json:"endpoint,omitempty"`
	Region       string `json:"region,omitempty"`
	BucketName   string `json:"bucketName,omitempty"`
	ObjectPrefix string `json:"objectPrefix,omitempty"`
	AccessKey    string `json:"accessKey,omitempty"`
	SecretKey    string `json:"secretKey,omitempty"`
}

// IngestorStatus defines the observed state of Ingestor.
type IngestorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Format=date-time
	LastIngested metav1.Time `json:"lastIngested,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Ingestor is the Schema for the ingestors API.
type Ingestor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IngestorSpec   `json:"spec,omitempty"`
	Status IngestorStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// IngestorList contains a list of Ingestor.
type IngestorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Ingestor `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Ingestor{}, &IngestorList{})
}
