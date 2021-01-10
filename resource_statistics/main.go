package main

import (
	"flag"
	"fmt"
	"github.com/idoubi/goz"
	"log"
)

var (
	endpoint        string   /* the prometheus http://url:port*/
	statefulSetName string
	namespace       string
	interval        int
)


func init() {
	flag.StringVar(&endpoint, "endpoint", "http://127.0.0.1:9090", "promethues url")
	flag.StringVar(&statefulSetName, "statefulset", "default", "statefulset Name")
	flag.StringVar(&namespace, "namespace", "default", "namespace")
	flag.IntVar(&interval, "interval", 5, "interval time")
}


func main() {
	curl := PromCurl{endpoint, namespace, nil}

	responseBody, err := curl.Get("/api/v1/query", goz.Options{
		Query: map[string]interface{}{
			"query": "container_cpu_usage_seconds_total",
			"time": 1610078841,
		},
	})
	if err != nil {
		log.Fatal("curl.Get error")
	}
	contents := responseBody.GetContents()


}

