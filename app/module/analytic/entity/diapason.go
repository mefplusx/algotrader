package entity

import (
	"fmt"
	"robot/common"
)

type Diapason struct {
	TAKE_LINE_BY_SLIDEDOWN float64
	PERC_IN                float64
	TAKE_AFTER             float64
	SLIDEDOWN              float64
	STOP                   float64
	VALUE_USDT             float64
	REPEAT                 int64

	ShortPerc float64
	LongPerc  float64

	CntShortTransactions int
	CntLongTransactions  int

	TotalUSD float64

	// чем выше коэффичиент, тем стабильнее чередование сделок красный/зеленый
	ToggleCoef int64
}

func (d *Diapason) TotalCntTransactions() int {
	return d.CntLongTransactions + d.CntShortTransactions
}

func (d *Diapason) TotalSimplePerc() float64 {
	return d.ShortPerc + d.LongPerc
}

func (d *Diapason) Print() {
	fmt.Println(common.FINDER_SEPARATOR)
	fmt.Println("ANALYTIC_TAKE_LINE_BY_SLIDEDOWN :", d.TAKE_LINE_BY_SLIDEDOWN)
	fmt.Println("ANALYTIC_PERC_IN                :", d.PERC_IN)
	fmt.Println("ANALYTIC_TAKE_AFTER             :", d.TAKE_AFTER)
	fmt.Println("ANALYTIC_SLIDEDOWN              :", d.SLIDEDOWN)
	fmt.Println("ANALYTIC_STOP                   :", d.STOP)
	fmt.Println("ANALYTIC_VALUE_USDT             :", d.VALUE_USDT)
	fmt.Println("ANALYTIC_REPEAT                 :", d.REPEAT)
	fmt.Println("")
	fmt.Println("short% :", int(d.ShortPerc), ", cnt:", d.CntShortTransactions)
	fmt.Println("long%  :", int(d.LongPerc), ", cnt:", d.CntLongTransactions)
	fmt.Println("total$ :", int(d.TotalUSD))

	if d.ToggleCoef > 0 {
		fmt.Println("toggle coef :", d.ToggleCoef)
	}
}
