// package database

// import (
// 	"context"
// 	"ianag5j/ecommerce-back-go/create-order/pkg/order/models"
// 	"os"

// 	"github.com/aws/aws-sdk-go-v2/aws"
// 	"github.com/aws/aws-sdk-go-v2/config"
// 	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
// 	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
// )

// type Table interface {
// 	AddOrder(order models.Order)
// }

// type table struct {
// 	DynamoDbClient *dynamodb.Client
// 	TableName      string
// }

// func New() Table {
// 	cfg, err := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
// 		o.Region = "us-east-1"
// 		return nil
// 	})
// 	if err != nil {
// 		panic(err)
// 	}

// 	svc := dynamodb.NewFromConfig(cfg)

// 	return &table{
// 		DynamoDbClient: svc,
// 		TableName:      os.Getenv("ORDERS_TABLE"),
// 	}
// }

// func (table table) AddOrder(order models.Order) {
// 	item, err := attributevalue.MarshalMap(order)
// 	if err != nil {
// 		panic("error marshal order")
// 	}

// 	_, err = table.DynamoDbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
// 		TableName: aws.String(table.TableName), Item: item,
// 	})

// 	if err != nil {
// 		panic("error on save order")
// 	}
// }
