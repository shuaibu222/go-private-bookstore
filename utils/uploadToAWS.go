package utils

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/shuaibu222/go-bookstore/config"
)

var coverLinkCreated string
var bookLinkCreated string
var userLinkCreated string

func UploadBook(filepath string, resultChan chan<- string) {

	ext := path.Ext(filepath)

	log.Println(ext)

	if ext != ".pdf" && ext != ".txt" {
		log.Println("File type not supported")
		return
	}
	key := uuid.NewString() + ext
	log.Println(key)

	Region := "af-south-1"

	bucket := "onlinebooks"

	file, err := os.Open(filepath)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error opening file:", err)
		return
	}

	defer file.Close()

	svc, _ := config.AwsS3Config()

	// Read the contents of the file into a buffer
	var buf bytes.Buffer

	if _, err := io.Copy(&buf, file); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading file:", err)
		return
	}

	// This uploads the contents of the buffer to S3
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(buf.Bytes()),
	})

	if err != nil {
		fmt.Println("Error uploading file:", err)
		return
	}

	fmt.Println("File uploaded successfully!!!")

	// Generate the S3 file link based on your S3 bucket and key
	bookLinkCreated = fmt.Sprintf("https://%s.%s.s3.amazonaws.com/%s", bucket, Region, key)

	// If the upload is successful, send the path to the result channel
	if bookLinkCreated != "" {
		resultChan <- bookLinkCreated
	}

}

func UploadCoverImage(filepath string, resultChan chan<- string) {

	ext := path.Ext(filepath)

	log.Println(ext)

	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		log.Println("File type not supported")
		return
	}

	key := uuid.NewString() + ext
	log.Println(key)

	Region := "af-south-1"

	bucket := "onlinebooks"

	file, err := os.Open(filepath)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error opening file:", err)
		return
	}

	defer file.Close()

	svc, _ := config.AwsS3Config()

	// Read the contents of the file into a buffer
	var buf bytes.Buffer

	if _, err := io.Copy(&buf, file); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading file:", err)
		return
	}

	// This uploads the contents of the buffer to S3
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(buf.Bytes()),
	})

	if err != nil {
		fmt.Println("Error uploading file:", err)
		return
	}

	fmt.Println("File uploaded successfully!!!")

	// Generate the S3 file link based on your S3 bucket and key
	coverLinkCreated = fmt.Sprintf("https://%s.%s.s3.amazonaws.com/%s", bucket, Region, key)

	// If the upload is successful, send the path to the result channel
	if coverLinkCreated != "" {
		resultChan <- coverLinkCreated
	}

}

func UploadUserImage(filepath string) string {
	ext := path.Ext(filepath)

	log.Println(ext)

	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		log.Println("File type not supported")

	} else {
		key := uuid.NewString() + ext
		log.Println(key)

		Region := "af-south-1"

		bucket := "onlinebooks"

		file, err := os.Open(filepath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error opening file:", err)
		}
		defer file.Close()

		svc, _ := config.AwsS3Config()

		// Read the contents of the file into a buffer
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, file); err != nil {
			fmt.Fprintln(os.Stderr, "Error reading file:", err)
		}

		// This uploads the contents of the buffer to S3
		_, err = svc.PutObject(&s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
			Body:   bytes.NewReader(buf.Bytes()),
		})
		if err != nil {
			fmt.Println("Error uploading file:", err)
		}

		fmt.Println("File uploaded successfully!!!")

		// Generate the S3 file link based on your S3 bucket and key
		userLinkCreated = fmt.Sprintf("https://%s.%s.s3.amazonaws.com/%s", bucket, Region, key)

	}

	return userLinkCreated
}

// make it to be a controller for deleting an object
func DeleteFromS3(filepath string, key string) error {
	linkKey := key

	bucket := "onlinebooks"

	svc, _ := config.AwsS3Config()

	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(linkKey),
	})

	if err != nil {
		return err
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(linkKey),
	})

	if err != nil {
		log.Println("Your waited request failed")
		return err
	}

	log.Printf("Object %q successfully deleted\n", linkKey)

	return nil

}
