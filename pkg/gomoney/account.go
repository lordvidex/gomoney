package gomoney

type Account struct {
	Id          int64
	Title       string
	Description string
	Balance     float64
	Currency    Currency
	IsBlocked   bool
}

func (a *Account) CanTransfer() bool {
	return a.Balance > 0 && !a.IsBlocked
}
