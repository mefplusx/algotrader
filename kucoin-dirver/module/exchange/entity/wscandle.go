package entity

import (
	"strconv"

	"github.com/Kucoin/kucoin-go-sdk"
)

// time     : 0
// open     : 1
// close    : 2
// high     : 3
// low      : 4
// volume   : 5
// turnover : 6
type WSCandle struct {
	// Symbol string            `json:"symbol"`
	KLine kucoin.KLineModel `json:"candles"`
	// Time   int64             `json:"time"`
}

func (wsc *WSCandle) GetTime() int64 {
	if wsc.KLine == nil || len(wsc.KLine) != 7 {
		return 0
	}

	time, _ := strconv.ParseInt((wsc.KLine)[0], 10, 64)
	return time
}

func (wsc *WSCandle) GetVol() float64 {
	if wsc.KLine == nil || len(wsc.KLine) != 7 {
		return 0
	}

	volume, _ := strconv.ParseFloat((wsc.KLine)[5], 64)
	return volume
}

func (wsc *WSCandle) GetClose() float64 {
	if wsc.KLine == nil || len(wsc.KLine) != 7 {
		return 0
	}

	close, _ := strconv.ParseFloat((wsc.KLine)[2], 64)
	return close
}
