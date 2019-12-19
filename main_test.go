package main

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {

	bytes, err := json.Marshal(SaveSoundRequest{
		File: "sound",
	})

	request := events.APIGatewayProxyRequest{}
	request.Headers = map[string]string{"Content-Type": "application/json"}
	request.Body = string(bytes)

	expectedResponse := events.APIGatewayProxyResponse{}
	expectedResponse.StatusCode = 400
	response, err := Handler(request)
	assert.Equal(t, response.Headers, expectedResponse.Headers)

	request.Headers = map[string]string{"Content-Type": "multipart/form-data"}
	expectedResponse.StatusCode = 200

	response, err = Handler(request)
	assert.Equal(t, response.Headers, expectedResponse.Headers)
	assert.Equal(t, err, nil)
}
