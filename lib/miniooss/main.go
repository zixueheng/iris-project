/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2024-01-02 16:16:54
 * @LastEditTime: 2024-01-04 10:47:28
 */
package main

import (
	"context"
	"io"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// https://min.io/docs/minio/linux/developers/go/minio-go.html

const (
	endpoint        = "127.0.0.1:9000"
	accessKeyID     = "FEsccxISzT1tgovzQR3t"
	secretAccessKey = "mRUz9a0wiIyNRFvgKQ6Kj3kr6HthjZFI2j46dvLg"
	useSSL          = false

	bucketName = "bucket1"
	location   = ""
)

var minioClient *minio.Client

func init() {
	var err error
	minioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()
	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}
}

func main() {
	// testUploadFile()

	// testGetObject()

	testGetUrl()
}

func testUploadFile() {
	ctx := context.Background()
	// Make a new bucket called testbucket.

	// Upload the test file
	// Change the value of filePath if the file is in another location
	objectName := "images/5ce8174dcc854.jpg"
	filePath := "C:/Users/root/Desktop/img/5ce8174dcc854.jpg"
	contentType := "image/jpeg" //"application/octet-stream"

	// Upload the test file with FPutObject
	info, err := minioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully created %+v\n", info)

	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
}

func testGetObject() {
	var (
		ctx        = context.Background()
		objectName = "images/5ce8174dcc854.jpg"
	)

	object, err := minioClient.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		log.Println(err)
		return
	}
	defer object.Close()

	localFile, err := os.Create("D:/go-projects/iris-project/lib/miniooss/5ce8174dcc854.jpg")
	if err != nil {
		log.Println(err)
		return
	}
	defer localFile.Close()

	if _, err = io.Copy(localFile, object); err != nil {
		log.Println(err)
		return
	}
}

func testGetUrl() {
	var (
		ctx        = context.Background()
		objectName = "images/5ce8174dcc854.jpg"
	)

	// Set request parameters for content-disposition.
	reqParams := make(url.Values)
	// reqParams.Set("response-content-disposition", "attachment; filename=\"your-filename.txt\"")
	// reqParams.Set("content-type", "image/jpeg")

	// Generates a presigned url which expires in a day.
	presignedURL, err := minioClient.PresignedGetObject(ctx, bucketName, objectName, time.Second*24*60*60, reqParams)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Successfully generated presigned URL", presignedURL)
}
