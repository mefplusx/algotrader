package driver

import (
	"log"

	"kucoin-dirver/module/exchange/entity"

	"github.com/Kucoin/kucoin-go-sdk"
)

type Api struct {
	a *kucoin.ApiService
}

func NewApi(apiKey, apiSecret, apiPassphrase string) Api {
	// kucoin.DebugMode = true
	// kucoin.SetLoggerDirectory(".")
	a := kucoin.NewApiService(
		kucoin.ApiKeyOption(apiKey),
		kucoin.ApiSecretOption(apiSecret),
		kucoin.ApiPassPhraseOption(apiPassphrase),
		kucoin.ApiKeyVersionOption(kucoin.ApiKeyVersionV2))

	return Api{a}
}

func (a Api) UpdateCandleSubscribe(currency string, candleChan chan entity.WSCandle, reconnectChan chan bool) {
	log.Println("WS connect start")
	rsp, err := a.a.WebSocketPublicToken()
	if err != nil {
		log.Println("WebSocketPublicToken err:", err)
		reconnectChan <- true
		return
	}

	tk := &kucoin.WebSocketTokenModel{}
	if err := rsp.ReadData(tk); err != nil {
		log.Println("ReadData err:", err)
		reconnectChan <- true
		return
	}

	c := a.a.NewWebSocketClient(tk)

	mc, ec, err := c.Connect()
	if err != nil {
		log.Println("Connect err:", err)
		reconnectChan <- true
		return
	}

	ch := kucoin.NewSubscribeMessage("/market/candles:"+currency+"_1min", false)
	if err := c.Subscribe(ch); err != nil {
		log.Println("Subscribe err:", err)
		reconnectChan <- true
		return
	}
	log.Println("WS connect success")
	hasBeenFirstRead := false

	for {
		select {
		case err := <-ec:
			c.Stop()
			log.Printf("Error: %s\n", err.Error())
			reconnectChan <- true
			return
		case msg := <-mc:
			candle := &entity.WSCandle{}
			if candle == nil {
				c.Stop()
				log.Printf("Error: model == nil\n")
				reconnectChan <- true
				return
			}
			if err := msg.ReadData(candle); err != nil {
				log.Printf("Failure to read: %s\n", err.Error())
				reconnectChan <- true
				return
			}
			candleChan <- *candle
			if !hasBeenFirstRead {
				hasBeenFirstRead = true
				log.Println("First read from ws success:", candle.GetClose())
			}
		}
	}
}
