package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/minio/minio-go"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Container struct {
	Image   string           `json:"image"`
	Name    string           `json:"name"`
	Command []string         `json:"command"`
	Ports   []map[string]int `json:"ports"`
}

type Pod struct {
	Kind       string                 `json:"kind"`
	ApiVersion string                 `json:"apiVersion"`
	Metadata   map[string]string      `json:"metadata"`
	Spec       map[string][]Container `json:"spec"`
}

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

	//TODO: adds a 6 character sha hash to the name so that files of same name don't get overwritten.
	objectName := fileName
	contentType := "application/octet-stream"

	// Upload the zip file with FPutObject
	n, err := minioClient.PutObject(bucketName, objectName, file, contentType)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, n)
	//Get binaryURL from minio for the object that we just uploaded
	url, err := minioClient.PresignedGetObject(bucketName, objectName, time.Hour, nil)

	cmdstr := createCommandString(url.String(), objectName)

	ports := []map[string]int{
		map[string]int{
			"hostPort":      8000,
			"containerPort": 8000,
		},
	}
	container := Container{"extrasalt/wgettu", "binary", []string{"sh", "-c", cmdstr}, ports}
	metadata := map[string]string{
		"name":      "goo",
		"namespace": "default",
	}
	pod := Pod{"Pod", "v1", metadata,
		map[string][]Container{"containers": []Container{container}}}

	var b []byte
	reader := bytes.NewBuffer(b)
	encoder := json.NewEncoder(reader)
	encoder.SetEscapeHTML(false)
	encoder.Encode(pod)

	req, err := http.NewRequest("POST", "http://localhost:8001/api/v1/namespaces/default/pods", reader)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	fmt.Println("%v", req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
	return url, nil

}

func createCommandString(url, filename string) string {

	//"wget -O /bin/#{filename} '#url' && chmod +x /bin/{#filename} && {#filename}"

	return fmt.Sprintf("wget -O /bin/%[2]s '%[1]s' && chmod +x /bin/%[2]s && %[2]s", url, filename)

}

func getShortHash(f io.Reader) string {

	hash := sha256.New()
	io.Copy(hash, f)
	key := hex.EncodeToString(hash.Sum(nil))

	return key[:6]

}
