package entity

import (
	"robot/common"
	commonEntity "robot/common/entity"
)

type Calculator struct {
	Moves  []commonEntity.Moment
	Orders map[commonEntity.MomentPath]LimitOrder
}

func (c *Calculator) GetTransactions(
	continueTrancation Transaction,
	takeAfter, slidedown float64, TAKE_LINE_BY_SLIDEDOWN float64, stopPerc float64) (LimitOrder, Transaction, []Transaction) {
	transactions := []Transaction{}
	transaction := continueTrancation
	appendTransaction := func(status uint8, out commonEntity.Moment) {
		if transaction.In.Close != 0 {
			transaction.Status = status
			transaction.Out = out
			transactions = append(transactions, transaction)
			transaction = Transaction{}
		}
	}

	order := LimitOrder{}
	for _, moveItem := range c.Moves {
		//market order как точка отсчета на вход, если не в транзакции
		if _order, ok := c.Orders[moveItem.Path]; ok {
			if transaction.In.Close == 0 {
				order = _order
			}
		}

		if order.PriceForIn != 0 { //ищем вход
			//или не нарушает продолжение движения логике входа
			if order.Type == common.TRANSACTOIN_TYPE_LONG {
				if moveItem.Close > order.Direction.Close {
					order = LimitOrder{}
					continue
				}
			} else {
				if moveItem.Close < order.Direction.Close {
					order = LimitOrder{}
					continue
				}
			}

			//ожидание входа
			if transaction.In.Close == 0 {
				if order.Type == common.TRANSACTOIN_TYPE_LONG {
					if order.PriceForIn >= moveItem.Close { //вход по лонгу
						transaction.InLong(order, moveItem, stopPerc, takeAfter)
						order = LimitOrder{} //если вход, затираем ордер
						continue
					}
				} else {
					if order.PriceForIn <= moveItem.Close { //вход по шорту
						transaction.InShort(order, moveItem, stopPerc, takeAfter)
						order = LimitOrder{} //если вход, затираем ордер
						continue
					}
				}
			}
		} else if transaction.In.Close != 0 { //ищем выход
			if transaction.Type == common.TRANSACTOIN_TYPE_LONG {
				if moveItem.Close >= transaction.Take { //long take
					//вошли в зону тэйк
					//передвигаем стоп сеткой slidedown
					if slidedown == 0 {
						appendTransaction(TRANSACTOIN_STATUS_STOP, moveItem)
					} else {
						transaction.Take = moveItem.Close + (moveItem.Close / float64(100) * TAKE_LINE_BY_SLIDEDOWN)
						transaction.Stop = moveItem.Close - (moveItem.Close / float64(100) * slidedown)
					}
					continue
				} else if transaction.Stop > 0 && moveItem.Close <= transaction.Stop { //long stop
					appendTransaction(TRANSACTOIN_STATUS_STOP, moveItem)
					continue
				}
			} else {
				if moveItem.Close <= transaction.Take { //short take
					//вошли в зону тэйк
					//передвигаем стоп сеткой slidedown
					if slidedown == 0 {
						appendTransaction(TRANSACTOIN_STATUS_STOP, moveItem)
					} else {
						transaction.Take = moveItem.Close - (moveItem.Close / float64(100) * TAKE_LINE_BY_SLIDEDOWN)
						transaction.Stop = moveItem.Close + (moveItem.Close / float64(100) * slidedown)
					}
					continue
				} else if transaction.Stop > 0 && moveItem.Close >= transaction.Stop { //short stop
					appendTransaction(TRANSACTOIN_STATUS_STOP, moveItem)
					continue
				}
			}
		}
	}

	return order, transaction, transactions
}
