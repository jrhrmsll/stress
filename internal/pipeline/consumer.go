package pipeline

import (
	"log"
	"net/http"
	"sync"
	"time"

	"stress/internal"
	"stress/internal/config"
)

const (
	HeaderContentType   = "Content-Type"
	MIMEApplicationJSON = "application/json"
)

type Consumer struct {
	cfg       *config.Config
	requests  <-chan *internal.Request
	responses chan<- *internal.Response
	wg        *sync.WaitGroup
	logger    *log.Logger
}

func NewConsumer(
	cfg *config.Config,
	requests <-chan *internal.Request,
	responses chan<- *internal.Response,
	wg *sync.WaitGroup,
	logger *log.Logger,
) *Consumer {
	return &Consumer{
		cfg:       cfg,
		requests:  requests,
		responses: responses,
		wg:        wg,
		logger:    logger,
	}
}

func (c *Consumer) Execute() error {
	for request := range c.requests {
		response, err := c.execute(request)
		if err != nil {
			c.logger.Println(err)
		}

		c.responses <- response
	}

	c.wg.Done()

	return nil
}

func (c *Consumer) Interrupt(err error) {}

func (c *Consumer) execute(request *internal.Request) (*internal.Response, error) {
	start := time.Now().UTC()

	timeout, err := time.ParseDuration(c.cfg.Timeout())
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

	response := &internal.Response{
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
