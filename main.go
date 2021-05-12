package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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
	// Two ways of getting the package name
	var pack string
	var err error

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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	args := append([]string{"install", "--isolated"}, strings.Fields(pack)...)

	path, err := os.MkdirTemp("", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	args = append(args, []string{"-t", path}...)

	// From now on the server is invalidated
	defer cancel()

	log.Println("running", args)
	cmd := exec.Command("pip", args...)
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		log.Println(err)
		http.Error(w, "pip install "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("zipping files")
	zipFile := path + ".zip"
	cmd = exec.Command("zip", zipFile, "-r", path)
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		log.Println(err)
		http.Error(w, "zip "+err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: Clean pack files (to free memory)
	f, err := os.Open(zipFile)
	if err != nil {
		log.Printf("failed to open %q, %v\n", zipFile, err)
		http.Error(w, "open zip "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("uploading")
	// Upload zip or save to mounted volume
	// The session the S3 Uploader will use
	sess := session.Must(session.NewSession())

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("gopip"),
		Key:    aws.String(strings.TrimPrefix(zipFile, "/tmp/")),
		Body:   f,
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		log.Printf("failed to upload file, %v\n", err)
		http.Error(w, "upload "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("file uploaded to %s\n", result.Location)

	w.Write(successPage(result.Location))
}
