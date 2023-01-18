package handler

import (
	"encoding/json"
	"ianag5j/ecommerce-back-go/get-store/internal/processor"

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

func (h Handler) EventHandler(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	rb := h.p.Process(r)

	response, _ := json.Marshal(rb)

	return events.APIGatewayProxyResponse{
		Body: string(response),
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Headers": "*",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "OPTIONS,POST",
		},
		StatusCode: rb.StatusCode,
	}, nil
}
