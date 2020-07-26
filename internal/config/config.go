package config

import "time"

type Config struct {
	throughput int
	duration   string
	timeout    string
	url        string

	interval time.Duration
	requests int64
}

func NewConfig(
	throughput int,
	duration string,
	timeout string,
	url string,
) *Config {
	// intended interval between requests, with a constant
	// throughput of N requests per second (1s = 1e9ns)
	interval := time.Duration(1e9 / throughput)

	// total number of requests
	d, err := time.ParseDuration(duration)
	if err != nil {
		panic(err)
	}

	requests := int64(throughput) * int64(d.Seconds())

	return &Config{
		throughput: throughput,
		duration:   duration,
		timeout:    timeout,
		url:        url,

		interval: interval,
		requests: requests,
	}
}

func (cfg *Config) Throughput() int {
	return cfg.throughput
}

func (cfg *Config) Duration() string {
	return cfg.duration
}

func (cfg *Config) Timeout() string {
	return cfg.timeout
}

func (cfg *Config) Url() string {
	return cfg.url
}

func (cfg *Config) Interval() time.Duration {
	return cfg.interval
}

func (cfg *Config) IntervalAsInt64() int64 {
	return cfg.interval.Milliseconds()
}

func (cfg *Config) Requests() int64 {
	return cfg.requests
}
