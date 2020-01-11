package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Handler is executed by AWS Lambda in the main function. Once the request
// is processed, it returns an Amazon API Gateway response object to AWS Lambda
func Handler(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Printf("Body: %s", r.Body)
	log.Print("Params: ", r.QueryStringParameters)

	contentType, params, err := mime.ParseMediaType(r.Headers["Content-Type"])
	if err != nil {
		return response(400), err
	}

	log.Printf("Content-Type: %s", contentType)
	log.Printf("boundary: %s", params["boundary"])
	log.Printf("IsBase64Encoded: %t", r.IsBase64Encoded)

	reader := multipart.NewReader(strings.NewReader(r.Body), params["boundary"])
	part, err := reader.NextPart()
	if err != nil {
		return response(400), err
	}

	buf, err := ioutil.ReadAll(part)
	if err != nil {
		return response(500), err
	}

	if _, err := uploadFile(buf); err != nil {
		return response(400), err
	}
	log.Printf("success uploadFile")

	return response(200), nil
}

func main() {
	lambda.Start(Handler)
}

func uploadFile(data []byte) (*s3manager.UploadOutput, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1"),
	}))
	uploader := s3manager.NewUploader(sess)

	t := time.Now().Format(time.RFC3339)

	return uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("meeting-sounds"),
		Key:    aws.String(t + ".wav"),
		Body:   bytes.NewReader(data),
	})
}

func response(statusCode int) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
		StatusCode: statusCode,
	}
}
