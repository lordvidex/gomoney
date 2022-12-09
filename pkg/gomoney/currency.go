package gomoney

type Currency string

const (
	USD Currency = "USD"
	RUB Currency = "RUB"
	NGN Currency = "NGN"
)

func (c Currency) IsValid() bool {
	switch c {
	case NGN, RUB, USD:
		return true
	default:
		return false
	}
}
