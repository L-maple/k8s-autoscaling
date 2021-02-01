package statistics

import (
	"github.com/idoubi/goz"
	"k8s.io/apimachinery/pkg/util/json"
	"log"
	"strconv"
	"time"
)

type PodStatistics struct {
	Endpoint           string   /* pod statistics' data source: http://ip:port */
	PodName, Namespace string
}

func (c *PodStatistics) getUsageQuery(query string, timeRange int64) []interface{}  {
	curl := PromCurl{c.Endpoint, nil}
	responseBody, err := curl.Get("/api/v1/query_range", goz.Options{
		Query: map[string]interface{}{
			"query": query,
			"start": strconv.FormatInt(time.Now().Unix() - timeRange, 10),
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

func (c *PodStatistics) GetAvgCpuUtilizations(timeRange int64) []interface{} {
	query := "sum(rate(container_cpu_usage_seconds_total{image!=\"\"," +
			 "pod=" + "\"" + c.PodName + "\", namespace=\"" + c.Namespace +
			 "\"}[1m])) / " +
			 "sum(container_spec_cpu_quota{image!=\"\", pod=\"" + c.PodName + "\", namespace=\"" + c.Namespace + "\"} " +
				"/ container_spec_cpu_period{image!=\"\", pod=\"" + c.PodName + "\", namespace=\"" +c.Namespace + "\"})"
	result := c.getUsageQuery(query, timeRange)
	return result
}

func (c *PodStatistics) GetAvgMemoryUsages(timeRange int64) []interface{} {
	query := "sum(container_memory_rss{" + "image!=\"\", " +
		"pod=" + "\""+ c.PodName +"\", namespace=\"" + c.Namespace +
		"\"})"

	result := c.getUsageQuery(query, timeRange)
	return result
}

func (c *PodStatistics) GetAvgDiskUtilizations(timeRange int64) []interface{} {
	query := "disk_utilization{pod=" + "\"" + c.PodName +  "\", namespace=\"" + c.Namespace + "\"}"

	result := c.getUsageQuery(query, timeRange)
	return result
}

func (c *PodStatistics) GetLastCpuUtilization() float64 {
	cpuTimeAndUtilizations := c.GetAvgCpuUtilizations(3600)

	lastCpuTimeAndUtilization := cpuTimeAndUtilizations[len(cpuTimeAndUtilizations) - 1]
	utilization := lastCpuTimeAndUtilization.([]interface{})[1].(string)

	floatVal, err := strconv.ParseFloat(utilization, 64)
	if err != nil {
		log.Fatal("strconv.ParseFloat: ", err)
	}

	return floatVal
}

func (c *PodStatistics) GetLastMemoryUsage() int64 {
	memoryTimeAndUtilizations := c.GetAvgMemoryUsages(3600)

	lastMemoryTimeAndUtilization := memoryTimeAndUtilizations[len(memoryTimeAndUtilizations) - 1]
	utilization := lastMemoryTimeAndUtilization.([]interface{})[1].(string)

	intVal, err := strconv.ParseInt(utilization, 10, 64)
	if err != nil {
		log.Fatal("strconv.ParseInt", err)
	}

	return intVal
}

func (c *PodStatistics) GetLastDiskUtilization() float64 {
	diskTimeAndUtilizations := c.GetAvgDiskUtilizations(3600)

	lastDiskTimeAndUtilization := diskTimeAndUtilizations[len(diskTimeAndUtilizations) - 1]
	utilization := lastDiskTimeAndUtilization.([]interface{})[1].(string)

	floatVal, err := strconv.ParseFloat(utilization, 64)
	if err != nil {
		log.Fatal("strconv.ParseFloat: ", err)
	}
	return floatVal
}

