package main

import (
	"bytes"
	"encoding/json"
	"log"

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

	if r.Headers["Content-Type"] != "multipart/form-data" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid Content-Type",
		}, nil
	}

	// events.APIGatewayProxyRequest parse to SavSaveSoundRequest
	body := r.Body
	bytes := []byte(body)
	request := new(SaveSoundRequest)

	if err := json.Unmarshal(bytes, request); err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	log.Printf("success json.Unmarshal")

	if _, err := uploadFile(request.File); err != nil {
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

	return uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("meeting-sounds"),
		Key:    aws.String("hoge.wav"),
		Body:   bytes.NewReader([]byte(data)),
	})
}
