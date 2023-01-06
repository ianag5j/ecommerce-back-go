package processor

import (
	"encoding/json"
	order "ianag5j/ecommerce-back-go/create-order/pkg/order/models"
	store "ianag5j/ecommerce-back-go/create-order/pkg/store/models"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
)

type (
	Processor interface {
		Process(request events.APIGatewayProxyRequest) (response, error)
	}

	processor struct {
		store store.Table
		order order.Table
	}

	CartRequest struct {
		Id    string `json:"id"`
		Cant  int    `json:"cant"`
		Name  string `json:"name"`
		Price string `json:"price"`
	}

	BodyRequest struct {
		StoreName string        `json:"storeName"`
		Cart      []CartRequest `json:"cart"`
	}

	response struct {
		Message string      `json:"message,omitempty"`
		Order   order.Order `json:"order"`
	}
)

func New() Processor {
	store := store.Initialize()
	order := order.Initialize()
	return &processor{
		store: store,
		order: order,
	}
}

func (processor processor) Process(request events.APIGatewayProxyRequest) (response, error) {
	body := BodyRequest{}
	r := response{}
	json.Unmarshal([]byte(request.Body), &body)
	c, err := json.Marshal(body.Cart)
	if err != nil {
		r.Message = err.Error()
		return r, err
	}
	var amount float64
	for _, p := range body.Cart {
		pa, _ := strconv.ParseFloat(p.Price, 64)
		pc := float64(p.Cant)
		amount += pc * pa
	}

	store, err := processor.store.GetByName(body.StoreName)
	if err != nil {
		r.Message = err.Error()
		return r, err
	}

	o, err := processor.order.Create(amount, store.Id, "Uala", string(c))

	r.Order = o
	return r, err
}
