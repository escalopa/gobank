package util

const (
	EGP = "EGP"
	USD = "USD"
	RUB = "RUB"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case EGP, USD, RUB:
		return true
	}
	return false
}
