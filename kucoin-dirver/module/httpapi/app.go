package httpapi

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type App struct {
	Port string
}

func (a *App) OpenHost(functions ...string) {
	router := httprouter.New()

	router.GET("/buymarket", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		fmt.Println("buy..")
	})

	router.GET("/setstop", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		fmt.Println("setstop..")
	})

	log.Println("open host to: http://127.0.0.1:" + a.Port)
	http.ListenAndServe(":"+a.Port, router)
}
