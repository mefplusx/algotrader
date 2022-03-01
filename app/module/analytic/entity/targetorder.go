package entity

import (
	"fmt"
	"robot/common"
)

const TARGET_ORDER_LONG = "long"
const TARGET_ORDER_SHORT = "short"

type TargetOrder struct {
	Key     string
	Type    string
	PirceIn float64
	Take    float64
	Stop    float64
}

func (mo *LimitOrder) ToTargetOrder(PERC_IN, TAKE_AFTER, SLIDEDOWN, STOP float64) TargetOrder {
	to := TargetOrder{}

	lastMI := mo.Direction.GetLastMoment()
	to.Key = fmt.Sprintf("%d-%d", lastMI.Path.CandleTime, lastMI.Path.InCandleOrder)
	if mo.Type == common.TRANSACTOIN_TYPE_LONG {
		to.Type = TARGET_ORDER_LONG
		to.PirceIn = lastMI.Close - (lastMI.Close / float64(100) * PERC_IN)
		if STOP > 0 {
			to.Stop = lastMI.Close - (lastMI.Close / float64(100) * STOP)
		}
		to.Take = lastMI.Close + (lastMI.Close / float64(100) * TAKE_AFTER)
	} else {
		to.Type = TARGET_ORDER_SHORT
		to.PirceIn = lastMI.Close + (lastMI.Close / float64(100) * PERC_IN)
		if STOP > 0 {
			to.Stop = lastMI.Close + (lastMI.Close / float64(100) * STOP)
		}
		to.Take = lastMI.Close - (lastMI.Close / float64(100) * TAKE_AFTER)
	}

	return to
}
