package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Podrestarter
type Podrestarter struct {
	metav1.TypeMeta `json:",inline"`

	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PodReStarterSpec   `json:"spec,omitempty"`
	Status PodReStarterStatus `json:"status"`
}

type PodReStarterSpec struct {
	Type                string  `json:"type"`
	DeploymentName      string  `json:"deployment_name"`
	DeploymentNamespace string  `json:"deployment_namespace"`
	Replicas            int     `json:"replicas"`
	Images              []Image `json:"images"`
}

type Image struct {
	Image string `json:"image"`
}

type PodReStarterStatus struct {
	// Type 类型 restart upgrade 两种
	Type string `json:"type"`
	// Status 进行状态
	Status string `json:"status"`
	// Image 镜像
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
