package main

import (
	"bytes"
	"log"
	"mime"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
)

// Handler is executed by AWS Lambda in the main function. Once the request
// is processed, it returns an Amazon API Gateway response object to AWS Lambda
func Handler(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	contentType, params, err := mime.ParseMediaType(r.Headers["Content-Type"])
	log.Printf("Content-Type: %s", contentType)
	log.Printf("boundary: %s", params["boundary"])
	if err != nil || contentType != "multipart/form-data" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
		}, err
	}

	log.Printf("IsBase64Encoded: %t", r.IsBase64Encoded)

	accessor := core.RequestAccessor{}
	request, err := accessor.ProxyEventToHTTPRequest(r)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	log.Printf("success EventToRequest")
	log.Printf("request body: %v", request.Body)

	if err := request.ParseMultipartForm(4096); err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	log.Printf("success ParseMultipartForm")

	if _, err := uploadFile(request.FormValue("file")); err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	log.Printf("success uploadFile")

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}

func uploadFile(data string) (*s3manager.UploadOutput, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1"),
	}))
	uploader := s3manager.NewUploader(sess)

	t := time.Now().Format(time.RFC3339)

	return uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("meeting-sounds"),
		Key:    aws.String(t + ".wav"),
		Body:   bytes.NewReader([]byte(data)),
	})
}
