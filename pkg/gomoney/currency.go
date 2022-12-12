package gomoney

type Currency string

const (
	USD Currency = "USD"
	RUB          = "RUB"
	NGN          = "NGN"
)

func (c Currency) IsValid() bool {
	switch c {
	case NGN, RUB, USD:
		return true
	default:
		return false
	}
}
