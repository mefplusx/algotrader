package interfaces

import (
	"robot/common/entity"
)

type DataSource interface {
	SetAllSessions(allSessions map[int64][]entity.Candle)
	GetCurrentSessionCandles() []entity.Candle
	GetCurrentSession() int64
}
