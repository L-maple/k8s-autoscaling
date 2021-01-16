package statistics

import (
	"github.com/idoubi/goz"
	"log"
)


// PromCurl: a struct used for Get prometheus http request
type PromCurl struct {
	// endpoint is url:port,
	// namespace is the statefulSet/deployment's namespace
	endpoint, namespace  string

	// httpClient is the singleton client
	httpClient           *goz.Request
}

func (c *PromCurl) Get(path string, queryParams goz.Options) (goz.ResponseBody, error) {
	if c.httpClient == nil || queryParams.Query != nil  {
		c.httpClient = goz.NewClient(queryParams)
	}

	url := c.endpoint + path

	resp, err := c.httpClient.Get(url, queryParams)
	if err != nil {
		log.Fatal("c.httpClient.Get error: ", err)
	}

	return resp.GetBody()
}


type PromJsonParser struct {
	
}

type PromRespContent struct {
	Data    interface{}
	Status  string
}

func (p *PromJsonParser) parser() {

}
