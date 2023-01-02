package main

import (
	"encoding/json"
	"fmt"

	ds "ianag5j/ecommerce-back-go/uala-web-hook/services/dynamodb"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type BodyRequest struct {
	Status string `json:"status"`
}

type Response struct {
	Message string
	Error   string
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println(request.Body)
	fmt.Println(request.PathParameters)
	body := BodyRequest{}
	json.Unmarshal([]byte(request.Body), &body)

	dynamodbService := ds.New()
	isUpdated, err := dynamodbService.UpdateOrderStatus(request.PathParameters["orderId"], body.Status)
	res := Response{
		Message: "Success",
	}
	responseStatus := 200

	if !isUpdated {
		res = Response{
			Message: "Error on update order",
			Error:   err.Error(),
		}
		responseStatus = 500
	}
	response, _ := json.Marshal(res)
	return events.APIGatewayProxyResponse{
		Body:       string(response),
		Headers:    map[string]string{"Content-Type": "application/json"},
		StatusCode: responseStatus,
	}, nil
}

func main() {
	lambda.Start(handler)
}
