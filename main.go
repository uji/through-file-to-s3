package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

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

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       request.File,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
