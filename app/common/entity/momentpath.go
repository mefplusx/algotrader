package entity

import "fmt"

type MomentPath struct {
	CandleTime    int64
	InCandleOrder int64
}

func (mp *MomentPath) IsIN(bullets []MomentPath) bool {
	if mp.CandleTime == 0 || mp.InCandleOrder == 0 {
		return false
	}

	for _, bullet := range bullets {
		if bullet.CandleTime == mp.CandleTime &&
			bullet.InCandleOrder == mp.InCandleOrder {
			return true
		}
	}

	return false
}

func (mp *MomentPath) GetString() string {
	return fmt.Sprintf("%d.%d", mp.CandleTime, mp.InCandleOrder)
}

func (mp *MomentPath) IsEmpty() bool {
	return mp.CandleTime == 0 && mp.InCandleOrder == 0
}
