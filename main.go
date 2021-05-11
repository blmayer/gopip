package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/exec"

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
	packName := r.URL.Path[1:]

	path := "/tmp/" + packName
	err := os.Mkdir(path, os.ModeDir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("installing", packName, "into", path)
	cmd := exec.Command("pip", "install", "--isolated", packName, "-t", path)

	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("zipping files")
	cmd = exec.Command("zip", path+".zip", "-r", path)

	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		log.Println(err)
		return
	}

	// TODO: Clean pack files (to free memory)

	// Upload zip or save to mounted volume
	// The session the S3 Uploader will use
	sess := session.Must(session.NewSession())

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	f, err := os.Open(path + ".zip")
	if err != nil {
		log.Printf("failed to open %q, %v\n", path, err)
		return
	}

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("gopip"),
		Key:    aws.String(packName + ".zip"),
		Body:   f,
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		log.Printf("failed to upload file, %v\n", err)
		return
	}
	log.Printf("file uploaded to, %s\n", result.Location)

	w.Write([]byte(result.Location))
	cancel()
}
