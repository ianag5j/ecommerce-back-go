package main

import (
	"ianag5j/ecommerce-back-go/get-store/pkg/handler"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	h := handler.NewHandler()
	lambda.Start(h.EventHandler)
}
