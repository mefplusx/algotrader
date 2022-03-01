package internal

import (
	"bytes"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Sender struct {
	Url  string
	Body interface{}
}

func (s *Sender) Send() bool {
	jsonData, err := json.Marshal(s.Body)
	if err != nil {
		return true
	}

	_, err = http.Post(s.Url, "application/json", bytes.NewBuffer(jsonData))
	if err == nil {
		return true
	}

	return false
}
