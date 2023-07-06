package summarizer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSummarizerGetBalance(t *testing.T) {
	tests := []struct {
		name         string
		transactions transactions
		expected     float64
	}{
		{
			name:         "without transactions return 0",
			transactions: nil,
			expected:     0,
		},
		{
			name: "only credit transactions",
			transactions: transactions{
				{amount: 10},
				{amount: 20},
				{amount: 30.5},
			},
			expected: 60.5,
		},
		{
			name: "only debit transactions",
			transactions: transactions{
				{amount: -10},
				{amount: -20},
				{amount: -30.5},
			},
			expected: -60.5,
		},
		{
			name: "debit and credit transactions",
			transactions: transactions{
				{amount: -10},
				{amount: 20},
				{amount: -30.5},
			},
			expected: -20.5,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, summarizer{}.getBalance(test.transactions))

		})
	}
}

func TestSummarizerGetDebitAvg(t *testing.T) {
	tests := []struct {
		name         string
		transactions transactions
		expected     float64
	}{
		{
			name: "without debit transactions return 0",
			transactions: transactions{
				{amount: 10},
				{amount: 20},
				{amount: 30.5},
			},
			expected: 0,
		},
		{
			name: "only debit transactions",
			transactions: transactions{
				{amount: -10},
				{amount: -20},
				{amount: -60},
			},
			expected: -30,
		},
		{
			name: "debit and credit transactions",
			transactions: transactions{
				{amount: -10},
				{amount: 20},
				{amount: -60},
			},
			expected: -35,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, summarizer{}.getDebitAvg(test.transactions))

		})
	}
}

func TestSummarizerGetCreditAvg(t *testing.T) {
	tests := []struct {
		name         string
		transactions transactions
		expected     float64
	}{
		{
			name: "without credit transactions return 0",
			transactions: transactions{
				{amount: -10},
				{amount: -20},
				{amount: -30.5},
			},
			expected: 0,
		},
		{
			name: "only credit transactions",
			transactions: transactions{
				{amount: 10},
				{amount: 20},
				{amount: 60},
			},
			expected: 30,
		},
		{
			name: "debit and credit transactions",
			transactions: transactions{
				{amount: 10},
				{amount: -20},
				{amount: 60},
			},
			expected: 35,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, summarizer{}.getCreditAvg(test.transactions))

		})
	}
}

func TestSummarizerGetTotalTransactionsByMonth(t *testing.T) {
	tests := []struct {
		name         string
		transactions transactions
		expected     map[time.Month]int
	}{
		{
			name:         "no transactions",
			transactions: nil,
			expected:     map[time.Month]int{},
		},
		{
			name: "only transactions in january",
			transactions: transactions{
				{amount: 10, date: time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)},
				{amount: 20, date: time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)},
				{amount: 30, date: time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)},
			},
			expected: map[time.Month]int{
				time.January: 3,
			},
		},
		{
			name: "different months",
			transactions: transactions{
				{amount: -10, date: time.Date(2021, time.December, 1, 0, 0, 0, 0, time.UTC)},
				{amount: 20, date: time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)},
				{amount: -30, date: time.Date(2021, time.April, 1, 0, 0, 0, 0, time.UTC)},
			},
			expected: map[time.Month]int{
				time.December: 1,
				time.January:  1,
				time.April:    1,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, summarizer{}.getTotalTransactionsByMonth(test.transactions))
		})
	}
}
