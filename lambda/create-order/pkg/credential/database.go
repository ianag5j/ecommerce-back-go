package credential

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
	Database interface {
		Get(userId string, provider string) (Credential, error)
	}

	database struct {
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
		TableName:      os.Getenv("CREDENTIALS_TABLE"),
	}
}

func (d database) Get(userId string, provider string) (Credential, error) {

	credential := Credential{}
	output, err := d.DynamoDbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(d.TableName),
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
