package statistics

import (
	"github.com/idoubi/goz"
	"log"
	"time"
)


// PromCurl: a struct used for Get prometheus http request
type PromCurl struct {
	// endpoint is url:port
	endpoint  string

	// httpClient is the singleton client
	httpClient           *goz.Request
}

func (c *PromCurl) Get(path string, queryParams goz.Options) (goz.ResponseBody, error) {
	if c.httpClient == nil || queryParams.Query != nil  {
		c.httpClient = goz.NewClient(queryParams)
	}

	url := c.endpoint + path

	waitTime := 1
	for {
		resp, err := c.httpClient.Get(url, queryParams)
		if err != nil {
			log.Println("c.httpClient.Get error: ", err)

			time.Sleep(time.Duration(waitTime) * time.Second)
			waitTime <<= 1
			continue
		}

		if waitTime > 512 {
			log.Fatal("The process failed: c.httpClient.Get error:(")
		}

		return resp.GetBody()
	}

}

