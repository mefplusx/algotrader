package analytic

import (
	"fmt"
	"math"
	"robot/common"
	commonEntity "robot/common/entity"
	"robot/helper"
	"robot/module/analytic/entity"
	analyticHelper "robot/module/analytic/helper"
	"robot/module/analytic/interfaces"
	"robot/module/analytic/middleware"
	"runtime"
)

// АНАЛИТИКА
// - знает что делать с историей движения
// - покажет результат с выбранными диапазонами
// - найдет диапазоны получше
// - скажет когда входить в сделку short/long по маркету
// - скажет лимитный стоп на выход
type App struct {
	DataSource interfaces.DataSource
	Chart      interfaces.Chart

	FEE                    float64
	TAKE_LINE_BY_SLIDEDOWN float64
	PERC_IN                float64
	TAKE_AFTER             float64
	SLIDEDOWN              float64
	STOP                   float64
	VALUE_USDT             float64
	REPEAT                 int64
	LEVERAGE               int64
	ALLOW_SHORT            bool
	ALLOW_LONG             bool
}

func (a *App) GetTransactions(diapason entity.Diapason, joinedMoments [][]commonEntity.Moment) (entity.LimitOrder, entity.Transaction, []entity.Transaction) {
	tragetOrder := entity.LimitOrder{}
	targetTrancation := entity.Transaction{}
	transactions := []entity.Transaction{}
	for _, moments := range joinedMoments {
		directions := middleware.MomentsToDirections(moments)

		_, _, orders := middleware.GetLimitOrders(
			directions,
			int(diapason.REPEAT), diapason.VALUE_USDT,
			diapason.PERC_IN,
			a.ALLOW_SHORT, a.ALLOW_LONG)

		calculator := entity.Calculator{moments, orders}
		var transactionsByGroup []entity.Transaction
		tragetOrder, targetTrancation, transactionsByGroup = calculator.GetTransactions(
			targetTrancation,
			diapason.TAKE_AFTER, diapason.SLIDEDOWN, diapason.TAKE_LINE_BY_SLIDEDOWN, diapason.STOP)
		transactions = append(transactions, transactionsByGroup...)
	}

	return tragetOrder, targetTrancation, transactions
}

func (a *App) CreateChart(initUSD float64, forLastDays int64, imgPath string) {
	joinedMoments := a.DataSource.JoinMoments(forLastDays)
	_, _, transactions := a.GetTransactions(entity.Diapason{
		TAKE_LINE_BY_SLIDEDOWN: a.TAKE_LINE_BY_SLIDEDOWN,
		PERC_IN:                a.PERC_IN,
		TAKE_AFTER:             a.TAKE_AFTER,
		SLIDEDOWN:              a.SLIDEDOWN,
		STOP:                   a.STOP,
		VALUE_USDT:             a.VALUE_USDT,
		REPEAT:                 a.REPEAT,
	}, joinedMoments)

	prices := []float64{}
	for _, item := range transactions {
		initUSD += initUSD / 100 * (item.GetResultPerc() * float64(a.LEVERAGE))
		initUSD -= initUSD / 100 * (a.FEE * float64(a.LEVERAGE))
		prices = append(prices, initUSD)
	}

	a.Chart.CreateChart(prices, imgPath)
}

const FORMAT_FULL = "full"
const FORMAT_SIMPLE = "simple"

