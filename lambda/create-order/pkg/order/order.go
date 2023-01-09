package order

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type (
	Order struct {
		Id            string          `json:"id"`
		StoreId       string          `json:"storeId"`
		Amount        float64         `json:"amount"`
		Status        string          `json:"status"`
		StatusHistory []statusHistory `json:"statusHistory"`
		Cart          []cartRequest   `json:"cart"`
		PaymentMethod string          `json:"paymentMethod"`
		CreatedAt     string          `json:"createdAt"`
		UpdatedAt     string          `json:"updatedAt"`
	}

	statusHistory struct {
		CreatedAt string `json:"createdAt"`
		Status    string `json:"status"`
	}

	cartRequest struct {
		Id    string `json:"id"`
		Cant  int    `json:"cant"`
		Name  string `json:"name"`
		Price string `json:"price"`
	}
)

func Create(amount float64, storeId string, paymentMethod string, cart string) (Order, error) {
	sh := []statusHistory{{Status: "CREATED", CreatedAt: time.Now().Format(time.RFC3339)}}
	c := []cartRequest{}
	json.Unmarshal([]byte(cart), &c)
	o := Order{
		Id:            uuid.NewString(),
		StoreId:       storeId,
		Amount:        amount,
		Status:        "CREATED",
		StatusHistory: sh,
		PaymentMethod: paymentMethod,
		Cart:          c,
		CreatedAt:     time.Now().Format(time.RFC3339),
		UpdatedAt:     time.Now().Format(time.RFC3339),
	}

	return o, nil
}
