package model

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/viper"
)

type File struct {
	Id			int64		`json:"id" db:"id"`
	UserId		int64		`json:"user_id" db:"user_id"`
	Name		string		`json:"name" db:"name"`
	Hash		string		`json:"hash" db:"hash"`
}

type S3Config struct {
	AccessKeyId			string		`json:"AWS_ACCESS_KEY_ID"`
	SecretAccessKey		string		`json:"AWS_SECRET_ACCESS_KEY"`
	Bucket				string		`json:"AWS_S3_BUCKET"`
	Region				string		`json:"AWS_S3_REGION"`
}

func UploadFileToS3(fileHash string, fileBytes []byte) (error) {
	// Load Config
	s3Config := S3Config{}
	viper.SetConfigName("aws_s3") // name of config file (without extension)
	viper.AddConfigPath("../config/") // path to look for the config file in
	err := viper.MergeInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	err = viper.Unmarshal(&s3Config)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
	}

	sess, err := session.NewSession(&aws.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Open the file for use
	file, err := ioutil.TempFile("/tmp", "file_upload")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())
	err = ioutil.WriteFile(file.Name(), fileBytes, 0644)
	if err != nil { 
		// handle error
	}

	// Get file size and read the file content into a buffer
	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	// Config settings: this is where you choose the bucket, filename, content-type etc.
	// of the file you're uploading.
	_, err = s3.New(sess).PutObject(&s3.PutObjectInput{
		Bucket:					aws.String("basoukazuma-file-upload"),
		Key:					aws.String(fmt.Sprintf("files/%v", fileHash)),
		ACL:					aws.String("private"),
		Body:					bytes.NewReader(buffer),
		ContentLength:			aws.Int64(size),
		ContentType:			aws.String(http.DetectContentType(buffer)),
		ContentDisposition:		aws.String("attachment"),
		ServerSideEncryption:	aws.String("AES256"),
	})
	return err
}
