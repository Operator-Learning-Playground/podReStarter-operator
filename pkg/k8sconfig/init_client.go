package k8sconfig

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
)

// InitClient 初始化 client
func InitClient(config *rest.Config) kubernetes.Interface {
	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatal(err)
	}
	return c
}

var ClientSet kubernetes.Interface

func init() {
	ClientSet = InitClient(K8sRestConfig())
}
