package remotedataconnector

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"robot/common/entity"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type App struct {
	HostAddress string
}

//получить ссылки на json файлы дирректории
func (a *App) getFiles() map[int64]string {
	files := make(map[int64]string)
	resp, _ := http.Get(a.HostAddress + "/data/")
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	content := strings.ReplaceAll(strings.ToLower(string(body)), " ", "")
	items := regexp.MustCompile(`[\d\.json]+`).FindAllString(content, -1)
	for _, item := range items {
		if strings.Contains(item, ".json") {
			sessionId, _ := strconv.ParseInt(strings.ReplaceAll(item, ".json", ""), 10, 64)
			files[sessionId] = a.HostAddress + "/data/" + item
		}
	}

	return files
}

//загрузить файлы сессий
func (a *App) getSessionsByUrls(files map[int64]string) map[int64][]entity.Candle {
	sessions := make(map[int64][]entity.Candle)

	type bullet struct {
		Session int64
		Candles []entity.Candle
	}
	bulletChan := make(chan bullet)

	cnt := 0
	for sessionId, fileUrl := range files {
		cnt++

		go (func(_sessionId int64, _fileUrl string, _bulletChan chan bullet) {
			resp, _ := http.Get(_fileUrl)
			decoded := []entity.Candle{}
			json.NewDecoder(resp.Body).Decode(&decoded)
			resp.Body.Close()
			bulletChan <- bullet{_sessionId, decoded}
		})(sessionId, fileUrl, bulletChan)
	}

	fmt.Println("\nLEFT FILES IN LOADING:")
	for cnt > 0 {
		select {
		case current := <-bulletChan:
			sessions[current.Session] = current.Candles
			cnt--
			fmt.Print(cnt, " ")
		}
	}
	fmt.Println("")

	return sessions
}

//пересчитать momentPath
func (a *App) recalcMomentPath(sessions map[int64][]entity.Candle) map[int64][]entity.Candle {
	for sessionId, _ := range sessions {
		for candleId, candle := range sessions[sessionId] {
			for momentOrder, _ := range sessions[sessionId][candleId].Moments {
				sessions[sessionId][candleId].Moments[momentOrder].Path =
					entity.MomentPath{
						candle.Time,
						int64(momentOrder)}
			}
		}
	}

	return sessions
}

func (a *App) loadSessionsFromCache() map[int64][]entity.Candle {
	sessions := map[int64][]entity.Candle{}

	filedatacandlesPath, _ := os.Getwd()
	if _, err := os.Stat(filedatacandlesPath + "/module/remotedataconnector/.data/cache.json"); errors.Is(err, os.ErrNotExist) {
		return sessions
	}

	b, _ := os.ReadFile(filedatacandlesPath + "/module/remotedataconnector/.data/cache.json")
	json.Unmarshal(b, &sessions)

	return sessions
}

func (a *App) saveSessionsToCache(sessions map[int64][]entity.Candle) {
	sessionsJson, _ := json.Marshal(sessions)
	filedatacandlesPath, _ := os.Getwd()
	file, _ := os.Create(filedatacandlesPath + "/module/remotedataconnector/.data/cache.json")
	file.Write(sessionsJson)
	file.Close()
}

func (a *App) GetSessions(useCache bool) map[int64][]entity.Candle {
	if useCache {
		sessionsFromCache := a.loadSessionsFromCache()
		if len(sessionsFromCache) > 0 {
			return sessionsFromCache
		}
	}

	files := a.getFiles()
	sessions := a.getSessionsByUrls(files)
	sessions = a.recalcMomentPath(sessions)

	if useCache {
		a.saveSessionsToCache(sessions)
	}

	return sessions
}

func (a *App) GetLogs(nodes []string) map[string]string {
	logs := make(map[string]string)

	for _, node := range nodes {
		resp, err := http.Get(node + "/log/find")
		if err != nil {
			continue
		}

		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		logs[node] = string(body)
	}

	return logs
}
