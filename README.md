### Introduction

This repo is a summary for my recent work. 

My work is to monitorer k8s metrics using prometheus and auto-scaling resources, especially the storage resources.

I have wasted so much time on setting up the environment and related tools. And this is a summary I made.

So: 

1. This is a summary for Kubernetes setup and  installing kube-prometheus quickly（Actually, in just one step）. 

2. Apart from that, this repo also includes storage auto-scaling quick setup yml.
3. In the future, I will include more summary in this repo.



### Environment

Ubuntu18.04 is OK. No try on other system.



### Step

#### Step1: setup Kubernetes 

```shell
./k8s_setup.sh
```

After that, the K8s is set up. Run `kubectl get po` to test that.

---

#### Step2: setup kube-prometheus
```
cd my-kube-prometheus/manifests
kubectl apply -f .
kuebctl apply -f .
cd ../..
```


#### Step3: setup custom-metrics
```
cd custom-metrics-adapter/deploy/manifests/
k apply -f .
cd ../../..
```

---

#### Step4: setup Kube-prometheus

```shell
 k apply -f my-kube-prometheus/manifests/
```

**Note:** Some resources may not create correctly, just because of the wrong setup sequence. So running the upper command twice will solve problem.



And how to access the kube-prometheus? The following:



##### Accessing the dashboards
Prometheus, Grafana和Alertmanager dashboards可以通过kubectl port-forward快速访问:
**(1) Prometheus**

`kubectl --namespace monitoring port-forward svc/prometheus-k8s 9090`
然后可以通过 http://localhost:9090访问Prometheus。

**(2) Grafana**

`kubectl --namespace monitoring port-forward svc/grafana 3000`
然后可以通过 http://localhost:3000访问Grafana，默认的用户名和密码是 admin:admin

**(3) Alert Manager**

`kubectl --namespace monitoring port-forward svc/alertmanager-main 9093`
然后可以通过 http://localhost:9093来访问。

---

#### Step3: setup csi-related resources

```shell
./lvm/csi_setup.sh
```

So, csi-related resources are set up.

The directory `example` show an example of creating pv dynamicly.

