package restart

import (
	"context"
	"encoding/json"
	"fmt"
	jsonpatch "github.com/evanphx/json-patch"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
)

type patchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

// UpgradePodByImages 原地升级pod镜像
func UpgradePodByImages(pod *v1.Pod, clientSet kubernetes.Interface, images []string) {
	klog.Info("pod is upgrading !!")
	patchList := make([]*patchOperation, 0)
	for k, image := range images {
		p := &patchOperation{
			Op:    "replace",
			Path:  fmt.Sprintf("/spec/containers/%v/image", k),
			Value: image,
		}
		patchList = append(patchList, p)

	}
	patchBytes, err := json.Marshal(patchList)
	if err != nil {
		klog.Error(err)
		return
	}

	jsonPatch, err := jsonpatch.DecodePatch(patchBytes)
	if err != nil {
		klog.Error("DecodePatch error: ", err)
		return
	}
	jsonPatchBytes, err := json.Marshal(jsonPatch)
	if err != nil {
		klog.Error("json Marshal error: ", err)
		return
	}
	_, err = clientSet.CoreV1().Pods(pod.Namespace).
		Patch(context.TODO(), pod.Name, types.JSONPatchType,
			jsonPatchBytes, metav1.PatchOptions{})
	if err != nil {
		klog.Error("pod patch error: ", err)
		return
	}
}
