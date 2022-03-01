package entity

import (
	"robot/common/entity"
)

const (
	DIRECTION_UNKNOWN = iota
	DIRECTION_GROW
	DIRECTION_FALL
)

type Direction struct {
	Moments   []entity.Moment
	Direction uint8

	Close float64
	Value float64
}

func (d *Direction) GetLastMoment() entity.Moment {
	return d.Moments[len(d.Moments)-1]
}
