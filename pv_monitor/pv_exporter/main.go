package main

import (
	"context"
	"flag"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"path/filepath"
	"time"
)

type PodInfo struct {
	PVCNames             []string
	PVNames              []string
}

func (p *PodInfo) SetPVCNames(pvcNames []string) {
	p.PVCNames = pvcNames
}

func (p PodInfo) GetPVCNames() []string {
	return p.PVCNames
}

func (p *PodInfo) SetPVNames(pvNames []string) {
	p.PVNames = pvNames
}

func (p PodInfo) GetPVNames() []string {
	return p.PVNames
}


type StatefulSetInfo struct {
	StatefulSetName      string                /* the statefulSet name */
	PodNames             []string              /* the pods' name       */
	PodInfos             map[string]PodInfo    /* store the pod's info */
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

func (s *StatefulSetInfo) initializePodInfos() {
	s.PodInfos = make(map[string]PodInfo)
}

func (s *StatefulSetInfo) setPodInfo(podName string, podInfo PodInfo) {
	s.PodInfos[podName] = podInfo
}

func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func getInClusterClientSet() *kubernetes.Clientset {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return clientSet
}

func getClientSet() *kubernetes.Clientset {
	/* get the k8s clientset via config */
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
		log.Fatal("Error: statefulSetName not found")
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

	/* Initialize the pvc's names */
	setPodInfos(clientSet, pods, stsInfo)
}


func setPodInfos(clientSet *kubernetes.Clientset, pods *v1.PodList, stsInfo *StatefulSetInfo) {
	stsInfo.initializePodInfos()

	/* get all pods' pvc and pv info*/
	var podInfo PodInfo
	for _, pod := range pods.Items {
		var pvcNames []string

		for _,volume := range pod.Spec.Volumes {
			if volume.PersistentVolumeClaim == nil {
				continue
			}
			pvcName := volume.PersistentVolumeClaim.ClaimName
			pvcNames = append(pvcNames, pvcName)
		}
		podInfo.SetPVCNames(pvcNames)

		/* get a pvcs' pv info*/
		pvcClient := clientSet.CoreV1().PersistentVolumeClaims(namespaceName)
		pvcs, err := pvcClient.List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Fatal("Error: pvClient List error")
		}
		var pvNames []string
		for _, pvc := range pvcs.Items {
			pvcName, pvName := pvc.Name, pvc.Spec.VolumeName
			if _, found := Find(pvcNames, pvcName); found {
				pvNames = append(pvNames, pvName)
			}
		}
		podInfo.SetPVNames(pvNames)

		stsInfo.PodInfos[pod.Name] = podInfo
	}
}

func printStsInfo(stsInfo *StatefulSetInfo) {
	fmt.Println(stsInfo.getStatefulSetName())

	for _, podName := range stsInfo.getPodNames() {
		/* Print pod name */
		fmt.Println(podName)

		/* Print pvc names */
		for _, pvcName := range stsInfo.PodInfos[podName].GetPVCNames() {
			fmt.Print(pvcName, " ")
		}
		fmt.Println()

		/* Print pv names */
		for _, pvName := range stsInfo.PodInfos[podName].GetPVNames() {
			fmt.Print(pvName, " ")
		}
		fmt.Println()
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

	/* get k8s clientset */
	clientSet := getInClusterClientSet()

	for {
		/* store statefulSet's Pod info */
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
