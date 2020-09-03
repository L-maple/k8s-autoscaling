package main

import (
	"encoding/json"
	"flag"
	"github.com/idoubi/goz"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)


var (
	prometheusUrl string

	readsMetricsGaugeValues  map[string]float64       // store deploys' readsIO values
	writesMetricsGaugeValues map[string]float64       // store deploys' writesIO values
	clientSet  *kubernetes.Clientset
	deployReadGauges  map[string]*prometheus.Gauge
	deployWriteGauges map[string]*prometheus.Gauge
)

func getClientSet() *kubernetes.Clientset {
	// get the k8s clientset via config
	var kubeConfig* string

	if home := homedir.HomeDir(); home != "" {
		kubeConfig = flag.String("kubeconfig",
			filepath.Join(home, ".kube", "config"),
			"(optional) absolute path to the kubeconfig file")
	} else {
		kubeConfig = flag.String("kubeconfig",
			"",
			"abosolute path to the kubeconfig path")
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

func getPodsFromDeploy(deployLabels map[string]string, allPods *v1.PodList) ([]string, error) {
	var podNames []string

	for _, pod := range allPods.Items {
		for labelKey, labelValue := range pod.Labels {
			if labelValue == deployLabels[labelKey] {
				podNames = append(podNames, pod.Name)
			}
		}
	}

	return podNames, nil
}

func parseBody(body string) map[string]string {
	var res = make(map[string]string)

	var parsedResult = make(map[string]interface{})
	if err := json.Unmarshal([]byte(body), &parsedResult); err != nil {
		log.Fatal(err)
	}
	parsedResult = parsedResult["data"].(map[string]interface{})
	for _, item := range parsedResult["result"].([]interface{}) {
		itemMap := item.(map[string]interface{})
		metrics := itemMap["metric"].(map[string]interface{})
		pod     := metrics["pod"]

		values  := itemMap["value"].([]interface{})
		value   := values[1].(string)

		res[pod.(string)] = value
	}

	return res
}

func getPodsFsWritesIO() (map[string]string, error) {
	httpCli := goz.NewClient()
	resp, err := httpCli.Get(prometheusUrl + "/api/v1/query", goz.Options{
		Query: map[string]interface{} {
			"query": `sum(rate(container_fs_writes_bytes_total{image!=""}[1m])) by (pod, namespace)`,
		},
		Timeout: 30,
	})
	if err != nil {
		return nil, err
	}
	body, _ := resp.GetBody()

	return parseBody(body.GetContents()), nil
}

func getPodsFsReadsIO() (map[string]string, error) {
	httpCli := goz.NewClient()
	resp, err := httpCli.Get(prometheusUrl + "/api/v1/query", goz.Options{
		Query: map[string]interface{} {
			"query": `sum(rate(container_fs_reads_bytes_total{image!=""}[1m])) by (pod, namespace)`,
		},
		Timeout: 30,
	})
	if err != nil {
		return nil, err
	}
	body, _ := resp.GetBody()

	return parseBody(body.GetContents()), nil
}

func init() {
	readsMetricsGaugeValues  = make(map[string]float64)
	writesMetricsGaugeValues = make(map[string]float64)
	deployReadGauges         = make(map[string]*prometheus.Gauge)
	deployWriteGauges        = make(map[string]*prometheus.Gauge)

	flag.StringVar(&prometheusUrl, "url", "http://localhost:9090", "the prometheus url for cadvisor metrics")
    flag.Parse()
	clientSet = getClientSet()
}

func deployIOMetricsUpdater() {
	for {
		namespaces, err := clientSet.CoreV1().Namespaces().List(metav1.ListOptions{})
		if err != nil {
			log.Fatal(err)
		}

		var readsMetrics  = make(map[string]float64)
		var writesMetrics = make(map[string]float64)
		for _, namespace := range namespaces.Items {
			namespaceName := namespace.Name
			deploys, err := clientSet.AppsV1().Deployments(namespaceName).List(metav1.ListOptions{})
			if err != nil {
				log.Fatal(err)
			}

			var allPods *v1.PodList
			if allPods, err = clientSet.CoreV1().Pods(namespaceName).List(metav1.ListOptions{}); err != nil {
				log.Fatal(err)
			}

			for _, deploy := range deploys.Items {
				labels := deploy.Labels
				deployName := deploy.Name

				pods, err := getPodsFromDeploy(labels, allPods)
				if err != nil {
					continue
				}

				podsWritesIO, err := getPodsFsWritesIO()
				if err != nil {
					continue
				}
				podsReadsIO,  err := getPodsFsReadsIO()
				if err != nil {
					continue
				}
				for _, pod := range pods {
					if podsWritesIO[pod] != "" {
						val, err := strconv.ParseFloat(podsWritesIO[pod], 64)
						if err != nil {
							log.Fatal(err)
						}
						writesMetrics[deployName+"_writes_io"] = writesMetrics[deployName+"_writes_io"] + val
					}
				}

				for _, pod := range pods {
					if podsReadsIO[pod] != "" {
						val, err := strconv.ParseFloat(podsReadsIO[pod], 64)
						if err != nil {
							log.Fatal(err)
						}
						readsMetrics[deployName+"_reads_io"] = readsMetrics[deployName+"_reads_io"] + val
					}
				}
			}
		}

		// update all the global metrics gauge values
		for metricName, value := range writesMetrics {
			metricName = strings.ReplaceAll(metricName, "-", "_")
			writesMetricsGaugeValues[metricName] = value
			//fmt.Println("+++++++++: ", metricName, writesMetricsGaugeValues[metricName])
		}
		for metricName, value := range readsMetrics {
			metricName = strings.ReplaceAll(metricName, "-", "_")
			readsMetricsGaugeValues[metricName] = value
			//fmt.Println("--------: ", metricName, readsMetricsGaugeValues[metricName])
		}
		// re-expose the prometheus metrics' values
		for deployName, gauge := range deployReadGauges {
			//fmt.Println("=========: ", deployName + "_reads_io", readsMetricsGaugeValues[deployName + "_reads_io"])
			(*gauge).Set(readsMetricsGaugeValues[deployName + "_reads_io"])
		}
		for deployName, gauge := range deployWriteGauges {
			(*gauge).Set(writesMetricsGaugeValues[deployName + "_writes_io"])
		}

		time.Sleep(time.Duration(10) * time.Second)
	}
}

func deployIOMetricsRegister() {
	for {
		namespaces, err := clientSet.CoreV1().Namespaces().List(metav1.ListOptions{})
		if err != nil {
			log.Fatal(err)
		}

		for _, namespace := range namespaces.Items {
			namespaceName := namespace.Name
			deploys, _ := clientSet.AppsV1().Deployments(namespaceName).List(metav1.ListOptions{})
			for _, deploy := range deploys.Items {
				deployName := strings.ReplaceAll(deploy.Name, "-", "_")
				if deployReadGauges[deployName] == nil { // if deployment not exist
					readGauge := promauto.NewGauge(prometheus.GaugeOpts{
						Name: deployName + "_reads_io",
						Help: deployName + "_reads_io(bytes)",
					})
					readGauge.Set(0)
					deployReadGauges[deployName] = &readGauge
					_ = prometheus.Register(readGauge)

					writeGauge := promauto.NewGauge(prometheus.GaugeOpts{
						Name: deployName + "_writes_io",
						Help: deployName + "_writes_io(bytes)",
					})
					writeGauge.Set(0)
					deployWriteGauges[deployName] = &writeGauge
					_ = prometheus.Register(writeGauge)
				}
			}
		}

		time.Sleep(time.Duration(60) * time.Second)
	}
}

func main() {
	go deployIOMetricsRegister()

	go deployIOMetricsUpdater()

	http.Handle("/metrics", promhttp.Handler())
	_ = http.ListenAndServe(":20000", nil)
}
