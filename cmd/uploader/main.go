package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	s3Client *s3.S3
	s3Bucket string
	wg       sync.WaitGroup
)

func init() {
	acesskeyid := os.Getenv("AWS_ACCESS_KEY_ID")
	secretaccesskey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
		Credentials: credentials.NewStaticCredentials(
			acesskeyid,
			secretaccesskey,
			""),
	})
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}
	s3Client = s3.New(sess)
	s3Bucket = "your-bucket-name"
}

func main() {
	dir, err := os.Open("./tmp")
	if err != nil {
		log.Fatal(err)
	}
	defer dir.Close()

	uploadControl := make(chan struct{}, 5) // Limit to 5 concurrent uploads
	errorfileUpload := make(chan string, 5)

	go func() {
		for {
			select {
			case filename := <-errorfileUpload:
				uploadControl <- struct{}{} // Acquire a slot
				wg.Add(1)
				go uploadFile(filename, uploadControl, errorfileUpload)
			}
		}
	}()

	for {
		files, err := dir.Readdir(1)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
			continue
		}
		wg.Add(1)

		uploadControl <- struct{}{} // Acquire a slot

		go uploadFile(files[0].Name(), uploadControl, errorfileUpload)
	}
	wg.Wait()

}

func uploadFile(file string, uploadControl <-chan struct{}, errorfileUpload chan<- string) {
	defer wg.Done()

	completeFilePath := fmt.Sprintf("./tmp/%s", file)

	f, err := os.Open(completeFilePath)

	if err != nil {
		log.Fatal(err)
		<-uploadControl // Release the slot
		errorfileUpload <- completeFilePath
		return
	}

	defer f.Close()

	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(file),
		Body:   f,
	})
	if err != nil {
		log.Fatal(err)
		<-uploadControl // Release the slot
		errorfileUpload <- completeFilePath
		return

	}
	fmt.Printf("Successfully uploaded %s to %s\n", file, s3Bucket)
	<-uploadControl // Release the slot

}
