package main

import (
	"os"
	"robot/bootstrap"
	"robot/helper"
	"robot/module/analytic"
	"robot/module/bot"
	"robot/module/chart"
	"robot/module/filedatacandles"
	"robot/module/httpapi"
	"robot/module/remotedataconnector"
	"robot/module/vardatacandles"
	"time"

	_ "net/http/pprof"
)

var app bootstrap.Bootstrap

func init() {
	app = bootstrap.Bootstrap{}
	app.Config = &bootstrap.Config{}
	app.Config.Read("config.yaml")
	app.Vardatacandles = &vardatacandles.App{}

	filedatacandlesPath, _ := os.Getwd()
	filedatacandlesPath += "/module/filedatacandles/.data"
	app.FileDataCandles = &filedatacandles.App{
		app.Vardatacandles,
		filedatacandlesPath}
	app.FileDataCandles.AllToDataSource()

	app.Httpapi = &httpapi.App{
		DataSource: app.Vardatacandles,
		Port:       app.Config.HTTPAPI_PORT}

	app.RemoteDataConnector = &remotedataconnector.App{
		app.Config.REMOTEDATACONNECTOR_HOST}

	app.Chart = &chart.App{}

	app.Analytic = &analytic.App{
		app.Vardatacandles,
		app.Chart,
		app.Config.ANALYTIC_FEE,
		app.Config.ANALYTIC_TAKE_LINE_BY_SLIDEDOWN,
		app.Config.ANALYTIC_PERC_IN,
		app.Config.ANALYTIC_TAKE_AFTER,
		app.Config.ANALYTIC_SLIDEDOWN,
		app.Config.ANALYTIC_STOP,
		app.Config.ANALYTIC_VALUE_USDT,
		app.Config.ANALYTIC_REPEAT,
		app.Config.ANALYTIC_LEVERAGE,
		app.Config.ANALYTIC_ALLOW_SHORT,
		app.Config.ANALYTIC_ALLOW_LONG}

	botFileStatePath, _ := os.Getwd()
	botFileStatePath += "/module/bot/.data"
	app.Bot = &bot.App{
		Adviser:          app.Analytic,
		BotFileStatePath: botFileStatePath}
	app.Bot.Init()
}

func pprofRUN() {
	if helper.ExistsArgs("-pprof", "-p") {
		helper.RunPprof()
	}
}

func detailRUN() {
	const imgChartPath = "/tmp/chart-analytic.png"
	const initUSD float64 = 100
	if helper.ExistsArgs("-detailsimple", "-ds", "-detailfull", "-df") {
		if helper.ExistsArgs("-detailsimple", "-ds") {
			app.PrintDetail(initUSD, helper.GetInt64FromArgs(2, 999), analytic.FORMAT_SIMPLE)
		}
		if helper.ExistsArgs("-detailfull", "-df") {
			app.PrintDetail(initUSD, helper.GetInt64FromArgs(2, 999), analytic.FORMAT_FULL)
		}
		if helper.ExistsArgs("-chart", "-c") {
			app.CreateChart(initUSD, helper.GetInt64FromArgs(2, 999), imgChartPath, true)
		}
	}
}
func findRUN() {
	const initUSD float64 = 100
	if helper.ExistsArgs("-find", "-f") {
		app.Find(initUSD,
			helper.GetInt64FromArgs(2, 999),
			helper.GetInt64FromArgs(3, 100), //all
			helper.GetInt64FromArgs(4, 0),   //from
			helper.GetInt64FromArgs(5, 100)) //to
	}
}
func httpapiRUN() {
	if helper.ExistsArgs("-httpapi", "-ha") {
		go app.Httpapi.OpenHost(app.Config.HTTPAPI_HANDLERS...)
		for {
			time.Sleep(1 * time.Second)
		}
	}
}
func listfindsRUN() {
	if helper.ExistsArgs("-listfinds", "-lf", "-listfindssource", "-lfs") {
		app.ReadFindLogs(
			app.Config.REMOTEDATACONNECTOR_FINDNODES,
			helper.ExistsArgs("-listfindssource", "-lfs"))
	}
}

func main() {
	pprofRUN()
	detailRUN()
	findRUN()
	httpapiRUN()
	listfindsRUN()
}
