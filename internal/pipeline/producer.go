package pipeline

import (
	"fmt"
	"log"
	"time"

	"github.com/jrhrmsll/stress/internal"
	"github.com/jrhrmsll/stress/internal/config"
)

const header = `
        URL: %s
 Throughput: %d request/s
   Interval: %s
    Timeout: %s
   Requests: %d

`

type Producer struct {
	cfg      *config.Config
	requests chan<- *internal.Request
	logger   *log.Logger

	quit chan struct{}
}

func NewProducer(
	cfg *config.Config,
	requests chan<- *internal.Request,
	logger *log.Logger,
) *Producer {

	return &Producer{
		cfg:      cfg,
		requests: requests,
		logger:   logger,

		quit: make(chan struct{}),
	}
}

func (producer *Producer) Execute() error {
	producer.init()
	defer close(producer.requests)

	for i := int64(0); i < producer.cfg.Requests(); i++ {
		request := &internal.Request{
			Number: i + 1,
			URL:    producer.cfg.Url(),
		}

		select {
		case <-producer.quit:
			return nil
		case producer.requests <- request:
		}

		time.Sleep(producer.cfg.Interval())
	}

	return nil
}

func (producer *Producer) init() {
	fmt.Printf(
		header,
		producer.cfg.Url(),
		producer.cfg.Throughput(),
		producer.cfg.Interval(),
		producer.cfg.Timeout(),
		producer.cfg.Requests(),
	)
}

func (producer *Producer) Interrupt(err error) {
	close(producer.quit)
}
