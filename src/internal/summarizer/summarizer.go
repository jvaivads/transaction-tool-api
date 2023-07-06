package summarizer

import "time"

type summarizer struct {
}

func (s summarizer) getBalance(txns transactions) float64 {
	var balance float64
	for _, txn := range txns {
		balance += txn.amount
	}
	return balance
}

func (s summarizer) getDebitAvg(txns transactions) float64 {
	var (
		avg   float64
		count float64
	)

	for _, txn := range txns {
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

	for _, txn := range txns {
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
	for _, txn := range txns {
		totalTransactionsByMonth[txn.date.Month()]++
	}
	return totalTransactionsByMonth
}
