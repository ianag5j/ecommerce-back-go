package store

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type (
	Database interface {
		GetStoreByUser(userId string) (Store, error)
	}

	database struct {
		DynamoDbClient *dynamodb.Client
		TableName      string
	}

	Store struct {
		Id     string `dynamodbav:"Id"`
		Name   string `dynamodbav:"Name"`
		UserId string `dynamodbav:"UserId"`
	}
)

func New() Database {
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
		TableName:      os.Getenv("STORES_TABLE"),
	}
}

func (d database) GetStoreByUser(userId string) (Store, error) {
	store := Store{}
	output, err := d.DynamoDbClient.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(d.TableName),
		IndexName:              aws.String("UserIdIndex"),
		KeyConditionExpression: aws.String("UserId = :userId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userId": &types.AttributeValueMemberS{Value: userId},
		},
		Select: types.SelectAllProjectedAttributes,
		Limit:  aws.Int32(1),
	})
	if err != nil {
		fmt.Println("error on get store: ", err.Error())
		return store, err
	}
	if len(output.Items) == 0 {
		return store, errors.New("error store not found")
	}

	err = attributevalue.UnmarshalMap(output.Items[0], &store)
	if err != nil {
		fmt.Println("error on UnmarshalMap store: ", err.Error())
	}
	return store, err
}
