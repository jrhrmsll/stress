package processor

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"stress/internal/config"
)

const filename = "results.csv"

type CSVFileDump struct {
	cfg    *config.Config
	logger *log.Logger
}

func NewCSVFileDump(
	cfg *config.Config,
	logger *log.Logger,
) *CSVFileDump {
	return &CSVFileDump{
		cfg:    cfg,
		logger: logger,
	}
}

func (processor *CSVFileDump) Process(results map[int64]int64) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	defer file.Close()

	w := csv.NewWriter(file)

	err = w.Write(processor.meta())
	if err != nil {
		return err
	}

	for k, v := range results {
		if err := w.Write(processor.record(k, v)); err != nil {
			processor.logger.Println(err)
		}
	}

	w.Flush()

	if err := w.Error(); err != nil {
		return err
	}

	return nil
}

func (processor *CSVFileDump) meta() []string {
	return []string{
		strconv.FormatInt(int64(processor.cfg.Throughput()), 10),
		processor.cfg.Duration(),
		processor.cfg.Timeout(),
		processor.cfg.Url(),
		strconv.FormatInt(int64(processor.cfg.IntervalAsInt64()), 10),
		strconv.FormatInt(int64(processor.cfg.Requests()), 10),
	}
}

func (processor *CSVFileDump) record(k, v int64) []string {
	return []string{
		strconv.FormatInt(k, 10),
		strconv.FormatInt(v, 10),
	}
}