func (a *App) PrintDetail(initUSD float64, forLastDays int64, printFormat string) {
	joinedMoments := a.DataSource.JoinMoments(forLastDays)

	_, lastTransaction, transactions := a.GetTransactions(entity.Diapason{
		TAKE_LINE_BY_SLIDEDOWN: a.TAKE_LINE_BY_SLIDEDOWN,
		PERC_IN:                a.PERC_IN,
		TAKE_AFTER:             a.TAKE_AFTER,
		SLIDEDOWN:              a.SLIDEDOWN,
		STOP:                   a.STOP,
		VALUE_USDT:             a.VALUE_USDT,
		REPEAT:                 a.REPEAT,
	}, joinedMoments)

	cntInColorDirect := ""
	beforeAmount := 0
	beforeColor := ""
	t := analyticHelper.Table{}
	var totalLong, totalShort float64
	var complicated float64 = 100
	var totalPerc float64 = 0
	for _, item := range transactions {
		typo := "S"
		if item.Type == common.TRANSACTOIN_TYPE_LONG {
			typo = "L"
		}
		color := "red"
		resultPerc := item.GetResultPerc() - a.FEE

		initUSD += initUSD / 100 * (item.GetResultPerc() * float64(a.LEVERAGE))
		initUSD -= initUSD / 100 * (a.FEE * float64(a.LEVERAGE))

		complicated += complicated / float64(100) * resultPerc
		if resultPerc > 0 {
			color = "green"
		}
		if printFormat == FORMAT_FULL {
			t.Add(typo, fmt.Sprintf("%.2f%%", resultPerc), "color="+color, int(initUSD), int(item.In.Close), int(item.Out.Close))
		} else if printFormat == FORMAT_SIMPLE {
			totalPerc += item.GetResultPerc() - a.FEE

			if color == beforeColor || beforeColor == "" {
				cntInColorDirect += "."
			} else {
				t.Add(
					fmt.Sprintf("%.2f%%", totalPerc),
					fmt.Sprintf("%d$", int(beforeAmount)),
					cntInColorDirect,
					"color="+beforeColor)
				cntInColorDirect = "."
			}
			beforeColor = color
			beforeAmount = int(initUSD)
		}
		if item.Type == common.TRANSACTOIN_TYPE_LONG {
			totalLong += resultPerc
		} else {
			totalShort += resultPerc
		}
	}
	if printFormat == FORMAT_SIMPLE {
		t.Add(
			fmt.Sprintf("%.2f%%", totalPerc),
			fmt.Sprintf("%d$", int(beforeAmount)),
			cntInColorDirect,
			"color="+beforeColor)
	}
	t.Print()

	if lastTransaction.Status != entity.TRANSACTOIN_STATUS_STOP && len(joinedMoments) > 0 {
		if lastTransaction.Type == common.TRANSACTOIN_TYPE_LONG {
			fmt.Println("last long transaction in:", lastTransaction.In.Close)
		} else {
			fmt.Println("last short transaction in:", lastTransaction.In.Close)
		}
	}
	fmt.Println("simple%:", int(totalLong+totalShort), "complicated%:", int(complicated-100))
	fmt.Println("long%:", int(totalLong), "short%:", int(totalShort))
	fmt.Println("count transactions:", len(transactions))
}

func getVariants(
	RULES []string,
	TAKE_LINE_BY_SLIDEDOWN []float64,
	PERC_IN []float64,
	TAKE_AFTER []float64,
	SLIDEDOWN []float64,
	STOP []float64,
	VALUE_USDT []float64,
	REPEAT []int64,
	FEE float64,
) []entity.Diapason {
	variants := []entity.Diapason{}
	for _TAKE_LINE_BY_SLIDEDOWN := TAKE_LINE_BY_SLIDEDOWN[0]; _TAKE_LINE_BY_SLIDEDOWN <= TAKE_LINE_BY_SLIDEDOWN[1]; _TAKE_LINE_BY_SLIDEDOWN += TAKE_LINE_BY_SLIDEDOWN[2] {
		for _PERC_IN := PERC_IN[0]; _PERC_IN <= PERC_IN[1]; _PERC_IN += PERC_IN[2] {
			for _TAKE_AFTER := TAKE_AFTER[0]; _TAKE_AFTER <= TAKE_AFTER[1]; _TAKE_AFTER += TAKE_AFTER[2] {
				for _SLIDEDOWN := SLIDEDOWN[0]; _SLIDEDOWN <= SLIDEDOWN[1]; _SLIDEDOWN += SLIDEDOWN[2] {
					for _STOP := STOP[0]; _STOP <= STOP[1]; _STOP += STOP[2] {
						for _VALUE_USDT := VALUE_USDT[0]; _VALUE_USDT <= VALUE_USDT[1]; _VALUE_USDT += VALUE_USDT[2] {
							for _REPEAT := REPEAT[0]; _REPEAT <= REPEAT[1]; _REPEAT += REPEAT[2] {
								if (helper.Contains(RULES, "TAKE_AFTER > STOP > SLIDEDOWN") && !(_TAKE_AFTER > _STOP && _STOP > _SLIDEDOWN)) ||
									(helper.Contains(RULES, "TAKE_AFTER - SLIDEDOWN > FEE") && !(_TAKE_AFTER-_SLIDEDOWN > FEE)) {
									continue
								}
								variants = append(variants, entity.Diapason{
									TAKE_LINE_BY_SLIDEDOWN: _TAKE_LINE_BY_SLIDEDOWN,
									PERC_IN:                _PERC_IN,
									TAKE_AFTER:             _TAKE_AFTER,
									SLIDEDOWN:              _SLIDEDOWN,
									STOP:                   _STOP,
									VALUE_USDT:             _VALUE_USDT,
									REPEAT:                 _REPEAT,
								})
							}
						}
					}
				}
			}
		}
	}

	return variants
}

func variantsToGroups(variants []entity.Diapason, groupCNT int) [][]entity.Diapason {
	groups := [][]entity.Diapason{}

	for groupID := 0; groupID < groupCNT; groupID++ {
		groups = append(groups, []entity.Diapason{})
	}

	groupID := 0
	for _, variant := range variants {
		groups[groupID] = append(groups[groupID], variant)

		groupID++
		if groupID > groupCNT-1 {
			groupID = 0
		}
	}

	return groups
}

