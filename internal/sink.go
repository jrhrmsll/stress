package internal

import (
	"fmt"
	"log"
	"time"

	"github.com/codahale/hdrhistogram"
)

type Sink struct {
	throughput int
	duration   string
	url        string

	responses <-chan *Response
	logger    *log.Logger

	interval time.Duration
}

func NewSink(
	throughput int,
	duration string,
	url string,
	responses <-chan *Response,
	logger *log.Logger,
) *Sink {
	// intended interval between requests, with a constant
	// throughput of N requests per second (1s = 1e9ns)
	interval := time.Duration(1e9 / throughput)

	return &Sink{
		throughput: throughput,
		duration:   duration,
		url:        url,

		responses: responses,
		logger:    logger,

		interval: interval,
	}
}

func (s *Sink) Execute() error {
	hist := hdrhistogram.New(0, 1000000, 3)
	for response := range s.responses {
		response.Latency = response.End.Sub(response.Start)

		hist.RecordCorrectedValue(response.Latency.Milliseconds(), s.interval.Milliseconds())
	}

	fmt.Println(hist.CumulativeDistribution())

	return nil
}

func (s *Sink) Interrupt(err error) {}
