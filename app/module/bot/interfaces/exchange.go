package interfaces

type Exchange interface {
	BuyMarket(typo string) bool
	SetStop(typo string, price float64)
}
