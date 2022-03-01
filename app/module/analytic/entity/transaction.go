package entity

import (
	"robot/common"
	commonEntity "robot/common/entity"
)

const (
	TRANSACTOIN_STATUS_WAIT_IN = iota + 1
	TRANSACTOIN_STATUS_IN
	TRANSACTOIN_STATUS_STOP
	TRANSACTOIN_STATUS_BREACK
)

type Transaction struct {
	Type   uint8
	Status uint8

	Order LimitOrder

	In  commonEntity.Moment
	Out commonEntity.Moment

	Take float64
	Stop float64
}

func (t *Transaction) IsEmpty() bool {
	return t.In.Close == 0
}

func (t *Transaction) InLong(order LimitOrder, in commonEntity.Moment, stopPerc float64, takePerc float64) {
	t.Type = common.TRANSACTOIN_TYPE_LONG
	t.Order = order
	t.In = in
	t.Status = TRANSACTOIN_STATUS_IN
	if stopPerc > 0 {
		t.Stop = order.Direction.Close - (order.Direction.Close / float64(100) * stopPerc)
	}

	t.Take = order.Direction.Close + (order.Direction.Close / float64(100) * takePerc)
}

func (t *Transaction) InShort(order LimitOrder, in commonEntity.Moment, stopPerc float64, takePerc float64) {
	t.Type = common.TRANSACTOIN_TYPE_SHORT
	t.Order = order
	t.In = in
	t.Status = TRANSACTOIN_STATUS_IN
	if stopPerc > 0 {
		t.Stop = order.Direction.Close + (order.Direction.Close / float64(100) * stopPerc)
	}
	t.Take = order.Direction.Close - (order.Direction.Close / float64(100) * takePerc)
}

func (t *Transaction) GetResultPerc() float64 {
	inPrice := t.Order.PriceForIn
	var outPrice float64
	if t.Status == TRANSACTOIN_STATUS_BREACK {
		outPrice = t.Out.Close
	} else if t.Status == TRANSACTOIN_STATUS_STOP {
		outPrice = t.Stop
	}

	if inPrice == 0 {
		return 0
	}

	result := ((float64(100) / inPrice) * outPrice) - float64(100)
	if t.Type == common.TRANSACTOIN_TYPE_SHORT {
		result *= -1
	}

	return result
}
