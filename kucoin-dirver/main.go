package main

import (
	"kucoin-dirver/bootstrap"
	"kucoin-dirver/module/exchange"
	"kucoin-dirver/module/hostconnector"
	"kucoin-dirver/module/httpapi"
	"os"
	"strings"
	"time"

	_ "net/http/pprof"
)

var app bootstrap.Bootstrap

func init() {
	app = bootstrap.Bootstrap{}
	app.Config = &bootstrap.Config{}
	app.Config.Read("config.yaml")

	app.Exchange = &exchange.App{}
	app.Exchange.Init(
		app.Config.EXCHANGE_API_KEY,
		app.Config.EXCHANGE_API_SECRET,
		app.Config.EXCHANGE_API_PASSPHRASE)

	app.Httpapi = &httpapi.App{
		Port: app.Config.HTTPAPI_PORT}

	app.Hostconnector = &hostconnector.App{
		app.Config.HOSTCONNECTOR_ADDRESS}
}

func main() {
	currency := strings.ToUpper(os.Args[1])
	app.Sync(currency)

	// не думаю что на споте я буду что-то делать. но пусть будет.
	// go app.Httpapi.OpenHost()

	for {
		time.Sleep(1 * time.Second)
	}
}
