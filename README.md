## podrestart-operator 
### 基于k8s-operator的pod原地升级原地重启operator

### 项目思路与设计
设计背景：服务部署时常常会对已有服务进行镜像更新或容器重启等操作，如果直接使用kubectl apply命令，会让Pod重新进入调度流程，导致重新调度节点或ip重新分配等。
在此原因上开发此operator，让原地操作更加方便。

思路：使用k8s-clientSet中的patch操作实现。


### 项目功能
1. 支持deployment中的Pod**原地升级**镜像
2. 支持deployment中的Pod**原地重启**服务
```yaml
apiVersion: api.practice.com/v1alpha1
kind: Podrestarter
metadata:
  name: mypodrestarter
spec:
  restart: "false"                  # 布尔值：true代表pod原地重启 false代表原地升级
  deployment_name: my-deployment    # deployment名
  deployment_namespace: default     # deployment namespace
  replicas: 2                       # 升级或重启的副本数
  images:                           # 升级镜像，镜像顺序一定要 "按照原顺序"
    - image: nginx:1.19-alpine
    - image: busybox
```

