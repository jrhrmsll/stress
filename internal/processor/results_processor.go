package processor

type ResultsProcessor interface {
	Process(results map[int64]int64) error
}
