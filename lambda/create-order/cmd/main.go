package main

import (
	"ianag5j/ecommerce-back-go/create-order/internal/processor"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	p := processor.New()
	p.Process(request)
	return events.APIGatewayProxyResponse{
		// Body:       string(response),
		Headers:    map[string]string{"Content-Type": "application/json"},
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
