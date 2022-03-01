package vardatacandles

import (
	"robot/common/entity"
	"sort"
)

type App struct {
	currentSession int64
	sessions       map[int64][]entity.Candle
}

func (a *App) SetAllSessions(allSessions map[int64][]entity.Candle) {
	a.sessions = allSessions
}

func (a *App) GetAllSessions() map[int64][]entity.Candle {
	return a.sessions
}

func (a *App) GetCurrentSessionCandles() []entity.Candle {
	if _, isExists := a.sessions[a.currentSession]; isExists {
		return a.sessions[a.currentSession]
	}

	return []entity.Candle{}
}

func (a *App) GetCurrentSession() int64 {
	return a.currentSession
}

func (a *App) SetCurrentSession(currentSession int64) {
	a.currentSession = currentSession
	if a.sessions == nil {
		a.sessions = make(map[int64][]entity.Candle)
	}
	if _, isExists := a.sessions[a.currentSession]; !isExists {
		a.sessions[a.currentSession] = []entity.Candle{}
	}
}

func (a *App) SetCandle(candle entity.Candle) bool {
	l := len(a.sessions[a.currentSession])

	if l == 0 {
		a.sessions[a.currentSession] = append(
			a.sessions[a.currentSession], candle)
		return false
	}

	if a.sessions[a.currentSession][l-1].Time != candle.Time {
		a.sessions[a.currentSession] = append(
			a.sessions[a.currentSession], candle)
		return true
	}

	a.sessions[a.currentSession][l-1] = candle
	return false
}

//разделить на группы, через разрывы между свечами
func (a *App) JoinMoments(forLastDays int64) [][]entity.Moment {
	const CANDLE_INTERVAL = 60
	groups := [][]entity.Moment{}

	if len(a.sessions) == 0 {
		return groups
	}

	var timeFrom int64 = a.GetMaxTime() - (forLastDays * 24 * 60 * 60)
	keys := a.GetSortedSessionsKeys()
	for _, key := range keys {
		if len(a.sessions[key]) == 0 {
			continue
		}

		group := a.sessions[key][0].Moments
		lastGroupTime := a.sessions[key][0].Time
		for i := 1; i < len(a.sessions[key]); i++ {
			if timeFrom > a.sessions[key][i].Time {
				continue
			}
			if lastGroupTime+CANDLE_INTERVAL != a.sessions[key][i].Time && len(group) > 0 {
				groups = append(groups, group)
				group = []entity.Moment{}
			}
			group = append(group, a.sessions[key][i].Moments...)
			lastGroupTime = a.sessions[key][i].Time
		}

		if len(group) > 0 {
			groups = append(groups, group)
		}
	}

	return groups
}

func (a *App) GetMaxTime() int64 {
	var t int64 = 0

	for _, candles := range a.sessions {
		for _, candle := range candles {
			if candle.Time > t {
				t = candle.Time
			}
		}
	}

	return t
}

// сортировка сессий
func (a *App) GetSortedSessionsKeys() []int64 {
	keys := []int64{}

	for key, _ := range a.sessions {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	return keys
}
