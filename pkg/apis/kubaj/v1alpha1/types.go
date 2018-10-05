package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type SecretKMSList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []SecretKMS `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type SecretKMS struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              SecretKMSSpec   `json:"spec"`
	Status            SecretKMSStatus `json:"status,omitempty"`
}

type SecretKMSSpec struct {
	Secret   string      `json:"secret"`
	File     string      `json:"file"`
	Provider KMSProvider `json:"provider"`
}

type KMSProvider struct {
	GoogleCloud *GoogleCloudProvider `json:"google-cloud"`
}

type GoogleCloudProvider struct {
	Project  string `json:"project"`
	Location string `json:"location"`
	Keyring  string `json:"keyring"`
	Key      string `json:"key"`
	Data     string `json:"data"`
}

type SecretKMSStatus struct {
	// Fill me
}
