package statistics

import (
	"github.com/idoubi/goz"
	"k8s.io/apimachinery/pkg/util/json"
	"log"
	"strconv"
	"time"
)

type PodStatistics struct {
	PodName, Namespace string
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
			 "pod=" + "\""+ c.PodName +"\", namespace=\"" + c.Namespace +
			 "\"}[1m]))"

	result := c.getUtilizationQuery(query)
	return result
}

func (c *PodStatistics) GetAvgMemoryUtilizationQuery() []interface{} {
	query := "sum(container_memory_rss{" + "image!=\"\", " +
		"pod=" + "\""+ c.PodName +"\", namespace=\"" + c.Namespace +
		"\"})"

	result := c.getUtilizationQuery(query)
	return result
}

func (c *PodStatistics) GetAvgDiskUtilizationQuery() []interface{} {
	query := "disk_utilization{pod=" + "\"" + c.PodName +  "\", namespace=\"" + c.Namespace + "\"}"

	result := c.getUtilizationQuery(query)
	return result
}

func (c *PodStatistics) GetLastCpuUtilizationQuery() float64 {
	cpuTimeAndUtilizations := c.GetAvgCpuUtilizationQuery()

	utilization := cpuTimeAndUtilizations[len(cpuTimeAndUtilizations) - 1].(float64)
	return utilization
}

func (c *PodStatistics) GetLastMemoryUtilizationQuery() int64 {
	memoryTimeAndUtilizations := c.GetAvgMemoryUtilizationQuery()

	utilization := memoryTimeAndUtilizations[len(memoryTimeAndUtilizations) - 1].(int64)
	return utilization
}

func (c *PodStatistics) GetLastDiskUtilizationQuery() float64 {
	diskTimeAndUtilizations := c.GetAvgDiskUtilizationQuery()

	utilization := diskTimeAndUtilizations[len(diskTimeAndUtilizations) - 1].(float64)
	return utilization
}