const (
	SEARCH_TYPE_BSP = "best simple percent"
	SEARCH_TYPE_BTU = "best total usd"
	SEARCH_TYPE_T   = "toggle"
)

func (a *App) Find(
	initUSD float64, all, from, to int64,
	RULES []string,
	SEARCH_TYPE string,
	TAKE_LINE_BY_SLIDEDOWN []float64,
	PERC_IN []float64,
	TAKE_AFTER []float64,
	SLIDEDOWN []float64,
	STOP []float64,
	VALUE_USDT []float64,
	REPEAT []int64,
	MIN_COUNT_TRANSACTIONS int64,
	FOR_LAST_DAYS int64,
	FEE float64,
) entity.Diapason {
	cores := int(runtime.NumCPU())
	variants := getVariants(RULES, TAKE_LINE_BY_SLIDEDOWN, PERC_IN, TAKE_AFTER, SLIDEDOWN, STOP, VALUE_USDT, REPEAT, FEE)
	fmt.Println("all variants:", len(variants))

	variants = variants[int(from)*int(float64(len(variants))/float64(all)) : int(to)*int(math.Ceil(float64(len(variants))/float64(all)))]
	fmt.Println("for this node variants:", len(variants))

	groups := variantsToGroups(variants, cores)

	cntForThisNode := 0
	for _, group := range groups {
		cntForThisNode += len(group)
	}

	var befResultPerc float64 = 0
	bestChan := make(chan entity.Diapason)
	joinedMoments := a.DataSource.JoinMoments(FOR_LAST_DAYS)
	for _, variants := range groups {
		go func(_variants []entity.Diapason, _bestChan chan entity.Diapason) {
			best := entity.Diapason{}
			for variant, _ := range _variants {
				_variants[variant].TotalUSD = initUSD
				_, _, transactions := a.GetTransactions(_variants[variant], joinedMoments)
				for _, item := range transactions {
					_variants[variant].TotalUSD += _variants[variant].TotalUSD / 100 * (item.GetResultPerc() * float64(a.LEVERAGE))
					_variants[variant].TotalUSD -= _variants[variant].TotalUSD / 100 * (a.FEE * float64(a.LEVERAGE))

					resultPerc := item.GetResultPerc() - a.FEE

					if SEARCH_TYPE == SEARCH_TYPE_T {
						if variant > 0 {
							if (befResultPerc > 0 && resultPerc < 0) ||
								(befResultPerc < 0 && resultPerc > 0) {
								_variants[variant].ToggleCoef++
							} else {
								_variants[variant].ToggleCoef--
							}
						}
						befResultPerc = resultPerc
					}

					if item.Type == common.TRANSACTOIN_TYPE_LONG {
						_variants[variant].LongPerc += resultPerc
						_variants[variant].CntLongTransactions++
					} else {
						_variants[variant].ShortPerc += resultPerc
						_variants[variant].CntShortTransactions++
					}
				}

				if SEARCH_TYPE == SEARCH_TYPE_BSP {
					if _variants[variant].TotalCntTransactions() >= int(MIN_COUNT_TRANSACTIONS) &&
						_variants[variant].TotalSimplePerc() > best.TotalSimplePerc() {
						best = _variants[variant]
					}
				} else if SEARCH_TYPE == SEARCH_TYPE_T {
					if _variants[variant].TotalCntTransactions() >= int(MIN_COUNT_TRANSACTIONS) &&
						_variants[variant].ToggleCoef > best.ToggleCoef {
						best = _variants[variant]
					}
				} else if SEARCH_TYPE == SEARCH_TYPE_BTU {
					if _variants[variant].TotalCntTransactions() >= int(MIN_COUNT_TRANSACTIONS) &&
						_variants[variant].TotalUSD > best.TotalUSD {
						best = _variants[variant]
					}
				}
			}

			_bestChan <- best
		}(variants, bestChan)
	}

	best := entity.Diapason{}
	for cores > 0 {
		select {
		case current := <-bestChan:
			if SEARCH_TYPE == SEARCH_TYPE_BSP {
				if current.TotalSimplePerc() > best.TotalSimplePerc() {
					best = current
				}
			} else if SEARCH_TYPE == SEARCH_TYPE_T {
				if current.ToggleCoef > best.ToggleCoef {
					best = current
				}
			} else if SEARCH_TYPE == SEARCH_TYPE_BTU {
				if current.TotalUSD > best.TotalUSD {
					best = current
				}
			}

			cores--
		}
	}

	return best
}

// последний вход в сделку short/long
// если нет входа тип отдаст TRANSACTOIN_TYPE_EMPTY
func (a *App) GetLastIn() (int, commonEntity.MomentPath) {
	return common.TRANSACTOIN_TYPE_EMPTY, commonEntity.MomentPath{}
}

// лимитный стоп на выход для последнего входа
// если уже не всделке выход будет 0
func (a *App) GetStopForLastIn() float64 {
	return 0
}
