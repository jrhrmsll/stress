package internal

import (
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	HeaderContentType   = "Content-Type"
	MIMEApplicationJSON = "application/json"
	httpClientTimeout   = "30s"
)

type Consumer struct {
	requests  <-chan *Request
	responses chan<- *Response
	wg        *sync.WaitGroup
	logger    *log.Logger

	stop chan struct{}
}

func NewConsumer(
	requests <-chan *Request,
	responses chan<- *Response,
	wg *sync.WaitGroup,
	stop chan struct{},
	logger *log.Logger,
) *Consumer {
	return &Consumer{
		requests:  requests,
		responses: responses,
		wg:        wg,
		stop:      stop,
		logger:    logger,
	}
}

func (c *Consumer) Execute() error {
	for request := range c.requests {
		response, err := execute(request)
		if err != nil {
			c.logger.Println(err)
		}

		c.responses <- response
	}

	c.wg.Done()

	return nil
}

func (c *Consumer) Interrupt(err error) {}

func execute(request *Request) (*Response, error) {
	start := time.Now().UTC()

	timeout, err := time.ParseDuration(httpClientTimeout)
	if err != nil {
		return nil, err
	}

	httpClient := http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest(http.MethodGet, request.URL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set(HeaderContentType, MIMEApplicationJSON)

	resp, err := httpClient.Do(req)

	response := &Response{
		Number: request.Number,
		Start:  start,
		End:    time.Now().UTC(),
		Error:  err,
	}

	if err != nil {
		return response, err
	}

	defer resp.Body.Close()

	return response, nil
}
