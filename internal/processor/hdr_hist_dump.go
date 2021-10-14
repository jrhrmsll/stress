package processor

import (
	"fmt"
	"io"
	"strconv"

	"github.com/HdrHistogram/hdrhistogram-go"

	"github.com/jrhrmsll/stress/internal/config"
)

type HdrHistDump struct {
	cfg       *config.Config
	writer    io.Writer
	histogram *hdrhistogram.Histogram
}

func NewHdrHistDump(cfg *config.Config, writer io.Writer) *HdrHistDump {
	return &HdrHistDump{
		cfg:       cfg,
		writer:    writer,
		histogram: hdrhistogram.New(0, 1000000, 3),
	}
}

func (processor *HdrHistDump) Process(results map[int64]int64) error {
	for k, v := range results {
		for i := int64(0); i < v; i++ {
			processor.histogram.RecordCorrectedValue(k, processor.cfg.IntervalAsInt64())
		}
	}

	percentiles := []string{"50", "75", "95", "99", "99.9", "99.95", "99.99"}

	fmt.Fprintln(processor.writer, " Response times percentiles:")

	for _, p := range percentiles {
		v, _ := strconv.ParseFloat(p, 64)
		fmt.Fprintf(processor.writer, "%9sth: %6dms\n", p, processor.histogram.ValueAtQuantile(v))
	}

	fmt.Fprintln(processor.writer)

	return nil
}
