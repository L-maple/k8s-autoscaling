package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"net"
	"path/filepath"
	"time"

	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	pb "github.com/k8s-autoscaling/pv_monitor/pv_monitor"
)


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

	/* port */
	promPort = ":30001"          /* For RequestPVNames grpc */
	pvRequestPort = ":30002"     /* For ReplyPVInfos grpc */

	/* metric name to expose */
	diskUtilizationMetric = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "disk_utilization_total", // TODO: 命名不规范
		Help: "pv_disk_utilization_total",
	})

	/* global statefulSet's Pod info */
	stsInfoGlobal StatefulSetInfo

	/* disk utilization */
	diskUtilizations map[string]float64
)

func init() {
	flag.IntVar(&intervalTime, "interval", 15, "exporter interval")
	flag.StringVar(&namespaceName, "namespace", "default", "statefulset's namespace")
	flag.StringVar(&statefulsetName, "statefulset", "default", "statefulset's name")
}


func setDiskUtilizationMetric() {
	go func() {
		if diskUtilizations == nil {
			diskUtilizationMetric.Set(0)
		} else {
			diskUtilizationTotal := 0.0
			for _, diskUtilization := range diskUtilizations {
				diskUtilizationTotal += diskUtilization
			}
			diskUtilizationMetric.Set(diskUtilizationTotal)
		}

		time.Sleep(time.Duration(intervalTime) * time.Second)
	}()
}

func recordStsInfo(clientSet *kubernetes.Clientset) {
	go func() {
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

			stsInfoGlobal = stsInfo

			printStsInfo(&stsInfo)

			time.Sleep(time.Duration(intervalTime) * time.Second)
		}
	}()
}

func RegisterPVRequestServer() {
	go func() {
		lis, err := net.Listen("tcp", pvRequestPort)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		pb.RegisterPVServiceServer(s, &server{})
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}

func main() {
	flag.Parse()

	/* get k8s clientset */
	var clientSet *kubernetes.Clientset

	//clientSet = getInClusterClientSet()
	clientSet = getClientSet()

	/* Record StatefulSet information */
	recordStsInfo(clientSet)

	/* Register grpc server */
	RegisterPVRequestServer()

	/* Set disk utilization metric & exposed at 30001 */
	setDiskUtilizationMetric()
	http.Handle("/metrics", promhttp.Handler())
	_ = http.ListenAndServe(promPort, nil)
}
