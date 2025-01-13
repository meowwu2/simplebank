package util

const(
	USD="USD"
	EUR="EUR"
	RMB="RMB"
)

func IsSupportCurrency(currency string) bool{
	switch currency{
	case EUR,USD,RMB:
		return true
	}
	return false
}