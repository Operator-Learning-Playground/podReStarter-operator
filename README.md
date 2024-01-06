## podrestart-operator 
### pod 原地升级镜像原地重启容器 operator

### 项目思路与设计
设计背景：服务部署时常常会对已有服务进行镜像更新或容器重启等操作，如果直接使用 kubectl apply 命令，会让 Pod 重新进入调度流程，导致重新调度节点或 ip 重新分配等。
在此原因上开发此 operator，让原地操作更加方便。

思路：使用 k8s-clientSet 中的 patch 操作实现。


### 项目功能
1. 支持 Deployment 中的 Pod **原地升级**镜像
2. 支持 Deployment 中的 Pod **原地重启**容器
3. 支持修改部分副本

```yaml
apiVersion: api.practice.com/v1alpha1
kind: Podrestarter
metadata:
  name: mypodrestarter
spec:
  type: restart                     # restart 代表pod原地重启容器 upgrade 代表原地升级镜像
  # 目前仅支持 deployment 操作
  deployment_name: my-deployment    # deployment name
  deployment_namespace: default     # deployment namespace
  replicas: 2                       # 升级或重启的副本数，数量需要小于等于 deployment 副本数
  images:                           # 升级镜像，镜像顺序一定要 "按照原顺序"
    - image: nginx:1.19-alpine
    - image: busybox
```

```bash
[root@VM-0-16-centos podrestarteroperator]# kubectl get podrestarters.api.practice.com
NAME             TYPE        AGE
mypodrestarter   ReStarter   9m34s
mypodupgrade     Upgrade     57s
```

附注：使用 patch 原地升级 pod 有一些限制，需要特别注意，否则会遇到 patch error
参考 [issue](https://github.com/Operator-Learning-Playground/podReStarter-operator/issues/1)
