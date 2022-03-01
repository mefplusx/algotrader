package middleware

import (
	commonEntity "robot/common/entity"
	"robot/module/analytic/entity"
)

const (
	DIRECTION_UNKNOWN = iota
	DIRECTION_GROW
	DIRECTION_FALL
)

func MomentsToDirections(moments []commonEntity.Moment) []entity.Direction {
	directions := []entity.Direction{}

	direction := entity.Direction{}
	for i, _ := range moments {
		if direction.Direction != DIRECTION_UNKNOWN {
			//смена направления
			if (direction.Direction == DIRECTION_GROW && direction.Close > moments[i].Close) ||
				(direction.Direction == DIRECTION_FALL && direction.Close < moments[i].Close) {
				directions = append(directions, direction)

				var nd uint8
				if direction.Direction == DIRECTION_GROW {
					nd = DIRECTION_FALL
				} else {
					nd = DIRECTION_GROW
				}
				direction = entity.Direction{
					[]commonEntity.Moment{}, nd, 0, 0}
			}
		} else {
			if i > 0 {
				if moments[i].Close > moments[i-1].Close {
					direction.Direction = DIRECTION_GROW
				} else if moments[i].Close < moments[i-1].Close {
					direction.Direction = DIRECTION_FALL
				}
			}
		}

		direction.Moments = append(direction.Moments, moments[i])
		direction.Value += moments[i].Value
		direction.Close = moments[i].Close
	}

	if len(direction.Moments) > 0 {
		directions = append(directions, direction)
	}

	return directions
}
