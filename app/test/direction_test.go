package test

import (
	commonEntity "robot/common/entity"
	"robot/module/analytic/entity"
	"robot/module/analytic/middleware"
	"testing"
)

func TestDirection(t *testing.T) {
	moments := []commonEntity.Moment{}

	moments = append(moments, commonEntity.Moment{1, 1, commonEntity.MomentPath{}}) //+7cnt v7
	moments = append(moments, commonEntity.Moment{1, 1, commonEntity.MomentPath{}})
	moments = append(moments, commonEntity.Moment{1, 1, commonEntity.MomentPath{}})
	moments = append(moments, commonEntity.Moment{1, 2, commonEntity.MomentPath{}})
	moments = append(moments, commonEntity.Moment{1, 3, commonEntity.MomentPath{}})
	moments = append(moments, commonEntity.Moment{1, 4, commonEntity.MomentPath{}})
	moments = append(moments, commonEntity.Moment{1, 5, commonEntity.MomentPath{}})
	moments = append(moments, commonEntity.Moment{1, 4, commonEntity.MomentPath{}}) //-5cnt v5
	moments = append(moments, commonEntity.Moment{1, 3, commonEntity.MomentPath{}})
	moments = append(moments, commonEntity.Moment{1, 3, commonEntity.MomentPath{}})
	moments = append(moments, commonEntity.Moment{1, 3, commonEntity.MomentPath{}})
	moments = append(moments, commonEntity.Moment{1, 2, commonEntity.MomentPath{}})
	moments = append(moments, commonEntity.Moment{1, 3, commonEntity.MomentPath{}}) //+3 v2
	moments = append(moments, commonEntity.Moment{0.5, 4, commonEntity.MomentPath{}})
	moments = append(moments, commonEntity.Moment{0.5, 5, commonEntity.MomentPath{}})
	moments = append(moments, commonEntity.Moment{1, 4, commonEntity.MomentPath{}}) //-1 v1

	directions := middleware.MomentsToDirections(moments)

	if len(directions) != 4 {
		t.Fatalf(`count directions != 4, current value %v`, len(directions))
	}

	if directions[0].Value != 7 {
		t.Fatalf(`value fail 7`)
	}
	if directions[1].Value != 5 {
		t.Fatalf(`value fail 5`)
	}
	if directions[2].Value != 2 {
		t.Fatalf(`value fail 2`)
	}
	if directions[3].Value != 1 {
		t.Fatalf(`value fail 1`)
	}

	if directions[0].Direction != entity.DIRECTION_GROW {
		t.Fatalf(`direction fail 0`)
	}
	if directions[1].Direction != entity.DIRECTION_FALL {
		t.Fatalf(`direction fail 1`)
	}
	if directions[2].Direction != entity.DIRECTION_GROW {
		t.Fatalf(`direction fail 2`)
	}
	if directions[3].Direction != entity.DIRECTION_FALL {
		t.Fatalf(`direction fail 3`)
	}

	if len(directions[0].Moments) != 7 {
		t.Fatalf(`count on direction fail 7`)
	}
	if len(directions[1].Moments) != 5 {
		t.Fatalf(`count on direction fail 5`)
	}
	if len(directions[2].Moments) != 3 {
		t.Fatalf(`count on direction fail 3`)
	}
	if len(directions[3].Moments) != 1 {
		t.Fatalf(`count on direction fail 1`)
	}
}
