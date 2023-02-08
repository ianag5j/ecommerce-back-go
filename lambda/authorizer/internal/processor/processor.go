package processor

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"ianag5j/ecommerce-back-go/authorizer/pkg/jwk"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type (
	Processor interface {
		Process(e events.APIGatewayV2CustomAuthorizerV2Request) (events.APIGatewayV2CustomAuthorizerSimpleResponse, error)
	}

	processor struct {
	}
)

func New() Processor {
	return &processor{}
}

func (p processor) Process(e events.APIGatewayV2CustomAuthorizerV2Request) (events.APIGatewayV2CustomAuthorizerSimpleResponse, error) {
	j, _ := json.Marshal(e)
	fmt.Println(bytes.NewBuffer(j).String())

	if e.Type != "REQUEST" {
		return events.APIGatewayV2CustomAuthorizerSimpleResponse{
			IsAuthorized: false,
		}, errors.New(`expected "event.type" parameter to have value "REQUEST"`)
	}

	ats := strings.Split(e.Headers["authorization"], " ")[1]
	at, err := jwk.ValidateToken(ats)

	if err != nil {
		fmt.Println(err)
		return events.APIGatewayV2CustomAuthorizerSimpleResponse{
			IsAuthorized: false,
		}, nil
	}

	userId, _ := at.Get("sub")

	return events.APIGatewayV2CustomAuthorizerSimpleResponse{
		IsAuthorized: true,
		Context: map[string]interface{}{
			"userId": userId,
		},
	}, nil
}
