package main

import (
	"context"
	"flag"
	"fmt"
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
func setStsInfo(clientSet *kubernetes.Clientset, pods *v1.PodList, stsInfo *StatefulSetInfo) {
	statfulSetClient := clientSet.AppsV1().StatefulSets(namespaceName)
	statefulSets, err := statfulSetClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

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

	for {
		// store statefulSet's Pod info
		stsInfo := StatefulSetInfo{}
		stsInfo.setStatefulSetName(statefulsetName)

		podClient := clientSet.CoreV1().Pods(namespaceName)
		pods, err := podClient.List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err)
		}

		/* Set statefulSet's podNames */
		setStsInfo(clientSet, pods, &stsInfo)

		printStsInfo(&stsInfo)

		time.Sleep(time.Duration(intervalTime) * time.Second)
	}
}
