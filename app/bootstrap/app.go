package bootstrap

import (
	"fmt"
	"os/exec"
	"robot/common"
	"robot/module/analytic"
	"robot/module/bot"
	"robot/module/chart"
	"robot/module/filedatacandles"
	"robot/module/httpapi"
	"robot/module/remotedataconnector"
	"robot/module/vardatacandles"
	"strings"
)

type Bootstrap struct {
	Config              *Config
	Vardatacandles      *vardatacandles.App
	FileDataCandles     *filedatacandles.App
	Httpapi             *httpapi.App
	RemoteDataConnector *remotedataconnector.App
	Chart               *chart.App
	Analytic            *analytic.App
	Bot                 *bot.App
}

func (a Bootstrap) PrintDetail(initUSD float64, forLastDays int64, printFormat string) {
	sessions := a.RemoteDataConnector.GetSessions(
		a.Config.REMOTEDATACONNECTOR_USE_CACHE)
	a.Vardatacandles.SetAllSessions(sessions)
	a.Analytic.PrintDetail(initUSD, forLastDays, printFormat)
}

func (a Bootstrap) CreateChart(initUSD float64, forLastDays int64, imgPath string, andOpen bool) {
	a.Analytic.CreateChart(initUSD, forLastDays, imgPath)
	if andOpen {
		cmd := exec.Command("eog", imgPath)
		cmd.Start()
		cmd.Wait()
	}
}

func (a Bootstrap) Find(initUSD float64, forLastDays int64, all, from, to int64) {
	sessions := a.RemoteDataConnector.GetSessions(
		a.Config.REMOTEDATACONNECTOR_USE_CACHE)
	a.Vardatacandles.SetAllSessions(sessions)
	best := a.Analytic.Find(
		initUSD, all, from, to,
		a.Config.FINDER_RULES,
		a.Config.FINDER_TYPE,
		a.Config.FINDER_TAKE_LINE_BY_SLIDEDOWN,
		a.Config.FINDER_PERC_IN,
		a.Config.FINDER_TAKE_AFTER,
		a.Config.FINDER_SLIDEDOWN,
		a.Config.FINDER_STOP,
		a.Config.FINDER_VALUE_USDT,
		a.Config.FINDER_REPEAT,
		a.Config.FINDER_MIN_COUNT_TRANSACTIONS,
		forLastDays,
		a.Config.ANALYTIC_FEE,
	)

	best.Print()
}

// func (a Bootstrap) SyncBot(typo string, botmode bool) {
// 	candleChan := make(chan entity.Candle)
// 	sessionChan := make(chan int64)

// 	for {
// 		select {
// 		case candle := <-candleChan:
// 			if a.Vardatacandles.GetCurrentSession() != 0 {
// 				isNew := a.Vardatacandles.SetCandle(candle)
// 				if isNew {
// 					a.FileDataCandles.CurrentSessionToFile()
// 				}

// 				if botmode {
// 					a.Bot.Do()
// 				}
// 			}
// 		case session := <-sessionChan:
// 			if a.Vardatacandles.GetCurrentSession() != 0 {
// 				a.FileDataCandles.CurrentSessionToFile()
// 			}
// 			a.Vardatacandles.SetCurrentSession(session)
// 		}
// 	}
// }

func (a Bootstrap) ReadFindLogs(nodes []string, source bool) {
	logs := a.RemoteDataConnector.GetLogs(nodes)
	for node, log := range logs {
		fmt.Print(node + ":")

		if source {
			fmt.Println()
			fmt.Println(log)
			continue
		}

		if parts := strings.Split(log, common.FINDER_SEPARATOR); len(parts) == 2 {
			fmt.Println(parts[1])
		} else {
			fmt.Println("") // nope....
		}
	}
}
