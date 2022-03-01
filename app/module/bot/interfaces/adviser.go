package interfaces

import (
	commonEntity "robot/common/entity"
)

type Adviser interface {
	// последний вход в сделку short/long
	// если нет входа тип отдаст TRANSACTOIN_TYPE_EMPTY
	GetLastIn() (int, commonEntity.MomentPath)

	// лимитный стоп на выход для последнего входа
	// если уже не в сделке выход будет 0
	GetStopForLastIn() float64
}
