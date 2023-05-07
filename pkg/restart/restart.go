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
	"log"
	"time"
)

// RestartPodByImage 原地重启pod的方式
func RestartPodByImage(pod *v1.Pod, clientSet kubernetes.Interface) {
	klog.Info("pod is restarting...")

	restartImage := pod.Spec.Containers[0].Image

	// 改成任意一个镜像
	randomImage := "nginx:1.18-alpine"
	if restartImage == randomImage {
		randomImage = "nginx:1.19-alpine"
	}
	patch := fmt.Sprintf(`[{"op": "replace", "path": "/spec/containers/0/image", "value": "%v"}]`, randomImage)
	patchBytes := []byte(patch)

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
		log.Fatalln(err)
	}

	// 延迟
	time.Sleep(time.Second * 30)

	// 再次使用patch换回原来image
	restartPatch := fmt.Sprintf(`[{"op": "replace", "path": "/spec/containers/0/image", "value": "%v"}]`, restartImage)
	restartPatchBytes := []byte(restartPatch)

	restartJsonPatch, err := jsonpatch.DecodePatch(restartPatchBytes)
	if err != nil {
		klog.Error("DecodePatch error: ", err)
		return
	}
	restartJsonPatchBytes, err := json.Marshal(restartJsonPatch)
	if err != nil {
		klog.Error("json Marshal error: ", err)
		return
	}
	_, err = clientSet.CoreV1().Pods(pod.Namespace).
		Patch(context.TODO(), pod.Name, types.JSONPatchType,
			restartJsonPatchBytes, metav1.PatchOptions{})
	if err != nil {
		log.Fatalln(err)
	}
}
