package order

import (
	"errors"
	"strconv"
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

func Create(cart []CartRequest, storeId string, paymentMethod string) (Order, error) {
	o := Order{}
	sh := []statusHistory{{Status: "CREATED", CreatedAt: time.Now().Format(time.RFC3339)}}

	ta, err := getTotalAmount(cart)
	if err != nil {
		return o, err
	}

	o = Order{
		Id:            uuid.NewString(),
		StoreId:       storeId,
		Amount:        ta,
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

func getTotalAmount(c []CartRequest) (float64, error) {
	ta := 0.0

	//TODO: validate amount with products amounts in database
	for _, p := range c {
		pa, err := strconv.ParseFloat(p.Price, 64)
		if err != nil {
			return ta, errors.New("error on parse amount")
		}
		pc := float64(p.Cant)
		ta += pc * pa
	}
	return ta, nil
}
