package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Podrestarter
type Podrestarter struct {
	metav1.TypeMeta `json:",inline"`

	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec PodReStarterSpec `json:"spec,omitempty"`
}

type PodReStarterSpec struct {
	Restart             string         `json:"restart"`
	DeploymentName      string         `json:"deployment_name"`
	DeploymentNamespace string         `json:"deployment_namespace"`
	Replicas            int            `json:"replicas"`
	Images              []Image        `json:"images"`
	Containers          []v1.Container `json:"containers"`
}

type Image struct {
	Image string `json:"image"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PodReStarterList
type PodrestarterList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Podrestarter `json:"items"`
}
