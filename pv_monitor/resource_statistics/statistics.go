package main

import (
	"github.com/idoubi/goz"
	"k8s.io/apimachinery/pkg/util/json"
	"log"
	"strconv"
	"time"
)

type PodStatistics struct {
	podName, namespace string
}

func (c *PodStatistics) getUtilizationQuery(query string) []interface{}  {
	curl := PromCurl{endpoint, namespace, nil}
	responseBody, err := curl.Get("/api/v1/query_range", goz.Options{
		Query: map[string]interface{}{
			"query": query,
			"start": strconv.FormatInt(time.Now().Unix() - 3600, 10),
			"end": strconv.FormatInt(time.Now().Unix(), 10),
			"step": "14",
		},
	})
	if err != nil {
		log.Fatal("curl.Get error")
	}

	jsonData := make(map[string]interface{})
	err = json.Unmarshal(responseBody, &jsonData)
	if err != nil {
		log.Fatal("json.Unmarshal: ", err)
	}

	data := jsonData["data"].(map[string]interface{})
	values := data["result"].([]interface{})[0].(map[string]interface{})["values"]
	valuesSlice := values.([]interface{})
	for _, tsAndUtilization := range valuesSlice {
		//res := tsAndUtilization.([]interface{})
		//timestamp, utilization := res[0].(int64), res[1].(string)
		_ = tsAndUtilization.([]interface{})
	}

	return values.([]interface{})
}

func (c *PodStatistics) GetAvgCpuUtilizationQuery() []interface{} {
	query := "sum(rate(container_cpu_usage_seconds_total{" + "image!=\"\", " +
			 "pod=" + "\""+ c.podName +"\", namespace=\"" + c.namespace +
			 "\"}[1m]))"

	result := c.getUtilizationQuery(query)
	return result
}

func (c *PodStatistics) GetAvgMemoryUtilizationQuery() []interface{} {
	query := "sum(container_memory_rss{" + "image!=\"\", " +
		"pod=" + "\""+ c.podName +"\", namespace=\"" + c.namespace +
		"\"})"

	result := c.getUtilizationQuery(query)
	return result
}

func (c *PodStatistics) GetAvgDiskUtilizationQuery() []interface{} {
	query := "disk_utilization{pod=" + "\"" + c.podName +  "\", namespace=\"" + c.namespace + "\"}"

	result := c.getUtilizationQuery(query)

	return result
}

