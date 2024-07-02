package main

import (
	"fmt"
	"net/http"
	"time"
)

func VerySlowFunction() {
	time.Sleep(6 * time.Second)
}

func Process(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	go VerySlowFunction()

	select {
	case <-time.After(10 * time.Second):
		w.Write([]byte("OK"))
	case <-ctx.Done():
		fmt.Println("Cancelled")
	}

}

func main() {
	http.ListenAndServe(":8080", http.HandlerFunc(Process))
}
