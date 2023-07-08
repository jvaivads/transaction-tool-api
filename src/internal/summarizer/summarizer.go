package summarizer

import (
	"encoding/json"
	"time"
)

type summarizer struct {
}

func (s summarizer) getBalance(txns transactions) float64 {
	var balance float64
	for _, txn := range txns.items {
		balance += txn.amount
	}
	return balance
}

func (s summarizer) getDebitAvg(txns transactions) float64 {
	var (
		avg   float64
		count float64
	)

	for _, txn := range txns.items {
		if txn.amount < 0 {
			avg += txn.amount
			count++
		}
	}

	if count > 0 {
		return avg / count
	}
	return 0
}

func (s summarizer) getCreditAvg(txns transactions) float64 {
	var (
		avg   float64
		count float64
	)

	for _, txn := range txns.items {
		if txn.amount > 0 {
			avg += txn.amount
			count++
		}
	}

	if count > 0 {
		return avg / count
	}
	return 0
}

func (s summarizer) getTotalTransactionsByMonth(txns transactions) map[time.Month]int {
	totalTransactionsByMonth := make(map[time.Month]int)
	for _, txn := range txns.items {
		totalTransactionsByMonth[txn.date.Month()]++
	}
	return totalTransactionsByMonth
}

func (s summarizer) resume(txns transactions) Resume {
	return Resume{
		Balance:                  s.getBalance(txns),
		CreditAvg:                s.getCreditAvg(txns),
		DebitAvg:                 s.getDebitAvg(txns),
		TotalTransactionsByMonth: s.getTotalTransactionsByMonth(txns),
	}
}

type Resume struct {
	Balance                  float64            `json:"balance"`
	CreditAvg                float64            `json:"credit_average"`
	DebitAvg                 float64            `json:"debit_average"`
	TotalTransactionsByMonth map[time.Month]int `json:"total_transactions_by_month"`
}

func (r Resume) String() string {
	b, _ := json.Marshal(r)
	return string(b)
}
