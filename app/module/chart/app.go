package chart

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

const timeScale = 20
const priceScale = 1

type App struct {
}

func (a *App) CreateChart(prices []float64, imgPath string) {
	var minPrice float64 = 999999999
	var maxPrice float64 = 0

	for _, price := range prices {
		if maxPrice < price {
			maxPrice = price
		}
		if minPrice > price {
			minPrice = price
		}
	}

	width := int(float64(len(prices))) * timeScale
	heigth := int(float64(maxPrice-minPrice)) * priceScale

	img := createImage(width, heigth)
	drawPrices(img, prices)
	saveImage(img, imgPath)
}

func createImage(width, heigth int) *image.RGBA {
	c := color.RGBA{0, 0, 0, 255}
	var img = image.NewRGBA(image.Rect(0, 0, width, heigth))
	for w := 0; w < width; w++ {
		for h := 0; h < heigth; h++ {
			img.Set(w, h, c)
		}
	}

	return img
}

func fillRect(img *image.RGBA, x, y, x2, y2 int, c *color.RGBA) {
	for _x := x; _x < x2; _x++ {
		for _y := y; _y < y2; _y++ {
			img.Set(_x, _y, *c)
		}
	}
}

func saveImage(img *image.RGBA, path string) {
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	png.Encode(f, img)
	f.Close()
}

func point(img *image.RGBA, x, y, size int, c *color.RGBA) {
	for _x := x; _x < x+size; _x++ {
		for _y := y; _y < y+size; _y++ {
			img.Set(_x, _y, c)
		}
	}
}

func drawPrices(img *image.RGBA, prices []float64) {
	const strokeLine = 10

	var minPrice float64 = 999999999
	var maxPrice float64 = 0

	for _, price := range prices {
		if maxPrice < price {
			maxPrice = price
		}
		if minPrice > price {
			minPrice = price
		}
	}

	var bef float64 = 0
	for timeSource, priceSource := range prices {
		time := int(float64(timeSource) * timeScale)
		price := maxPrice - (float64(float64(priceSource)*priceScale) - minPrice)

		if priceSource > bef {
			point(img, time, int(price), strokeLine, &color.RGBA{0, 255, 0, 255})
		} else {
			point(img, time, int(price), strokeLine, &color.RGBA{255, 0, 0, 255})
		}

		bef = priceSource
	}
}
