package entity

import (
	"robot/common"
	commonEntity "robot/common/entity"
)

type LimitOrder struct {
	Direction  Direction
	Type       uint8
	PriceForIn float64
}

func (mo *LimitOrder) GetId() int64 {
	return mo.Direction.Moments[0].Path.CandleTime
}

func (mo *LimitOrder) GetMoment() commonEntity.Moment {
	return mo.Direction.Moments[len(mo.Direction.Moments)-1]
}

func (mo *LimitOrder) IsEmpty() bool {
	return mo.Type == common.TRANSACTOIN_TYPE_EMPTY
}
