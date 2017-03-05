package main

import (
	"bytes"
	"io"
	"log"
	"net/url"
	"time"
)

func uploadFile(fileName string, file io.Reader) (*url.URL, string, error) {

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

	var b []byte
	buf := bytes.NewBuffer(b)

	fileCopy := io.TeeReader(file, buf)

	//Adds a 6 character sha hash to the name so that files of same name don't get overwritten.
	objectName := fileName + "-" + getShortHash(fileCopy)
	contentType := "application/octet-stream"

	//Upload the file in buffer to minio
	n, err := minioClient.PutObject(bucketName, objectName, buf, contentType)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, n)

	//Get binaryURL from minio for the object that we just uploaded
	url, err := minioClient.PresignedGetObject(bucketName, objectName, time.Hour, nil)

	return url, objectName, nil

}
