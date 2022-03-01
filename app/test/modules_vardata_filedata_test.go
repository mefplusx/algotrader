package test

import (
	"os"
	"robot/common/entity"
	"robot/module/filedatacandles"
	"robot/module/vardatacandles"
	"testing"
)

func TestFileDataCandlesVardatacandles(t *testing.T) {
	dir := t.TempDir()

	vd := vardatacandles.App{}
	vd.SetCurrentSession(1)

	vd.SetCandle(entity.Candle{1, []entity.Moment{
		entity.Moment{1, 1, entity.MomentPath{}},
		entity.Moment{1, 1, entity.MomentPath{}},
	}})
	app := filedatacandles.App{&vd, dir}
	vd.SetCandle(entity.Candle{2, []entity.Moment{
		entity.Moment{1, 1, entity.MomentPath{}},
	}})
	if len(app.DataSource.GetCurrentSessionCandles()) != 2 {
		t.Fatalf("error vardatacandles list")
	}

	app.CurrentSessionToFile()
	files, _ := os.ReadDir(dir)
	if len(files) != 1 {
		t.Fatalf("append file error, empty: %v", dir)
	}

	app.AllToDataSource()
	if len(app.DataSource.GetCurrentSessionCandles()) != 2 {
		t.Fatalf("error AllToDataSource")
	}
}
