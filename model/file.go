package model

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spf13/viper"
)

type File struct {
	Id			int64		`json:"id" db:"id"`
	UserId		int64		`json:"user_id" db:"user_id"`
	Name		string		`json:"name" db:"name"`
	Hash		string		`json:"hash" db:"hash"`
}

type AwsConfig struct {
	AwsAccessKeyId			string		`mapstructure:"aws_access_key_id"`
	AwsSecretAccessKey		string		`mapstructure:"aws_secret_access_key"`
	AwsRegion				string		`mapstructure:"aws_region"`
	S3Bucket				string		`mapstructure:"s3_bucket"`
}

func UploadFileToS3(fileHash string, fileBytes []byte) (error) {
	// Load Config
	viperAwsConfig := viper.Sub("services.aws")
	awsConfig := AwsConfig{}
	err := viperAwsConfig.Unmarshal(&awsConfig)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
		return err
	}

	// Get AWS Session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsConfig.AwsRegion),
		Credentials: credentials.NewStaticCredentials(awsConfig.AwsAccessKeyId, awsConfig.AwsSecretAccessKey, "")})
	if err != nil {
		log.Fatal(err)
		return err
	}

	// Open the file for use
	file, err := ioutil.TempFile("/tmp", "file_upload")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer os.Remove(file.Name())
	err = ioutil.WriteFile(file.Name(), fileBytes, 0644)
	if err != nil { 
		return err
	}

	// Get file size and read the file content into a buffer
	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	// Config settings: this is where you choose the bucket, filename, content-type etc.
	// of the file you're uploading.
	_, err = s3.New(sess).PutObject(&s3.PutObjectInput{
		Bucket:					aws.String(awsConfig.S3Bucket),
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

func GetBytesOfFileFromS3(fileHash string) ([]byte, error) {
	// Load Config
	viperAwsConfig := viper.Sub("services.aws")
	awsConfig := AwsConfig{}
	err := viperAwsConfig.Unmarshal(&awsConfig)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
		return nil, err
	}

	item   := fmt.Sprintf("files/%v", fileHash)

	// AWS Download Session
	sess, _ := session.NewSession(&aws.Config{
			Region: aws.String(awsConfig.AwsRegion),
			Credentials: credentials.NewStaticCredentials(awsConfig.AwsAccessKeyId, awsConfig.AwsSecretAccessKey, "")})
	downloader := s3manager.NewDownloader(sess)

	// Download the file
	fileLocation := fmt.Sprintf("/tmp/%v", fileHash)
	file, err := os.Create(fileLocation)
	if err != nil {
		fmt.Printf("Unable to create file %q, %v", fileLocation, err)
		return nil, err
	}
	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket:		aws.String(awsConfig.S3Bucket),
			Key:		aws.String(item)})
	if err != nil {
		fmt.Printf("Unable to download item %q, %v", item, err)
		return nil, err
	}

	// Read File
	fileBytes, err := ioutil.ReadFile(fileLocation)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	return fileBytes, err
}
