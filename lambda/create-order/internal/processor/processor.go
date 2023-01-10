package processor

import (
	"encoding/json"
	"errors"
	uala "ianag5j/ecommerce-back-go/create-order/pkg/clients"
	credential "ianag5j/ecommerce-back-go/create-order/pkg/credential/models"
	"ianag5j/ecommerce-back-go/create-order/pkg/order"
	store "ianag5j/ecommerce-back-go/create-order/pkg/store/models"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
)

type (
	Processor interface {
		Process(request events.APIGatewayProxyRequest) (response, error)
	}

	processor struct {
		store      store.Table
		order      order.Table
		credential credential.Table
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
	store := store.Initialize()
	order := order.New()
	credential := credential.Initialize()
	return &processor{
		store:      store,
		order:      order,
		credential: credential,
	}
}

func (p processor) Process(request events.APIGatewayProxyRequest) (response, error) {
	body := BodyRequest{}
	r := response{}
	json.Unmarshal([]byte(request.Body), &body)

	s, err := p.store.GetByName(body.StoreName)
	if err != nil {
		return r, err
	}

	ta, err := getTotalAmount(body)
	if err != nil {
		return r, err
	}

	o, err := order.Create(ta, s.Id, "Uala", body.Cart)
	if err != nil {
		return r, err
	}
	r.Order = &o

	_, err = p.order.Save(o)
	if err != nil {
		return r, err
	}

	c, err := p.credential.Get(s.UserId, "Uala")
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
	order.UpdateStatus(&o, "PENDING")
	_, err = p.order.Update(o)

	return r, err
}

func getTotalAmount(b BodyRequest) (float64, error) {
	ta := 0.0

	//TODO: validate amount with products amounts in database
	for _, p := range b.Cart {
		pa, err := strconv.ParseFloat(p.Price, 64)
		if err != nil {
			return ta, errors.New("error on parse amount")
		}
		pc := float64(p.Cant)
		ta += pc * pa
	}
	return ta, nil
}
