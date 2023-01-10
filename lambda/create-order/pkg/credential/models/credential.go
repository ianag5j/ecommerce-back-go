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
		Get(userId string, provider string) (Credential, error)
	}

	table struct {
		DynamoDbClient *dynamodb.Client
		TableName      string
	}

	Credential struct {
		UserId               string `dynamodbav:"UserId"`
		Provider             string `dynamodbav:"Provider"`
		ExternalClientId     string `dynamodbav:"externalClientId"`
		ExternalClientSecret string `dynamodbav:"externalClientSecret"`
		ExternalUserName     string `dynamodbav:"externalUserName"`
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
		TableName:      os.Getenv("CREDENTIALS_TABLE"),
	}
}

func (table table) Get(userId string, provider string) (Credential, error) {

	credential := Credential{}
	output, err := table.DynamoDbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(table.TableName),
		Key: map[string]types.AttributeValue{
			"UserId":   &types.AttributeValueMemberS{Value: userId},
			"Provider": &types.AttributeValueMemberS{Value: provider},
		},
	})
	if err != nil {
		fmt.Println("error on get credential: ", err.Error())
		return credential, err
	}
	err = attributevalue.UnmarshalMap(output.Item, &credential)
	if err != nil {
		fmt.Println("error on UnmarshalMap credential: ", err.Error())
	}
	return credential, err
}
