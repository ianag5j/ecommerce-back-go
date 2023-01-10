package processor

import (
	"encoding/json"
	uala "ianag5j/ecommerce-back-go/create-order/pkg/clients"
	"ianag5j/ecommerce-back-go/create-order/pkg/credential"
	"ianag5j/ecommerce-back-go/create-order/pkg/order"
	"ianag5j/ecommerce-back-go/create-order/pkg/store"

	"github.com/aws/aws-lambda-go/events"
)

type (
	Processor interface {
		Process(request events.APIGatewayProxyRequest) (response, error)
	}

	processor struct {
		s store.Database
		o order.Client
		c credential.Database
	}

	BodyRequest struct {
		StoreName string              `json:"storeName"`
		Cart      []order.CartRequest `json:"cart"`
	}

	response struct {
		Order     *order.Order   `json:"order"`
		UalaOrder uala.UalaOrder `json:"ualaOrder"`
	}
)

func New() Processor {
	return &processor{
		s: store.New(),
		o: order.New(),
		c: credential.New(),
	}
}

func (p processor) Process(request events.APIGatewayProxyRequest) (response, error) {
	body := BodyRequest{}
	r := response{}
	json.Unmarshal([]byte(request.Body), &body)

	s, err := p.s.GetByName(body.StoreName)
	if err != nil {
		return r, err
	}

	o, err := p.o.Create(body.Cart, s.Id, "Uala")
	if err != nil {
		return r, err
	}
	r.Order = &o

	c, err := p.c.Get(s.UserId, "Uala")
	if err != nil {
		return r, err
	}

	u := uala.New(c)
	uo, err := u.CreateOrder(o, s)
	if err != nil {
		return r, err
	}

	r.UalaOrder = uo
	o.ExternalId = uo.Uuid

	err = p.o.UpdateStatus(&o, "PENDING")
	return r, err
}
