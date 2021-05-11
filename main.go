package main

import (
	"fmt"
	"net/http"
	"os"
)

var count = 0

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", install)

	err := http.ListenAndServe(":" + port, nil)
	if err != nil {
		panic(err)
	}
}

func install(w http.ResponseWriter, r *http.Request) {
	count++
	out := fmt.Sprintln("request", count)
	w.Write([]byte(out))
}
