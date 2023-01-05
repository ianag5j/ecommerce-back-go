package processor

import (
	"encoding/json"
	"fmt"
	"ianag5j/ecommerce-back-go/create-order/pkg/order/models"
	store "ianag5j/ecommerce-back-go/create-order/pkg/store/models"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
)

type (
	BodyRequest struct {
		StoreName string            `json:"storeName"`
		Cart      map[string]string `json:"cart"`
	}

	Processor interface {
		Process(request events.APIGatewayProxyRequest)
	}

	processor struct {
		store *store.Table
	}
)

func New() Processor {
	store := store.Initialize()
	return &processor{
		store: &store,
	}
}

func (processor processor) Process(request events.APIGatewayProxyRequest) {
	sd := store.Initialize()

	body := BodyRequest{}
	json.Unmarshal([]byte(request.Body), &body)

	fmt.Println(body.Cart)
	store := sd.GetByName(body.StoreName)
	fmt.Println(store.Id, store.Name)
	order := models.Order{
		Id:      uuid.NewString(),
		StoreId: store.Id,
		// Amount: ,
		PaymentMethod: "Uala",
		// Cart: ,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}
	fmt.Println(order)
	// d := database.New()
	// d.AddOrder(order)
}
