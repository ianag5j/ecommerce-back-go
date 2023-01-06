package models

import (
	"context"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
)

type (
	Order struct {
		Id            string  `dynamodbav:"Id"`
		StoreId       string  `dynamodbav:"StoreId"`
		Amount        float64 `dynamodbav:"Amount"`
		Status        string  `dynamodbav:"Status"`
		Cart          string  `dynamodbav:"Cart"`
		PaymentMethod string  `dynamodbav:"PaymentMethod"`
		CreatedAt     string  `dynamodbav:"CreatedAt"`
		UpdatedAt     string  `dynamodbav:"UpdatedAt"`
	}

	CartRequest struct {
		Id    string `json:"id"`
		Cant  int    `json:"cant"`
		Name  string `json:"name"`
		Price string `json:"price"`
	}

	Table interface {
		Create(amount float64, storeId string, paymentMethod string, cart string) (Order, error)
	}

	table struct {
		DynamoDbClient *dynamodb.Client
		TableName      string
	}
)

func Initialize() Table {
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

func (table table) Create(amount float64, storeId string, paymentMethod string, cart string) (Order, error) {
	order := Order{
		Id:            uuid.NewString(),
		StoreId:       storeId,
		Amount:        amount,
		Status:        "PENDING",
		PaymentMethod: paymentMethod,
		Cart:          cart,
		CreatedAt:     time.Now().Format(time.RFC3339),
		UpdatedAt:     time.Now().Format(time.RFC3339),
	}

	item, err := attributevalue.MarshalMap(order)
	if err != nil {
		return order, err
	}

	_, err = table.DynamoDbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(table.TableName), Item: item,
	})
	return order, err
}
