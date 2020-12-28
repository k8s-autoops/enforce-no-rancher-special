# enforce-no-rancher-special

禁用一些 Rancher 自定义的功能，防止出现与标准 Kubernetes 不符合的行为

## 使用方式

* 初始化 `admission-bootstrapper` 
  参照此文档 https://github.com/k8s-autoops/admission-bootstrapper ，完成 `admission-bootstrapper` 的初始化步骤
* 部署以下 YAML

```yaml
# create serviceaccount
apiVersion: v1
kind: ServiceAccount
metadata:
  name: enforce-no-rancher-special
  namespace: autoops
---
# create clusterrole
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
  name: enforce-no-rancher-special
rules:
  - apiGroups: [""]
    resources: ["namespaces"]
    verbs: ["get"]
---
# create clusterrolebinding
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: enforce-no-rancher-special
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: enforce-no-rancher-special
subjects:
  - kind: ServiceAccount
    name: enforce-no-rancher-special
    namespace: autoops
---
# create job
apiVersion: batch/v1
kind: Job
metadata:
  name: install-enforce-no-rancher-special
  namespace: autoops
spec:
  template:
    spec:
      serviceAccount: admission-bootstrapper
      containers:
        - name: admission-bootstrapper
          image: autoops/admission-bootstrapper
          env:
            - name: ADMISSION_NAME
              value: enforce-no-rancher-special
            - name: ADMISSION_IMAGE
              value: autoops/enforce-no-rancher-special
            - name: ADMISSION_ENVS
              value: ""
            - name: ADMISSION_SERVICE_ACCOUNT
              value: "enforce-no-rancher-special"
            - name: ADMISSION_MUTATING
              value: "true"
            - name: ADMISSION_IGNORE_FAILURE
              value: "false"
            - name: ADMISSION_SIDE_EFFECT
              value: "None"
            - name: ADMISSION_RULES
              value: '[{"operations":["CREATE"],"apiGroups":[""], "apiVersions":["*"], "resources":["services"]},{"operations":["CREATE","UPDATE"],"apiGroups":["*"], "apiVersions":["*"], "resources":["ingresses"]}]'
      restartPolicy: OnFailure
```

## Credits

Guo Y.K., MIT License
