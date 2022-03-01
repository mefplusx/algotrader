package filedatacandles

import (
	"fmt"
	"os"
	"robot/common/entity"
	"robot/module/filedatacandles/interfaces"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const ext = ".json"

type App struct {
	DataSource interfaces.DataSource
	Path       string
}

func (a *App) AllToDataSource() {
	allSessions := make(map[int64][]entity.Candle)

	files, _ := os.ReadDir(a.Path)
	for _, file := range files {
		fileName := file.Name()
		if !strings.Contains(fileName, ext) {
			continue
		}

		sessionTime, _ := strconv.ParseInt(fileName[:len(fileName)-len(ext)], 10, 64)
		b, _ := os.ReadFile(a.Path + "/" + fileName)
		candles := []entity.Candle{}
		json.Unmarshal(b, &candles)
		allSessions[sessionTime] = candles
	}

	a.DataSource.SetAllSessions(allSessions)
}

func (a *App) CurrentSessionToFile() {
	candles := a.DataSource.GetCurrentSessionCandles()
	candlesJson, _ := json.Marshal(candles)

	session := a.DataSource.GetCurrentSession()
	path := fmt.Sprintf("%s/%v"+ext, a.Path, session)
	file, _ := os.Create(path)
	file.Write(candlesJson)
	file.Close()
}
