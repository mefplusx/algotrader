package hostconnector

import (
	"fmt"
	"kucoin-dirver/common/entity"
	"kucoin-dirver/module/hostconnector/internal"
	"strings"
)

type App struct {
	HostAddress string
}

func (a *App) SendCandle(currency string, session int64, candle entity.Candle) bool {
	url := strings.Join([]string{
		a.HostAddress,
		"set-candle",
		strings.ToLower(currency),
		fmt.Sprintf("%d", session)}, "/")

	sender := internal.Sender{url, candle}
	return sender.Send()
}
