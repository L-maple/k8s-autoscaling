package main

import (
	"context"
	"flag"
	"fmt"
	"k8s.io/api/apps/v1beta1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"path/filepath"
	"time"
)

type StatefulSetInfo struct {
	StatefulSetName      string     // the statefulSet name
	PodNames             []string   // the pods' name
}

func (s *StatefulSetInfo) setStatefulSetName(name string) {
	s.StatefulSetName = name
}

func (s *StatefulSetInfo) getStatefulSetName() string {
	return s.StatefulSetName
}

func (s *StatefulSetInfo) appendPodName(podName string) {
	s.PodNames = append(s.PodNames, podName)
}

func (s *StatefulSetInfo) getPodNames() []string {
	return s.PodNames
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

/*
 * Set StatefulSet's info
 */
func setStsInfo(pods *v1.PodList, statefulSets *v1beta1.StatefulSetList, stsInfo *StatefulSetInfo) {
	/*
	 * Judge whether the input statefulSet exists
	 */
	isFound := false
	matchLabels := make(map[string]string)
	for _, statefulSet := range statefulSets.Items {
		stsName, stsLabels := statefulSet.Name, statefulSet.Spec.Selector.MatchLabels
		if stsName == statefulsetName {
			isFound = true
			matchLabels = stsLabels
			break
		}
	}
	if isFound == false {
		log.Fatal("error: statefulSetName not found")
	}

	/*
	 * Search all pods' name in this statefulSetName
	 */
	for _, pod := range pods.Items {
		podName, podLabels := pod.Name, pod.Labels
		for podLabelKey, podLabelValue := range podLabels {
			if matchLabels[podLabelKey] == podLabelValue {
				stsInfo.appendPodName(podName)
				break
			}

		}
	}
}

func printStsInfo(stsInfo *StatefulSetInfo) {
	fmt.Println(stsInfo.getStatefulSetName())

	for index, podName := range stsInfo.getPodNames() {
		fmt.Println(index, " ", podName)
	}
}


var (
	/* input parameter */
	intervalTime     int
	namespaceName    string
	statefulsetName  string
)

func init() {
	flag.IntVar(&intervalTime, "interval", 15, "exporter interval")
	flag.StringVar(&namespaceName, "namespace", "default", "statefulset's namespace")
	flag.StringVar(&statefulsetName, "statefulset", "default", "statefulset's name")
}

func main() {
	flag.Parse()

	// get k8s clientset
	clientSet := getClientSet()

	// store statefulSet's Pod info
	stsInfo := StatefulSetInfo{}
	stsInfo.setStatefulSetName(statefulsetName)

	podClient := clientSet.CoreV1().Pods(namespaceName)
	podCtx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	pods, err := podClient.List(podCtx, metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	statfulSetClient := clientSet.AppsV1beta1().StatefulSets(namespaceName)
	stsCtx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	statefulSets, err := statfulSetClient.List(stsCtx, metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	/* Set statefulSet's podNames */
	setStsInfo(pods, statefulSets, &stsInfo)

	printStsInfo(&stsInfo)
}
