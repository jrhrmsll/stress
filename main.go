package main

import (
	"context"
	"flag"
	"log"
	"os"
	"sync"
	"syscall"

	"github.com/oklog/run"

	"github.com/jrhrmsll/stress/internal"
	"github.com/jrhrmsll/stress/internal/config"
	"github.com/jrhrmsll/stress/internal/pipeline"
	"github.com/jrhrmsll/stress/internal/processor"
)

func main() {
	var (
		throughput int
		duration   string
		timeout    string
		url        string
	)

	flag.IntVar(&throughput, "throughput", 1, "Throughput (requests per second).")
	flag.StringVar(&duration, "duration", "60s", "Test duration.")
	flag.StringVar(&timeout, "timeout", "30s", "Request timeout.")
	flag.StringVar(&url, "url", "http://localhost/", "URL of the system under test.")

	flag.Parse()

	cfg := config.NewConfig(throughput, duration, timeout, url)

	ctx := context.Background()

	logger := log.New(os.Stdout, "", log.LstdFlags)

	app := run.Group{}
	app.Add(run.SignalHandler(ctx, os.Interrupt, os.Kill, syscall.SIGTERM))

	// Channels
	requests := make(chan *internal.Request, throughput)
	responses := make(chan *internal.Response, throughput)

	producer := pipeline.NewProducer(cfg, requests, logger)
	app.Add(producer.Execute, producer.Interrupt)

	var wg sync.WaitGroup
	for i := 0; i < throughput; i++ {
		wg.Add(1)

		consumer := pipeline.NewConsumer(cfg, requests, responses, &wg, logger)
		app.Add(consumer.Execute, consumer.Interrupt)
	}

	processors := make([]processor.ResultsProcessor, 0)
	processors = append(
		processors,
		processor.NewHdrHistDump(cfg, os.Stdout),
		processor.NewCSVFileDump(cfg, logger),
	)

	sink := pipeline.NewSink(cfg, responses, processors, logger)
	app.Add(sink.Execute, sink.Interrupt)

	go func() {
		wg.Wait()
		close(responses)
	}()

	err := app.Run()
	if err != nil {
		logger.Println(err)
	}
}
