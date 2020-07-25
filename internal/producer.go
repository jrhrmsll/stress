package internal

import (
	"fmt"
	"log"
	"time"
)

const header = `
URL:        %s
Throughput: %d request/s
Interval:   %s
Requests:   %d

`

type Producer struct {
	throughput int
	duration   string
	url        string

	requests chan<- *Request

	logger *log.Logger

	interval time.Duration
	total    int64
	quit     chan struct{}
}

func NewProducer(
	throughput int,
	duration string,
	url string,
	requests chan<- *Request,
	logger *log.Logger,
) *Producer {
	// intended interval between requests, with a constant
	// throughput of N requests per second (1s = 1e9ns)
	interval := time.Duration(1e9 / throughput)

	// total number of requests
	d, err := time.ParseDuration(duration)
	if err != nil {
		panic(err)
	}

	total := int64(throughput) * int64(d.Seconds())

	return &Producer{
		throughput: throughput,
		duration:   duration,
		url:        url,

		requests: requests,

		logger: logger,

		interval: interval,
		total:    total,
		quit:     make(chan struct{}),
	}
}

func (p *Producer) Execute() error {
	p.init()
	defer close(p.requests)

	for i := int64(0); i < p.total; i++ {
		request := &Request{
			Number: i + 1,
			URL:    p.url,
		}

		select {
		case <-p.quit:
			return nil
		case p.requests <- request:
		}

		time.Sleep(p.interval)
	}

	return nil
}

func (p *Producer) init() {
	fmt.Printf(
		header,
		p.url,
		p.throughput,
		p.interval,
		p.total,
	)
}

func (p *Producer) Interrupt(err error) {
	close(p.quit)
}
