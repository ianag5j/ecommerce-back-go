package store

import (
	"ianag5j/ecommerce-back-go/get-store/pkg/dto"
	"net/http"
)

type (
	Client interface {
		GetStoreByUserId(userId string) (Store, dto.RequestError)
	}

	client struct {
		d Database
	}

	Store struct {
		Id     string `json:"id,omitempty"`
		Name   string `json:"name,omitempty"`
		UserId string `json:"user_id,omitempty"`
	}
)

func New() Client {
	return client{
		d: NewDatabase(),
	}
}

func (c client) GetStoreByUserId(userId string) (Store, dto.RequestError) {
	s, err := c.d.GetStoreByUserId(userId)
	if err != nil {
		return Store{}, dto.RequestError{
			StatusCode: http.StatusInternalServerError,
			ErrorCode:  dto.INTERNAL_ERROR,
			Err:        err.Error(),
		}
	}

	return s, dto.RequestError{}
}
