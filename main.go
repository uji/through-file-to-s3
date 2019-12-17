package main

import (
	"bytes"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// SaveSoundRequest is struct for parse request body
type SaveSoundRequest struct {
	File string `json:"file"`
}

// Handler is executed by AWS Lambda in the main function. Once the request
// is processed, it returns an Amazon API Gateway response object to AWS Lambda
func Handler(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// events.APIGatewayProxyRequest parse to SavSaveSoundRequest
	body := r.Body
	bytes := []byte(body)
	request := new(SaveSoundRequest)
	err := json.Unmarshal(bytes, request)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	uploadFile(request.File)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       request.File,
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

	return uploader.Upload(&s3manager.UploadInput{
		ACL:    aws.String("public-read"),
		Bucket: aws.String("test-bucket"),
		Key:    aws.String("hoge.png"),
		Body:   bytes.NewReader([]byte(data)),
	})
}
