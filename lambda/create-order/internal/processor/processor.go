package processor

import (
	"encoding/json"
	uala "ianag5j/ecommerce-back-go/create-order/pkg/clients"
	credential "ianag5j/ecommerce-back-go/create-order/pkg/credential/models"
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
		store      store.Table
		order      order.Table
		credential credential.Table
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
		Message   string         `json:"message,omitempty"`
		Order     order.Order    `json:"order"`
		UalaOrder uala.UalaOrder `json:"ualaOrder"`
	}
)

func New() Processor {
	store := store.Initialize()
	order := order.Initialize()
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
	cart, err := json.Marshal(body.Cart)
	if err != nil {
		r.Message = err.Error()
		return r, err
	}

	s, err := p.store.GetByName(body.StoreName)
	if err != nil {
		r.Message = err.Error()
		return r, err
	}

	//TODO: validate amount with products amounts in database
	var amount float64
	for _, p := range body.Cart {
		pa, _ := strconv.ParseFloat(p.Price, 64)
		pc := float64(p.Cant)
		amount += pc * pa
	}
	o, err := p.order.Create(amount, s.Id, "Uala", string(cart))
	r.Order = o

	c, err := p.credential.Get(s.UserId, "Uala")

	u := uala.New(c)
	uo, err := u.CreateOrder(o, s)
	if err == nil {
		r.UalaOrder = uo
	}
	return r, err
}
