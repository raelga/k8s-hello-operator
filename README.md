# Kubernetes Operator for raelga/hello-http

Sample kubernetes operator for raelga/hello-http using the Operator Framerwok.

Is a work in progress and only for learning purposes.

```
k8s-hello-operator on  master [⇡?] 
❯ kubectl create -f deploy/crds/hellohttp_v1alpha1_hellohttpservice_crd.yaml 
customresourcedefinition.apiextensions.k8s.io/hellohttpservices.hellohttp.rael.io created
```

```
k8s-hello-operator on  master [⇡?] 
❯ kubectl create -f deploy/                                                 
deployment.apps/k8s-hello-operator created
role.rbac.authorization.k8s.io/k8s-hello-operator created
rolebinding.rbac.authorization.k8s.io/k8s-hello-operator created
serviceaccount/k8s-hello-operator created
```

```
k8s-hello-operator on  master [⇡?] 
❯ kubectl create -f deploy/crds/hellohttp_v1alpha1_hellohttpservice_cr.yaml 
hellohttpservice.hellohttp.rael.io/example-hellohttpservice created
```

```
k8s-hello-operator on  master [⇡?] 
❯ k get pods                                                               
NAME                                 READY   STATUS              RESTARTS   AGE
example-hellohttpservice-rs-jm88m    0/1     ContainerCreating   0          2s
k8s-hello-operator-c95cdf577-nfddb   1/1     Running             0          8s
```

```
k8s-hello-operator on  master [⇡?] 
❯ k get pods
NAME                                 READY   STATUS    RESTARTS   AGE
example-hellohttpservice-rs-jm88m    1/1     Running   0          4s
k8s-hello-operator-c95cdf577-nfddb   1/1     Running   0          10s
```

```
k8s-hello-operator on  master [⇡?] 
❯ k get rs  
NAME                           DESIRED   CURRENT   READY   AGE
example-hellohttpservice-rs    1         1         1       10s
k8s-hello-operator-c95cdf577   1         1         1       17s
```