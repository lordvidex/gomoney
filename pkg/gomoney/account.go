package gomoney

type Account struct {
	Id          int64
	Title       string
	Description string
	Balance     float64
	Currency    Currency
	IsBlocked   bool
}

func (a *Account) CanTransfer(amount float64) bool {
	return a.Balance >= amount && !a.IsBlocked
}
