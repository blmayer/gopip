package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

var (
	ctx    context.Context
	cancel context.CancelFunc

	count = 0
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
	// waiting indefinitely for connections to return to idle and
	// then shut down.
	err := srv.Shutdown(context.Background())
	if err != nil {
		log.Println(err)
	}
}

func install(w http.ResponseWriter, r *http.Request) {
	// Two ways of getting the package name
	var pack string
	var err error

	if count > 0 {
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write(errorPage("Error", "Not first request!"))
		cancel()
		return
	}

	if r.Method == http.MethodGet {
		switch r.URL.Path {
		case "/favicon.ico":
			http.NotFound(w, r)
			return
		case "/", "/index.html":
			w.Write([]byte(index))
			return
		}

		// Get from URL
		pack, err = url.PathUnescape(r.URL.Path[1:])
	} else if r.Method == http.MethodPost {
		// Get from form
		err = r.ParseForm()
		if err != nil {
			log.Println("Parse Form", err)
		}
		pack = r.Form.Get("package")
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorPage("parse form error", err.Error()))
		return
	}

	args := []string{"install", "--isolated"}
	args = append(args, strings.Fields(pack)...)

	path := "/tmp/package"
	err = os.Mkdir(path, os.ModeDir)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorPage("", err.Error()))
		return
	}
	args = append(args, []string{"-t", path}...)

	// From now on the server is invalidated
	defer cancel()
	count++

	log.Println("running", args)
	cmd := exec.Command("pip", args...)

	err = cmd.Run()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorPage("pip install error", err.Error()))
		return
	}

	log.Println("zipping files")
	zipFile := "package.zip"

	cmd = exec.Command("zip", zipFile, "-r", path[5:])
	cmd.Dir = "/tmp"
	err = cmd.Run()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorPage("zip error", err.Error()))
		return
	}

	f, err := os.Open("tmp/" + zipFile)
	if err != nil {
		log.Printf("failed to open %q, %v\n", zipFile, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorPage("open zip failed", err.Error()))
		return
	}

	// Return file to user
	content, err := ioutil.ReadAll(f)
	if err != nil {
		log.Printf("failed to read %q, %v\n", zipFile, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorPage("read zip failed", err.Error()))
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=pack.zip")
	w.Write(content)
	log.Println("done")
}
