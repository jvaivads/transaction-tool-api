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
			transactions: transactions{},
			expected:     0,
		},
		{
			name: "only credit transactions",
			transactions: transactions{
				items: []transaction{
					{amount: 10},
					{amount: 20},
					{amount: 30.5},
				},
			},
			expected: 60.5,
		},
		{
			name: "only debit transactions",
			transactions: transactions{
				items: []transaction{
					{amount: -10},
					{amount: -20},
					{amount: -30.5},
				},
			},
			expected: -60.5,
		},
		{
			name: "debit and credit transactions",
			transactions: transactions{
				items: []transaction{
					{amount: -10},
					{amount: 20},
					{amount: -30.5},
				},
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
				items: []transaction{
					{amount: 10},
					{amount: 20},
					{amount: 30.5},
				},
			},
			expected: 0,
		},
		{
			name: "only debit transactions",
			transactions: transactions{
				items: []transaction{
					{amount: -10},
					{amount: -20},
					{amount: -60},
				},
			},
			expected: -30,
		},
		{
			name: "debit and credit transactions",
			transactions: transactions{
				items: []transaction{
					{amount: -10},
					{amount: 20},
					{amount: -60},
				},
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
				items: []transaction{
					{amount: -10},
					{amount: -20},
					{amount: -30.5},
				},
			},
			expected: 0,
		},
		{
			name: "only credit transactions",
			transactions: transactions{
				items: []transaction{
					{amount: 10},
					{amount: 20},
					{amount: 60},
				},
			},
			expected: 30,
		},
		{
			name: "debit and credit transactions",
			transactions: transactions{
				items: []transaction{
					{amount: 10},
					{amount: -20},
					{amount: 60},
				},
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
			transactions: transactions{},
			expected:     map[time.Month]int{},
		},
		{
			name: "only transactions in january",
			transactions: transactions{
				items: []transaction{
					{amount: 10, date: time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)},
					{amount: 20, date: time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)},
					{amount: 30, date: time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)},
				},
			},
			expected: map[time.Month]int{
				time.January: 3,
			},
		},
		{
			name: "different months",
			transactions: transactions{
				items: []transaction{
					{amount: -10, date: time.Date(2021, time.December, 1, 0, 0, 0, 0, time.UTC)},
					{amount: 20, date: time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)},
					{amount: -30, date: time.Date(2021, time.April, 1, 0, 0, 0, 0, time.UTC)},
				},
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

func TestSummarizerResume(t *testing.T) {
	user := User{
		UserID: 1,
		Name:   "name",
		Email:  "email",
	}

	tests := []struct {
		name         string
		transactions transactions
		expected     Resume
	}{
		{
			name: "resume",
			transactions: transactions{
				items: []transaction{
					{amount: -10, date: time.Date(2021, time.December, 1, 0, 0, 0, 0, time.UTC)},
					{amount: 15, date: time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)},
					{amount: -60, date: time.Date(2021, time.April, 1, 0, 0, 0, 0, time.UTC)},
				},
			},
			expected: Resume{
				User:      user,
				Balance:   "-55.00",
				CreditAvg: "15.00",
				DebitAvg:  "-35.00",

				MonthTransactions: []MonthTransaction{
					{time.January.String(), 1},
					{time.April.String(), 1},
					{time.December.String(), 1},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, summarizer{}.resume(user, test.transactions))
		})
	}
}

func TestResumeToHTML(t *testing.T) {
	resume := Resume{
		User:      User{Name: "name"},
		Balance:   "1.12",
		CreditAvg: "1.4",
		DebitAvg:  "1.55",
		MonthTransactions: []MonthTransaction{
			{
				Month:             time.January.String(),
				TotalTransactions: 5,
			},
			{
				Month:             time.December.String(),
				TotalTransactions: 5,
			},
		},
	}

	tests := []struct {
		name        string
		tmpl        string
		resume      Resume
		result      string
		expectedErr bool
	}{
		{
			name:        "error parsing",
			resume:      resume,
			tmpl:        "<h1>Balance: {{.Bad}}</h1>",
			expectedErr: true,
		},
		{
			name:   "error parsing",
			resume: resume,
			tmpl:   resumeHTMLTemplate,
			result: "<body>    <img src=\"https://blog.storicard.com/wp-content/uploads/2019/07/Stori-horizontal-11.jpg\"> " +
				"   <p>Hello name</p>    <h1>Balance: 1.12</h1>    <h1>Credit Average: 1.4</h1>   " +
				" <h1>Debit Average: 1.55</h1>    <h1>Transactions by month</h1>    " +
				"<p>         <p>January: 5</p> <p>December: 5</p>    </p></body>",
			expectedErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := test.resume.ToHTML(test.tmpl)
			assert.Equal(t, test.result, result)
			assert.Equal(t, test.expectedErr, err != nil)
		})
	}
}
