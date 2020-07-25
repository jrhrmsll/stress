package internal

import (
	"fmt"
	"log"
	"strconv"
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

	percentiles := []string{"50", "75", "95", "99", "99.9", "99.95", "99.99"}
	fmt.Println(" Response times percentiles:")

	for _, p := range percentiles {
		v, _ := strconv.ParseFloat(p, 64)
		fmt.Printf("%9sth: %6dms\n", p, hist.ValueAtQuantile(v))
	}

	fmt.Println()

	return nil
}

func (s *Sink) Interrupt(err error) {}
