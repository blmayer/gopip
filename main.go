package main

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"os"
	"os/exec"
)

var (
	ctx    context.Context
	cancel context.CancelFunc
)

func main() {
	ctx, cancel = context.WithCancel(context.Background())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr: ":" + port,
	}

	http.HandleFunc("/", install)

	go func() {
		err := http.ListenAndServe(":"+port, nil)
		if err != http.ErrServerClosed {
			log.Println(err)
		}
	}()

	<-ctx.Done() // wait for the signal to gracefully shutdown the server

	// gracefully shutdown the server:
	// waiting indefinitely for connections to return to idle and then shut down.
	err := srv.Shutdown(context.Background())
	if err != nil {
		log.Println(err)
	}
}

func install(w http.ResponseWriter, r *http.Request) {
	packName := r.URL.Path[1:]

	path := "/tmp/" + packName
	err := os.Mkdir(path, os.ModeDir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("installing", packName, "into", path)
	cmd := exec.Command("pip", "install", "--isolated", packName, "-t", path)

	var pipOut bytes.Buffer
	cmd.Stdout = &pipOut

	err = cmd.Run()
	if err != nil {
		log.Println(err)
	}

	log.Println("zipping files")
	cmd = exec.Command("zip", path+".zip", "-r", path)
	var zipOut bytes.Buffer
	cmd.Stdout = &zipOut

	err = cmd.Run()
	if err != nil {
		log.Println(err)
	}

	// TODO: Clean pack files (to free memory)

	// Read zip and save

	// Upload zip or save to mounted volume

	w.Write(append(pipOut.Bytes(), zipOut.Bytes()...))
	cancel()
}
