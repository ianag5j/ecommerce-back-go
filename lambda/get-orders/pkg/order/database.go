package order

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type (
	Database interface {
		GetOrders(storeId string) ([]OrderModel, error)
	}

	database struct {
		DynamoDbClient *dynamodb.Client
		TableName      string
	}

	OrderModel struct {
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

func (t database) GetOrders(storeId string) ([]OrderModel, error) {
	o := []OrderModel{}
	r, err := t.DynamoDbClient.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(t.TableName),
		IndexName:              aws.String("StoreIdIndex"),
		KeyConditionExpression: aws.String("StoreId = :storeId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":storeId": &types.AttributeValueMemberS{Value: storeId},
		},
		Select: types.SelectAllProjectedAttributes,
		Limit:  aws.Int32(50),
	})
	if err != nil {
		return o, err
	}

	err = attributevalue.UnmarshalListOfMaps(r.Items, &o)
	return o, err
}
