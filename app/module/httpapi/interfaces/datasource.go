package interfaces

import (
	"robot/common/entity"
)

type DataSource interface {
	GetAllSessions() map[int64][]entity.Candle
}
