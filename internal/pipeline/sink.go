package pipeline

import (
	"log"

	"stress/internal"
	"stress/internal/config"
	"stress/internal/processor"
)

type Sink struct {
	cfg *config.Config

	responses  <-chan *internal.Response
	processors []processor.ResultsProcessor
	logger     *log.Logger

	times map[int64]int64
}

func NewSink(
	cfg *config.Config,
	responses <-chan *internal.Response,
	processors []processor.ResultsProcessor,
	logger *log.Logger,
) *Sink {
	return &Sink{
		cfg: cfg,

		responses:  responses,
		processors: processors,
		logger:     logger,

		times: make(map[int64]int64),
	}
}

func (s *Sink) Execute() error {
	for response := range s.responses {
		response.Latency = response.End.Sub(response.Start)
		s.times[response.Latency.Milliseconds()]++
	}

	for _, p := range s.processors {
		if err := p.Process(s.times); err != nil {
			s.logger.Println(err)
		}
	}

	return nil
}

func (s *Sink) Interrupt(err error) {}
