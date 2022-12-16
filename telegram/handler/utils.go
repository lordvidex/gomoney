package handler

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/lordvidex/gomoney/pkg/gomoney"
)

func BeautifulUserData(user *gomoney.User) string {
	return tablify([]string{"Name", "Phone"}, [][]string{{user.Name, user.Phone}})
}

func BeautifulTransferSummary(summaries []gomoney.TransactionSummary) string {
	var buf bytes.Buffer
	for i := 0; i < len(summaries); i++ {
		buf.WriteString(fmt.Sprintf("`Account: %s\nBalance: %.2f %s\n`", summaries[i].Account.Title, summaries[i].Account.Balance, summaries[i].Account.Currency))
		table := BeautifulTransactions(summaries[i].Transactions)
		buf.WriteString(table + "\n")
	}
	return buf.String()
}

func BeautifulTransactions(tx []gomoney.Transaction) string {
	return tablify([]string{"From", "To", "Type", "Amount", "Date"}, func() [][]string {
		m := make([][]string, len(tx))
		fc := func(x *gomoney.Account) string {
			if x == nil {
				return ""
			}
			return x.Title
		}
		for i := 0; i < len(tx); i++ {
			m[i] = []string{
				fc(tx[i].From),
				fc(tx[i].To),
				emoji(tx[i].Type),
				fmt.Sprintf("%.2f", tx[i].Amount),
				tx[i].Created.Format("15:04, 02 Jan"),
			}
		}
		return m
	}())
}

func emoji(t gomoney.TransactionType) string {
	switch t {
	case gomoney.Deposit:
		return "⬇️"
	case gomoney.Withdrawal:
		return "⬆️"
	case gomoney.Transfer:
		return "↔️"
	default:
		return ""
	}
}

func BeautifulAccounts(acc []gomoney.Account) string {
	return tablify([]string{"Title", "Balance", "Currency"}, func() [][]string {
		arr := make([][]string, len(acc))
		for i := 0; i < len(acc); i++ {
			arr[i] = []string{acc[i].Title, fmt.Sprintf("%.2f", acc[i].Balance), string(acc[i].Currency)}
		}
		return arr
	}())
}

func tablify(header []string, data [][]string) string {
	var buf bytes.Buffer
	lens := make([]int, len(header))
	for i := 0; i < len(header); i++ {
		lens[i] = len([]rune(header[i])) + 2
	}
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[i]); j++ {
			rl := len([]rune(data[i][j])) + 2
			if lens[j] < rl {
				lens[j] = rl
			}
		}
	}

	all := func() int {
		sum := 1
		for i := 0; i < len(lens); i++ {
			sum += lens[i]
		}
		return sum
	}()

	divide := func(l int) {
		buf.WriteString("+")
		buf.WriteString(strings.Repeat("-", l-2))
		buf.WriteString("+\n")
	}

	buf.WriteString("```")
	divide(all)

	for i := 0; i < len(header); i++ {
		buf.WriteString(fmt.Sprintf("| %-*s", lens[i]-2, header[i]))
	}
	buf.WriteString("|\n")
	divide(all)
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[i]); j++ {
			buf.WriteString(fmt.Sprintf("| %-*s", lens[j]-2, data[i][j]))
		}
		buf.WriteString("|\n\n")
	}
	divide(all)
	buf.WriteString("```\n")
	return buf.String()
}
