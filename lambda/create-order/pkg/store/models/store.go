package models

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type (
	Table interface {
		GetByName(name string) (Store, error)
	}

	table struct {
		DynamoDbClient *dynamodb.Client
		TableName      string
	}

	Store struct {
		Id     string `dynamodbav:"Id"`
		Name   string `dynamodbav:"Name"`
		UserId string `dynamodbav:"UserId"`
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
		TableName:      os.Getenv("STORES_TABLE"),
	}
}

func (table table) GetByName(name string) (Store, error) {

	store := Store{}
	output, err := table.DynamoDbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(table.TableName),
		Key: map[string]types.AttributeValue{
			"Name": &types.AttributeValueMemberS{Value: name},
		},
	})
	if err != nil {
		fmt.Println("error on get store: ", err.Error())
		return store, err
	}
	err = attributevalue.UnmarshalMap(output.Item, &store)
	if err != nil {
		fmt.Println("error on UnmarshalMap store: ", err.Error())
	}
	return store, err
}
