package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"
)

var (
	count  = 0
	ctx    context.Context
	cancel context.CancelFunc
)

func main() {
	ctx, cancel = context.WithCancel(context.Background())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", install)

	go func() {
		err := http.ListenAndServe(":"+port, nil)
		if err != http.ErrServerClosed {
			println(err)
		}
	}()

	<-ctx.Done() // wait for the signal to gracefully shutdown the server
}

func install(w http.ResponseWriter, r *http.Request) {
	count++
	out := fmt.Sprintln("request", count)
	time.Sleep(5 * time.Second)
	w.Write([]byte(out))
	cancel()
}
