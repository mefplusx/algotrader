package bootstrap

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	HTTPAPI_PORT     string   `yaml:"HTTPAPI_PORT"`
	HTTPAPI_HANDLERS []string `yaml:"HTTPAPI_HANDLERS"`

	REMOTEDATACONNECTOR_USE_CACHE bool     `yaml:"REMOTEDATACONNECTOR_USE_CACHE"`
	REMOTEDATACONNECTOR_HOST      string   `yaml:"REMOTEDATACONNECTOR_HOST"`
	REMOTEDATACONNECTOR_FINDNODES []string `yaml:"REMOTEDATACONNECTOR_FINDNODES"`

	ANALYTIC_FEE                    float64 `yaml:"ANALYTIC_FEE"`
	ANALYTIC_TAKE_LINE_BY_SLIDEDOWN float64 `yaml:"ANALYTIC_TAKE_LINE_BY_SLIDEDOWN"`
	ANALYTIC_PERC_IN                float64 `yaml:"ANALYTIC_PERC_IN"`
	ANALYTIC_TAKE_AFTER             float64 `yaml:"ANALYTIC_TAKE_AFTER"`
	ANALYTIC_SLIDEDOWN              float64 `yaml:"ANALYTIC_SLIDEDOWN"`
	ANALYTIC_STOP                   float64 `yaml:"ANALYTIC_STOP"`
	ANALYTIC_VALUE_USDT             float64 `yaml:"ANALYTIC_VALUE_USDT"`
	ANALYTIC_REPEAT                 int64   `yaml:"ANALYTIC_REPEAT"`
	ANALYTIC_LEVERAGE               int64   `yaml:"ANALYTIC_LEVERAGE"`
	ANALYTIC_ALLOW_SHORT            bool    `yaml:"ANALYTIC_ALLOW_SHORT"`
	ANALYTIC_ALLOW_LONG             bool    `yaml:"ANALYTIC_ALLOW_LONG"`

	FINDER_TAKE_LINE_BY_SLIDEDOWN []float64 `yaml:"FINDER_TAKE_LINE_BY_SLIDEDOWN"`
	FINDER_PERC_IN                []float64 `yaml:"FINDER_PERC_IN"`
	FINDER_TAKE_AFTER             []float64 `yaml:"FINDER_TAKE_AFTER"`
	FINDER_SLIDEDOWN              []float64 `yaml:"FINDER_SLIDEDOWN"`
	FINDER_STOP                   []float64 `yaml:"FINDER_STOP"`
	FINDER_VALUE_USDT             []float64 `yaml:"FINDER_VALUE_USDT"`
	FINDER_REPEAT                 []int64   `yaml:"FINDER_REPEAT"`
	FINDER_MIN_COUNT_TRANSACTIONS int64     `yaml:"FINDER_MIN_COUNT_TRANSACTIONS"`
	FINDER_TYPE                   string    `yaml:"FINDER_TYPE"`
	FINDER_RULES                  []string  `yaml:"FINDER_RULES"`
}

func (c *Config) Read(fileName string) {
	path, _ := os.Getwd()
	b, _ := os.ReadFile(path + "/" + fileName)
	yaml.Unmarshal(b, c)
}
