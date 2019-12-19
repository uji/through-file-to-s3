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
	request.Body = string(bytes)

	expectedResponse := events.APIGatewayProxyResponse{}

	response, err := Handler(request)
	t.Logf("Error: %s", err.Error())

	assert.Equal(t, response.Headers, expectedResponse.Headers)
	// assert.Contains(t, response.Body, expectedResponse.Body)
	// assert.Equal(t, err, nil)
}
