package internal

import "time"

type Response struct {
	Number int64

	Start   time.Time
	End     time.Time
	Latency time.Duration

	Error error
}
