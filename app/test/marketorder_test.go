package test

import (
	commonEntity "robot/common/entity"
	"robot/module/analytic/middleware"
	"testing"
)

func TestGetLimitOrders(t *testing.T) {
	items := []commonEntity.Moment{}

	items = append(items, commonEntity.Moment{1, 1, commonEntity.MomentPath{}})
	items = append(items, commonEntity.Moment{1, 2, commonEntity.MomentPath{}}) // 2 not in
	items = append(items, commonEntity.Moment{0.1, 0.1, commonEntity.MomentPath{}})
	items = append(items, commonEntity.Moment{0.1, 0.1, commonEntity.MomentPath{}})
	items = append(items, commonEntity.Moment{0.1, 0.1, commonEntity.MomentPath{}}) //4 not in short by value
	items = append(items, commonEntity.Moment{1, 2, commonEntity.MomentPath{}})
	items = append(items, commonEntity.Moment{1, 3, commonEntity.MomentPath{}})
	items = append(items, commonEntity.Moment{1, 4, commonEntity.MomentPath{}})
	items = append(items, commonEntity.Moment{1, 5, commonEntity.MomentPath{}})
	items = append(items, commonEntity.Moment{1, 5, commonEntity.MomentPath{}})
	items = append(items, commonEntity.Moment{1, 5, commonEntity.MomentPath{}})
	items = append(items, commonEntity.Moment{1, 5, commonEntity.MomentPath{0, 1}}) // 4 in long, value ok
	//											   price 4.75
	items = append(items, commonEntity.Moment{1, 4, commonEntity.MomentPath{}})
	items = append(items, commonEntity.Moment{1, 3, commonEntity.MomentPath{}})
	items = append(items, commonEntity.Moment{1, 3, commonEntity.MomentPath{}})
	items = append(items, commonEntity.Moment{1, 3, commonEntity.MomentPath{}})
	items = append(items, commonEntity.Moment{1, 2, commonEntity.MomentPath{}})
	items = append(items, commonEntity.Moment{1, 2, commonEntity.MomentPath{}})
	items = append(items, commonEntity.Moment{1, 2, commonEntity.MomentPath{0, 2}}) // 3 in short, value ok
	//											   price 2.1
	items = append(items, commonEntity.Moment{0.5, 4, commonEntity.MomentPath{}})
	items = append(items, commonEntity.Moment{0.5, 5, commonEntity.MomentPath{}}) // 2 not in
	items = append(items, commonEntity.Moment{1, 4, commonEntity.MomentPath{}})   // 1 not in

	directions := middleware.MomentsToDirections(items)
	_, _, orders := middleware.GetLimitOrders(directions, 3, 1, 5, true, true)

	if len(orders) != 2 {
		t.Fatalf(`count orders != 2, current value %v`, len(directions))
	}

	if orders[commonEntity.MomentPath{0, 1}].PriceForIn != 4.75 {
		t.Fatalf(`fail price long in %v`, orders[commonEntity.MomentPath{0, 1}].PriceForIn)
	}

	if orders[commonEntity.MomentPath{0, 2}].PriceForIn != 2.1 {
		t.Fatalf(`fail price long in %v`, orders[commonEntity.MomentPath{0, 2}].PriceForIn)
	}
}
