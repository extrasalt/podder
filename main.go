package main

import (
	"github.com/gorilla/mux"
	"github.com/minio/minio-go"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

var minioClient *minio.Client
var err error

func main() {

	//Test bed values. Replace with real minio address and keys
	endpoint := "play.minio.io:9000"
	accessKeyID := "Q3AM3UQ867SPQQA43P2F"
	secretAccessKey := "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG"
	useSSL := true

	minioClient, err = minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Fatalln(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/getbinary", UserBinaryHandler)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	http.ListenAndServe(":3000", r)
}

func UserBinaryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	binary, header, err := r.FormFile("upload")

	if err != nil {
		panic(err)
	}

	url, err := uploadFile(header.Filename, binary)

	w.Write([]byte(url.String()))
}

func uploadFile(fileName string, file io.Reader) (*url.URL, error) {

	bucketName := "binary"
	location := "us-east-1" //As given in docs. Might change when we use our own server

	err = minioClient.MakeBucket(bucketName, location)
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, err := minioClient.BucketExists(bucketName)
		if err == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	}
	log.Printf("Successfully created %s\n", bucketName)

	objectName := fileName
	contentType := "application/octet-stream"

	// Upload the zip file with FPutObject
	n, err := minioClient.PutObject(bucketName, objectName, file, contentType)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, n)
	//Get binaryURL from minio for the object that we just uploaded
	url, err := minioClient.PresignedGetObject(bucketName, objectName, time.Minute, nil)

	return url, nil

}
