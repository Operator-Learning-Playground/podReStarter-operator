package controller

import (
	"context"
	podrestarterv1alpha1 "github.com/myoperator/podrestartoperator/pkg/apis/podrestarter/v1alpha1"
	"github.com/myoperator/podrestartoperator/pkg/k8sconfig"
	"github.com/myoperator/podrestartoperator/pkg/restart"
	"github.com/myoperator/podrestartoperator/pkg/util"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type PodReStarterController struct {
	client client.Client
	// 事件发送器
	EventRecorder record.EventRecorder
}

func NewPodReStarterController(client client.Client, event record.EventRecorder) *PodReStarterController {
	return &PodReStarterController{client: client, EventRecorder: event}
}

const (
	PodReStarterStatusSuccess = "Successful"
	PodReStarterStatusRunning = "Running"
	PodReStarterStatusFailed  = "Failure"

	ReStarter = "restart"
	Upgrade   = "upgrade"
)

// Reconcile 调协loop
func (r *PodReStarterController) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {

	podReStarter := &podrestarterv1alpha1.Podrestarter{}
	err := r.client.Get(ctx, req.NamespacedName, podReStarter)
	if err != nil {
		if client.IgnoreNotFound(err) != nil {
			klog.Error("get podReStarter error: ", err)
			return reconcile.Result{}, err
		}
		// 如果未找到的错误，不再进入调协
		return reconcile.Result{}, nil
	}

	// 代表升级已经成功，不再进入调协
	if podReStarter.Status.Status == PodReStarterStatusSuccess {
		klog.Info("restart or upgrade already successful...")
		return reconcile.Result{}, nil
	}

	if podReStarter.Spec.Type == ReStarter {
		podReStarter.Status.Type = "ReStarter"
	} else {
		podReStarter.Status.Type = "Upgrade"
	}
	podReStarter.Status.Status = PodReStarterStatusRunning

	err = r.client.Status().Update(ctx, podReStarter)
	if err != nil {
		klog.Error("update status error: ", err)
		return reconcile.Result{}, err
	}

	podList := util.GetPodsByDeployment(podReStarter.Spec.DeploymentName,
		podReStarter.Spec.DeploymentNamespace, k8sconfig.ClientSet)
	if len(podList) == 0 {
		klog.Error("nothing in podList...: ", err)
		return reconcile.Result{}, nil
	}
	// 避免用户填错
	var num = podReStarter.Spec.Replicas
	if len(podList) < podReStarter.Spec.Replicas {
		num = len(podList)
	}

	if podReStarter.Spec.Type == ReStarter {
		// 原地重启逻辑
		for i := 0; i < num; i++ {
			// pod原地重启
			err := restart.RestartPodByImage(ctx, &podList[i], k8sconfig.ClientSet)
			if err != nil {
				r.EventRecorder.Eventf(podReStarter, v1.EventTypeWarning, "PodReStarter Failed", "ReStarter Pod Failed")
				podReStarter.Status.Status = PodReStarterStatusFailed
				err = r.client.Status().Update(ctx, podReStarter)
				if err != nil {
					klog.Error("update status error: ", err)
					return reconcile.Result{}, err
				}
				return reconcile.Result{}, err
			}
		}

		// 重起成功后，修改状态
		podReStarter.Status.Status = PodReStarterStatusSuccess
		err = r.client.Status().Update(ctx, podReStarter)
		if err != nil {
			klog.Error("update status error: ", err)
			return reconcile.Result{}, err
		}
	} else if podReStarter.Spec.Type == Upgrade {
		imageList := make([]string, 0)
		for _, v := range podReStarter.Spec.Images {
			imageList = append(imageList, v.Image)
		}
		for i := 0; i < num; i++ {
			// pod原地升级
			// image 一定要按照原本的顺序
			err := restart.UpgradePodByImages(ctx, &podList[i], k8sconfig.ClientSet, imageList)
			if err != nil {
				r.EventRecorder.Eventf(podReStarter, v1.EventTypeWarning, "PodReStarter Failed", "Upgrade Pod Image Failed")
				podReStarter.Status.Status = PodReStarterStatusFailed
				err = r.client.Status().Update(ctx, podReStarter)
				if err != nil {
					klog.Error("update status error: ", err)
					return reconcile.Result{}, err
				}
				return reconcile.Result{}, err
			}
		}

		// 升级成功后，修改状态
		podReStarter.Status.Status = PodReStarterStatusSuccess
		err = r.client.Status().Update(ctx, podReStarter)
		if err != nil {
			klog.Error("update status error: ", err)
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}
