package order

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type (
	Table interface {
		Save(o Order) (orderModel, error)
	}

	table struct {
		DynamoDbClient *dynamodb.Client
		TableName      string
	}

	orderModel struct {
		Id            string  `dynamodbav:"Id"`
		StoreId       string  `dynamodbav:"StoreId"`
		Amount        float64 `dynamodbav:"Amount"`
		Status        string  `dynamodbav:"Status"`
		StatusHistory string  `dynamodbav:"StatusHistory"`
		Cart          string  `dynamodbav:"Cart"`
		PaymentMethod string  `dynamodbav:"PaymentMethod"`
		CreatedAt     string  `dynamodbav:"CreatedAt"`
		UpdatedAt     string  `dynamodbav:"UpdatedAt"`
	}
)

func New() Table {
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.Region = "us-east-1"
		return nil
	})
	if err != nil {
		panic(err)
	}

	svc := dynamodb.NewFromConfig(cfg)

	return &table{
		DynamoDbClient: svc,
		TableName:      os.Getenv("ORDERS_TABLE"),
	}
}

func (t table) Save(o Order) (orderModel, error) {
	shs, _ := json.Marshal(o.StatusHistory)
	c, _ := json.Marshal(o.Cart)
	order := orderModel{
		Id:            o.Id,
		StoreId:       o.StoreId,
		Amount:        o.Amount,
		Status:        o.Status,
		StatusHistory: string(shs),
		PaymentMethod: o.PaymentMethod,
		Cart:          string(c),
		CreatedAt:     o.CreatedAt,
		UpdatedAt:     o.UpdatedAt,
	}

	item, err := attributevalue.MarshalMap(order)
	if err != nil {
		return order, err
	}

	_, err = t.DynamoDbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(t.TableName), Item: item,
	})
	return order, err
}
