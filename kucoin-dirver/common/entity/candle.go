package entity

type Candle struct {
	Time    int64    `json:"t"`
	Moments []Moment `json:"m"`
}
