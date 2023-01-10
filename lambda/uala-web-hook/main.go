package main

import (
	"encoding/json"
	"fmt"

	ds "ianag5j/ecommerce-back-go/uala-web-hook/services/dynamodb"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type (
	BodyRequest struct {
		Status string `json:"status"`
	}

	Response struct {
		Message string
	}

	ResponseError struct {
		Message string
		Error   string
	}
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println(request.Body)
	fmt.Println(request.PathParameters)
	body := BodyRequest{}
	json.Unmarshal([]byte(request.Body), &body)

	dynamodbService := ds.New()
	isUpdated, err := dynamodbService.UpdateOrderStatus(request.PathParameters["orderId"], body.Status)

	responseStatus := 200
	response, _ := json.Marshal(Response{
		Message: "Success",
	})
	if !isUpdated {
		res := ResponseError{
			Message: "Error on update order",
			Error:   err.Error(),
		}
		response, _ = json.Marshal(res)
		responseStatus = 500
	}

	return events.APIGatewayProxyResponse{
		Body: string(response),
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Headers": "*",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "OPTIONS,POST",
		},
		StatusCode: responseStatus,
	}, nil
}

func main() {
	lambda.Start(handler)
}
