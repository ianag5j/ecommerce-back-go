package store

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
		CreateStore(s Store) (Store, error)
		GetStoreByName(storeName string) (Store, error)
		GetStoreByUserId(userId string) (Store, error)
	}

	database struct {
		DynamoDbClient *dynamodb.Client
		TableName      string
	}

	StoreModel struct {
		Id     string `dynamodbav:"Id"`
		Name   string `dynamodbav:"Name"`
		UserId string `dynamodbav:"UserId"`
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
		TableName:      os.Getenv("STORES_TABLE"),
	}
}

func (d database) CreateStore(s Store) (Store, error) {
	item, err := attributevalue.MarshalMap(StoreModel{
		Id:     s.Id,
		Name:   s.Name,
		UserId: s.UserId,
	})

	if err != nil {
		return s, err
	}

	_, err = d.DynamoDbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: &d.TableName,
		Item:      item,
	})

	return s, err
}

func (d database) GetStoreByName(name string) (Store, error) {
	s := Store{}
	o, err := d.DynamoDbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: &d.TableName,
		Key: map[string]types.AttributeValue{
			"Name": &types.AttributeValueMemberS{Value: name},
		},
	})
	if err != nil {
		fmt.Println("error on get store: ", err.Error())
		return s, err
	}

	err = attributevalue.UnmarshalMap(o.Item, &s)
	if err != nil {
		fmt.Println("error on UnmarshalMap store: ", err.Error())
	}
	return s, err
}

func (d database) GetStoreByUserId(userId string) (Store, error) {
	s := Store{}
	o, err := d.DynamoDbClient.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              &d.TableName,
		IndexName:              aws.String("UserIdIndex"),
		KeyConditionExpression: aws.String("UserId = :ui"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":ui": &types.AttributeValueMemberS{Value: userId},
		},
		Select: types.SelectAllProjectedAttributes,
		Limit:  aws.Int32(1),
	})

	if err != nil {
		fmt.Println("error on get store: ", err.Error())
		return s, err
	}

	if len(o.Items) == 0 {
		return s, nil
	}

	stores := []Store{}
	err = attributevalue.UnmarshalListOfMaps(o.Items, &stores)
	if err != nil {
		fmt.Println("error on UnmarshalMap store: ", err.Error())
		return s, err
	}

	return stores[0], err
}
