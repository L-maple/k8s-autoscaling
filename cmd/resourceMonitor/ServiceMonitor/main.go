package main

import (
	"flag"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

type podInfo struct {
	BossName      string     // the deployment or statefulSet resource name
	BossType      string     // ordinator type: (1) "deployment" (2) "statefulset"
	ContainerInfo []string
}

func (p *podInfo) setBossName(bossName string) {
	p.BossName = bossName
}

func (p *podInfo) getBossName() string {
	return p.BossName
}

func (p *podInfo) setBossType(bossType string) {
	p.BossType = bossType
}

func (p *podInfo) getBossType() string {
	return p.BossType
}


func getClientSet() *kubernetes.Clientset {
	// get the k8s clientset via config
	var kubeConfig* string

	if home := homedir.HomeDir(); home != "" {
		kubeConfig = flag.String("kubeconfig",
			filepath.Join(home, ".kube", "config"),
			"(optional) absolute path to the kubeconfig file");
	} else {
		kubeConfig = flag.String("kubeconfig",
			"",
			"abosolute path to the kubeconfig path");
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeConfig)
	if err != nil {
		panic(err)
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	return clientSet
}

func setPodMapInfoForDeployment(podMap map[string]*podInfo, clientSet *kubernetes.Clientset, pods *v1.PodList) {
	deployClient := clientSet.AppsV1().Deployments("default")
	deploys, err := deployClient.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, pod := range pods.Items {
		podName, podLabels := pod.Name, pod.Labels
		for podLabelKey, podLabelValue := range podLabels {
			for _, deploy := range deploys.Items {
				deployName, deployLabels := deploy.Name, deploy.Spec.Selector.MatchLabels
				if deployLabels[podLabelKey] == podLabelValue {
					podInfoStruct := new(podInfo)
					podInfoStruct.setBossName(deployName)
					podInfoStruct.setBossType("deployment")
					podMap[podName] = podInfoStruct
					break
				}
			}
		}
	}
}

/**
	function: set Pod's podInfo for StatefulSet
 */
func setPodMapInfoForStatefulSet(podMap map[string]*podInfo, clientSet *kubernetes.Clientset, pods *v1.PodList) {
	statfulSetClient := clientSet.AppsV1().StatefulSets("default")
	statefulSets, err := statfulSetClient.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, pod := range pods.Items {
		podName, podLabels := pod.Name, pod.Labels
		for podLabelKey, podLabelValue := range podLabels {
			for _, statefulSet := range statefulSets.Items {
				statefulSetName, statefulSetLabels := statefulSet.Name, statefulSet.Spec.Selector.MatchLabels
				if statefulSetLabels[podLabelKey] == podLabelValue {
					podInfoStruct := new(podInfo)
					podInfoStruct.setBossName(statefulSetName)
					podInfoStruct.setBossType("statefulset")
					podMap[podName] = podInfoStruct
					break
				}
			}
		}
	}
}

func printPodInfo(podMap map[string]*podInfo) {
	for podName, podInfo := range podMap {
		fmt.Println(podName, " => ", podInfo.BossName, ",Type: ", podInfo.BossType)
	}
}

func main() {
	// get k8s clientset
	clientSet := getClientSet()

	// map pod to deployment/statefulSet
	podMap := make(map[string]*podInfo)

	podClient := clientSet.CoreV1().Pods("default")
	pods, err := podClient.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	// pod -- Deployment
	setPodMapInfoForDeployment(podMap, clientSet, pods)

	// pod -- statefulSet
	setPodMapInfoForStatefulSet(podMap, clientSet, pods)

	printPodInfo(podMap)

	// pod -- containers
 
}
