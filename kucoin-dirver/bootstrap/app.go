package bootstrap

import (
	"kucoin-dirver/common/entity"
	"kucoin-dirver/module/exchange"
	"kucoin-dirver/module/hostconnector"
	"kucoin-dirver/module/httpapi"
	"log"
)

type Bootstrap struct {
	Config        *Config
	Exchange      *exchange.App
	Httpapi       *httpapi.App
	Hostconnector *hostconnector.App
}

func (a Bootstrap) Sync(currency string) {
	candleChan := make(chan entity.Candle)
	sessionChan := make(chan int64)
	go a.Exchange.UnbrokenSync(currency, candleChan, sessionChan)

	var currentSession int64
	for {
		select {
		case candle := <-candleChan:
			//! только синхронно
			success := a.Hostconnector.SendCandle(currency, currentSession, candle)
			if !success {
				log.Fatal("ERROR SEND: ", currency)
			}
		case newSession := <-sessionChan:
			currentSession = newSession
		}
	}
}
