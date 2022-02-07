# k8s WASM FaaS Builder Controller

Builer Controller 主要通过监听 ConfigMaps 中的 FaaS 函数资源，下发编译任务将函数编译为 WASM 可执行文件，并将 WASM 可执行文件传输到指定的地方

## 编译与部署

Builer Controller 以 Deployment 的形式部署在 k8s 集群中

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: k8s-wasm-faas-builder-controller
  namespace: default
  labels:
    app: k8s-wasm-faas-builder-controller
spec:
  replicas: 1
  selector:
    matchLabels:
      app: k8s-wasm-faas-builder-controller
  template:
    metadata:
      labels:
        app: k8s-wasm-faas-builder-controller
    spec:
      serviceAccountName: faas-wasm
      nodeSelector:
        faas-wasm-runtime: wasm
      containers:
        - name: k8s-faas-builder-controller
          image: redxiiikk/k8s-faas-builder-controller:latest
          env:
            - name: TARGET_HOST
              value: "x.x.x.x"
            - name: TARGET_PORT
              value: '22'
            - name: TARGET_USER
              value: root
            - name: TARGET_PATH
              value: /root/wasm/hello
```

> 为了监听 ConfigMaps 和创建 Builder Job , 需要提供一个服务账号，需要具有 ConfigMaps 和 Jobs 的权限.
>
> ```yaml
> kind: Role
> apiVersion: rbac.authorization.k8s.io/v1
> metadata:
>   namespace: default
>   name: faas-wasm
> rules:
>   - apiGroups: [ "", "batch" ]
>     resources: [ "configmaps", "jobs" ]
>     verbs: [ "get", "list", "watch", "create", "update", "patch", "delete" ]
> ---
> kind: RoleBinding
> apiVersion: rbac.authorization.k8s.io/v1
> metadata:
>   name: role-test-account-binding
>   namespace: default
> subjects:
>   - kind: ServiceAccount
>     name: faas-wasm
>     namespace: default
> roleRef:
>   kind: Role
>   name: faas-wasm
>   apiGroup: rbac.authorization.k8s.io
> ```
