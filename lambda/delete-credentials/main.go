package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"ianag5j/ecommerce-back-go/delete-credentials/utils"
)

type Response struct {
	Message string
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.Region = "us-east-1"
		return nil
	})
	if err != nil {
		panic(err)
	}

	svc := dynamodb.NewFromConfig(cfg)
	_, err = svc.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(os.Getenv("CREDENTIALS_TABLE")),
		Key: map[string]types.AttributeValue{
			"UserId":   &types.AttributeValueMemberS{Value: utils.GetUserId(request.Headers["authorization"])},
			"Provider": &types.AttributeValueMemberS{Value: "Uala"},
		},
	})

	res := Response{}

	if err != nil {
		fmt.Println(err)
		res.Message = err.Error()
		response, _ := json.Marshal(res)
		return events.APIGatewayProxyResponse{Body: string(response), StatusCode: 500}, nil
	}
	res.Message = "success"
	response, _ := json.Marshal(res)
	return events.APIGatewayProxyResponse{
		Body:       string(response),
		Headers:    map[string]string{"Content-Type": "application/json"},
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
