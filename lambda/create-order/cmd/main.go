package main

import (
	"encoding/json"
	"ianag5j/ecommerce-back-go/create-order/internal/processor"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	p := processor.New()
	s := 200
	r, err := p.Process(request)
	o, _ := json.Marshal(r)
	if err != nil {
		s = 500
		o, _ = json.Marshal(ErrorResponse{Error: err.Error()})
	}

	return events.APIGatewayProxyResponse{
		Body:       string(o),
		Headers:    map[string]string{"Content-Type": "application/json"},
		StatusCode: s,
	}, nil
}

func main() {
	lambda.Start(handler)
}
