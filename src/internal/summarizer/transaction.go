package summarizer

import "time"

type transactions struct {
	items  []transaction
	userID int64
}

type transaction struct {
	amount float64
	date   time.Time
}
