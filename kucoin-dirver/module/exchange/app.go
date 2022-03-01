package exchange

import (
	"fmt"
	"kucoin-dirver/common/entity"
	"kucoin-dirver/module/exchange/driver"
	Entity "kucoin-dirver/module/exchange/entity"
	"log"
	"time"
)

type App struct {
	typo string
	api  driver.Api
}

func (a *App) Init(key, secret, passphrase string) {
	a.api = driver.NewApi(key, secret, passphrase)
}

func (a *App) FuturesUnbrokenSync(pair string, candleChan chan entity.Candle, sessionChan chan int64) {
	fmt.Println("functional not implemented")
}

func (a *App) UnbrokenSync(pair string, candleChan chan entity.Candle, sessionChan chan int64) {
	candle := entity.Candle{}
	for {
		session := time.Now().Unix()
		sessionChan <- session
		log.Println("Open new session", session)
		wscandleChan := make(chan Entity.WSCandle)
		reconnectChan := make(chan bool)
		go a.api.UpdateCandleSubscribe(pair, wscandleChan, reconnectChan)

		for isWSConnected := true; isWSConnected; {
			select {
			case wscandle := <-wscandleChan:
				wscandleTime := wscandle.GetTime()
				wscandleClose := wscandle.GetClose()
				wscandleVol := wscandle.GetVol()

				if wscandleClose != 0 && wscandleTime != 0 && wscandleVol != 0 {
					moment := entity.Moment{
						wscandleVol,
						wscandleClose,
						entity.MomentPath{
							candle.Time, int64(len(candle.Moments)),
						}}
					if wscandleTime != candle.Time {
						candle = entity.Candle{wscandleTime, []entity.Moment{moment}}
					} else {
						candle.Moments = append(candle.Moments, moment)
					}

					candleChan <- candle
				}
			case <-reconnectChan:
				candle = entity.Candle{}
				log.Println("reconnect spy...")
				isWSConnected = false
			}
		}
	}
}

func (a *App) BuyMarket(typo string) bool {
	return true
}
func (a *App) SetStop(typo string, price float64) {

}
