package summarizer

import "time"

type transactions []transaction

type transaction struct {
	amount float64
	date   time.Time
}
