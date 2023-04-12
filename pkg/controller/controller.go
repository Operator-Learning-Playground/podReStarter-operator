package controller

import (
	"context"
	podrestarterv1alpha1 "github.com/myoperator/podrestartoperator/pkg/apis/podrestarter/v1alpha1"
	"github.com/myoperator/podrestartoperator/pkg/k8sconfig"
	"github.com/myoperator/podrestartoperator/pkg/restart"
	"github.com/myoperator/podrestartoperator/pkg/util"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)


type PodReStarterController struct {
	client.Client

}

func NewPodReStarterController() *PodReStarterController {
	return &PodReStarterController{}
}

// Reconcile 调协loop
func (r *PodReStarterController) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {

	podReStarter := &podrestarterv1alpha1.Podrestarter{}
	err := r.Get(ctx, req.NamespacedName, podReStarter)
	if err != nil {
		return reconcile.Result{}, err
	}
	klog.Info(podReStarter)
	podList := util.GetPodsByDeployment(podReStarter.Spec.DeploymentName,
		podReStarter.Spec.DeploymentNamespace, k8sconfig.ClientSet)
	// 避免用户填错
	var num = podReStarter.Spec.Replicas
	if len(podList) < podReStarter.Spec.Replicas {
		num = len(podList)
	}
	if podReStarter.Spec.Restart == "true" {
		// 原地重启逻辑
		for i := 0; i < num; i++ {
			// pod原地重启
			restart.RestartPodByImage(&podList[i], k8sconfig.ClientSet)
		}

	} else {
		imageList := make([]string, 0)
		for _, v := range podReStarter.Spec.Images {
			imageList = append(imageList, v.Image)
		}
		for i := 0; i < num; i++ {
			// pod原地升级
			// image 一定要按照原本的顺序
			restart.UpgradePodByImage(&podList[i], k8sconfig.ClientSet, imageList)

		}
	}

	return reconcile.Result{}, nil
}

// InjectClient 使用controller-runtime 需要注入的client
func(r *PodReStarterController) InjectClient(c client.Client) error {
	r.Client = c
	return nil
}



