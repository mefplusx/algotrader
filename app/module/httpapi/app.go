package httpapi

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"robot/module/httpapi/interfaces"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type App struct {
	DataSource interfaces.DataSource
	Port       string

	router *httprouter.Router
}

func (a *App) OpenHost(functions ...string) {
	a.router = httprouter.New()

	for _, function := range functions {
		reflect.ValueOf(a).MethodByName(function).Call([]reflect.Value{})
		log.Println("start handler:", function)
	}

	log.Println("open host to: http://127.0.0.1:" + a.Port)
	http.ListenAndServe(":"+a.Port, a.router)
}

func (a *App) DATA() {
	fs := http.FileServer(http.Dir("./module/filedatacandles/.data"))
	http.Handle("/data/", http.StripPrefix("/data/", fs))
}

func (a *App) LOG() {
	a.router.GET("/log/:file", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		file := strings.ReplaceAll(ps.ByName("file"), ".", "") + ".log"
		b, _ := os.ReadFile(file)
		fmt.Fprintf(w, string(b))
	})
}

func (a *App) SETCANDLE() {
	a.router.POST("/set-candle/:currency/:session", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		currency := ps.ByName("currency")
		session, _ := strconv.ParseInt(ps.ByName("session"), 10, 64)

		fmt.Println(currency, session)

		fmt.Fprintf(w, string("{\"success\": true}"))
	})
}
