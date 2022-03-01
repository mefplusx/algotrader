package interfaces

import (
	"kucoin-dirver/common/entity"
)

type DataSource interface {
	GetAllSessions() map[int64][]entity.Candle
}
