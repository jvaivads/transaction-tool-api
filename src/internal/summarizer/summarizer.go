package summarizer

import (
	"bytes"
	"fmt"
	"html/template"
	"sort"
	"strings"
	"time"
)

const (
	resumeHTMLTemplate = `
<body>
    <img src="https://blog.storicard.com/wp-content/uploads/2019/07/Stori-horizontal-11.jpg">
    <p>Hello {{.User.Name}}</p>
    <h1>Balance: {{.Balance}}</h1>
    <h1>Credit Average: {{.CreditAvg}}</h1>
    <h1>Debit Average: {{.DebitAvg}}</h1>
    <h1>Transactions by month</h1>
    <p>
        {{range .MonthTransactions}} <p>{{.Month}}: {{.TotalTransactions}}</p>{{end}}
    </p>
</body>`
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

func (s summarizer) resume(user User, txns transactions) Resume {
	var (
		months              []time.Month
		monthTransactions   []MonthTransaction
		transactionsByMonth = s.getTotalTransactionsByMonth(txns)
	)

	for month, _ := range transactionsByMonth {
		months = append(months, month)
	}

	sort.Slice(months, func(i, j int) bool {
		return months[i] < months[j]
	})

	for _, month := range months {
		monthTransactions = append(monthTransactions, MonthTransaction{
			Month:             month.String(),
			TotalTransactions: transactionsByMonth[month],
		})
	}

	return Resume{
		User:              user,
		Balance:           fmt.Sprintf("%.2f", s.getBalance(txns)),
		CreditAvg:         fmt.Sprintf("%.2f", s.getCreditAvg(txns)),
		DebitAvg:          fmt.Sprintf("%.2f", s.getDebitAvg(txns)),
		MonthTransactions: monthTransactions,
	}
}

type MonthTransaction struct {
	Month             string
	TotalTransactions int
}

type Resume struct {
	User              User
	Balance           string
	CreditAvg         string
	DebitAvg          string
	MonthTransactions []MonthTransaction
}

func (r Resume) ToHTML(tmpl string) (string, error) {
	tmpl = strings.ReplaceAll(tmpl, "\n", "")

	t, err := template.New("resume").Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("error creating resume template due to %w", err)
	}

	var buffer bytes.Buffer
	err = t.Execute(&buffer, r)
	if err != nil {
		return "", fmt.Errorf("error parsing resume due to %w", err)
	}

	return buffer.String(), nil
}
