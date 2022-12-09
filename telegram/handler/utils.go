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
		lens[i] = len(header[i]) + 2
	}
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[i]); j++ {
			if lens[j]-2 < len(data[i][j]) {
				lens[j] = len(data[i][j]) + 2
			}
		}
	}

	all := func() int {
		sum := -1
		for i := 0; i < len(lens); i++ {
			sum += lens[i] + 1
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
		buf.WriteString("|\n")
	}
	divide(all)
	buf.WriteString("```")
	return buf.String()
}
