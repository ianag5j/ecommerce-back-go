package services

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func setup() dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.Region = "us-east-1"
		return nil
	})
	if err != nil {
		panic(err)
	}

	svc := dynamodb.NewFromConfig(cfg)
	return *svc
}

func UpdateOrderStatus(orderId string, status string) bool {
	svc := setup()

	_, err := svc.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String(os.Getenv("ORDERS_TABLE")),
		Key: map[string]types.AttributeValue{
			"Id": &types.AttributeValueMemberS{Value: orderId},
		},
		UpdateExpression: aws.String("set #s = :status, #u = :updatedAt"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":status":    &types.AttributeValueMemberS{Value: status},
			":updatedAt": &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},
		},
		ExpressionAttributeNames: map[string]string{
			"#s": "Status",
			"#u": "UpdatedAt",
		},
	})
	fmt.Println(err)
	return err == nil
}
