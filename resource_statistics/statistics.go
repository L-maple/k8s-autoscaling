package main

import (
	"fmt"
	"github.com/idoubi/goz"
	"k8s.io/apimachinery/pkg/util/json"
	"log"
	"time"
)

type PodStatistics struct {
	podName, namespace string
}

func (c *PodStatistics) GetAvgCpuUtilizationQuery() string {
	query := "sum(rate(container_cpu_usage_seconds_total{" +
			 "pod=" + c.podName +", namespace=" + c.namespace +
			 "}[1m]))"

	curl := PromCurl{endpoint, namespace, nil}
	responseBody, err := curl.Get("/api/v1/query", goz.Options{
		Query: map[string]interface{}{
			"query": query,
			"time": time.Now().Unix(),
		},
	})
	if err != nil {
		log.Fatal("curl.Get error")
	}

	fmt.Println(responseBody)

	jsonData := make(map[string]interface{})
	err = json.Unmarshal(responseBody, &jsonData)
	if err != nil {
		log.Fatal("json.Unmarshal: ", err)
	}

	data := jsonData["data"].(map[string]interface{})

	fmt.Println(data["resultType"])

	return ""
}


