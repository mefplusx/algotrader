package entity

type Moment struct {
	Value float64 `json:"v"`
	Close float64 `json:"c"`

	Path MomentPath
}
