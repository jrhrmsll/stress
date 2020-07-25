package main

import (
	"context"
	"flag"
	"log"
	"os"
	"sync"
	"syscall"

	"github.com/oklog/run"

	"stress/internal"
)

var (
	throughput int
	duration   string
	url        string
)

func main() {
	flag.IntVar(&throughput, "t", 1, "Throughput (requests per second).")
	flag.StringVar(&duration, "d", "60s", "Test duration.")
	flag.StringVar(&url, "u", "http://localhost/", "URL of the system under test.")

	flag.Parse()

	ctx := context.Background()

	logger := log.New(os.Stdout, "", log.LstdFlags)

	app := run.Group{}
	app.Add(run.SignalHandler(ctx, os.Interrupt, os.Kill, syscall.SIGTERM))

	// Channels
	requests := make(chan *internal.Request, throughput)
	responses := make(chan *internal.Response, throughput)

	producer := internal.NewProducer(throughput, duration, url, requests, logger)
	app.Add(producer.Execute, producer.Interrupt)

	stop := make(chan struct{})
	var wg sync.WaitGroup

	for i := 0; i < throughput; i++ {
		wg.Add(1)

		consumer := internal.NewConsumer(requests, responses, &wg, stop, logger)
		app.Add(consumer.Execute, consumer.Interrupt)
	}

	sink := internal.NewSink(throughput, duration, url, responses, logger)
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
