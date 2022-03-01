package middleware

import (
	"robot/common"
	commonEntity "robot/common/entity"
	"robot/module/analytic/entity"
)

func GetLimitOrders(
	directions []entity.Direction,
	minCntRepeat int, minValue float64,
	percIn float64,
	allowShort, allowLong bool) (commonEntity.MomentPath, entity.LimitOrder, map[commonEntity.MomentPath]entity.LimitOrder) {

	lastMovesMomentPath := commonEntity.MomentPath{}
	lastLimitOrder := entity.LimitOrder{}
	orders := map[commonEntity.MomentPath]entity.LimitOrder{}

	for _, d := range directions {
		//объем слишком мал
		if d.Value <= minValue {
			continue
		}

		//не достаточно элементов внутри дирэкта
		l := len(d.Moments)
		if l < minCntRepeat {
			continue
		}

		//или последние n не равны последней цены директа
		for i := l - minCntRepeat; i < l; i++ {
			if d.Moments[i].Close != d.Close {
				break
			}
		}

		if allowLong && d.Direction == entity.DIRECTION_GROW {
			last := d.GetLastMoment()
			lastLimitOrder = entity.LimitOrder{d, common.TRANSACTOIN_TYPE_LONG, d.Close - (d.Close / float64(100) * percIn)}
			lastMovesMomentPath = last.Path
			orders[lastMovesMomentPath] = lastLimitOrder
		}

		if allowShort && d.Direction == entity.DIRECTION_FALL {
			last := d.GetLastMoment()
			lastLimitOrder = entity.LimitOrder{d, common.TRANSACTOIN_TYPE_SHORT, d.Close + (d.Close / float64(100) * percIn)}
			lastMovesMomentPath = last.Path
			orders[last.Path] = lastLimitOrder
		}
	}

	return lastMovesMomentPath, lastLimitOrder, orders
}
