package order

import (
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
		Cart          []CartRequest   `json:"cart"`
		PaymentMethod string          `json:"paymentMethod"`
		ExternalId    string          `json:"externalId,omitempty"`
		CreatedAt     string          `json:"createdAt"`
		UpdatedAt     string          `json:"updatedAt"`
	}

	statusHistory struct {
		CreatedAt string `json:"createdAt"`
		Status    string `json:"status"`
		Message   string `json:"message,omitempty"`
	}

	CartRequest struct {
		Id    string `json:"id"`
		Cant  int    `json:"cant"`
		Name  string `json:"name"`
		Price string `json:"price"`
	}
)

func Create(amount float64, storeId string, paymentMethod string, cart []CartRequest) (Order, error) {
	sh := []statusHistory{{Status: "CREATED", CreatedAt: time.Now().Format(time.RFC3339)}}
	o := Order{
		Id:            uuid.NewString(),
		StoreId:       storeId,
		Amount:        amount,
		Status:        "CREATED",
		StatusHistory: sh,
		PaymentMethod: paymentMethod,
		Cart:          cart,
		CreatedAt:     time.Now().Format(time.RFC3339),
		UpdatedAt:     time.Now().Format(time.RFC3339),
	}

	return o, nil
}

func UpdateStatus(o *Order, status string) {
	o.StatusHistory = append(o.StatusHistory, statusHistory{Status: status, CreatedAt: time.Now().Format(time.RFC3339)})
	o.Status = status
}
