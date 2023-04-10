package controller

import (
	"context"
	podrestarterv1alpha1 "github.com/myoperator/podrestartoperator/pkg/apis/podrestarter/v1alpha1"
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

	podReStarter := &podrestarterv1alpha1.PodReStarter{}
	err := r.Get(ctx, req.NamespacedName, podReStarter)
	if err != nil {
		return reconcile.Result{}, err
	}
	klog.Info(podReStarter)

	if podReStarter.Spec.Restart == "true" {
		// 原地重启逻辑
	} else {
		// 原地升级逻辑
	}

	return reconcile.Result{}, nil
}

// InjectClient 使用controller-runtime 需要注入的client
func(r *PodReStarterController) InjectClient(c client.Client) error {
	r.Client = c
	return nil
}



