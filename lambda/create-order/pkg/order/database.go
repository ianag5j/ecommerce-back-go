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
	Database interface {
		Save(o Order) (orderModel, error)
		Update(o Order) (orderModel, error)
		getOrderModel(o Order) (orderModel, error)
	}

	database struct {
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

func NewDatabase() Database {
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.Region = "us-east-1"
		return nil
	})
	if err != nil {
		panic(err)
	}

	svc := dynamodb.NewFromConfig(cfg)

	return &database{
		DynamoDbClient: svc,
		TableName:      os.Getenv("ORDERS_TABLE"),
	}
}

func (t database) Save(o Order) (orderModel, error) {
	order, err := t.getOrderModel(o)
	if err != nil {
		return order, err
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

func (t database) Update(o Order) (orderModel, error) {
	order, err := t.getOrderModel(o)
	if err != nil {
		return order, err
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

func (t database) getOrderModel(o Order) (orderModel, error) {
	order := orderModel{}
	shs, err := json.Marshal(o.StatusHistory)
	if err != nil {
		return order, err
	}
	c, err := json.Marshal(o.Cart)
	if err != nil {
		return order, err
	}

	order = orderModel{
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

	return order, err
}
