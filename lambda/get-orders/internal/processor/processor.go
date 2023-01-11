package processor

import (
	"ianag5j/ecommerce-back-go/get-orders/pkg/order"
	"ianag5j/ecommerce-back-go/get-orders/pkg/store"
	"ianag5j/ecommerce-back-go/get-orders/utils"

	"github.com/aws/aws-lambda-go/events"
)

type (
	Processor interface {
		Process(request events.APIGatewayProxyRequest) (response, error)
	}

	processor struct {
		s store.Database
		o order.Client
	}

	response struct {
		Orders []order.OrderModel `json:"order"`
	}
)

func New() Processor {
	return &processor{
		s: store.New(),
		o: order.New(),
	}
}

func (p processor) Process(request events.APIGatewayProxyRequest) (response, error) {
	r := response{}

	s, err := p.s.GetStoreByUser(utils.GetUserId(request.Headers["authorization"]))
	if err != nil {
		return r, err
	}

	o, err := p.o.GetOrders(s.Id)
	if err != nil {
		return r, err
	}

	r.Orders = o
	return r, err
}
