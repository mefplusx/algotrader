package interfaces

type Chart interface {
	CreateChart(prices []float64, imgPath string)
}
