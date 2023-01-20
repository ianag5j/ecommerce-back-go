package handler

import (
	"ianag5j/ecommerce-back-go/authorizer/internal/processor"

	"github.com/aws/aws-lambda-go/events"
)

type (
	Handler struct {
		p processor.Processor
	}
)

func NewHandler() Handler {
	return Handler{
		p: processor.New(),
	}
}

func (h Handler) EventHandler(e events.APIGatewayV2CustomAuthorizerV2Request) (events.APIGatewayV2CustomAuthorizerSimpleResponse, error) {
	return h.p.Process(e)
}
