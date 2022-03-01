package helper

import (
	"fmt"
	"net/http"
)

func RunPprof() {
	go func() {
		fmt.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	fmt.Println("go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30")
	fmt.Println("go tool pprof -alloc_space http://localhost:8080/debug/pprof/heap\n")
}
