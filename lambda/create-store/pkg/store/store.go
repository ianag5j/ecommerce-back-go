package store

import (
	"errors"
	"ianag5j/ecommerce-back-go/create-store/pkg/dto"
	"net/http"

	"github.com/google/uuid"
)

type (
	Client interface {
		CreateStore(storeName string, userId string) (Store, dto.RequestError)
		validateBeforeCreateStore(storeName string, userId string) dto.RequestError
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

func (c client) CreateStore(storeName string, userId string) (Store, dto.RequestError) {
	re := c.validateBeforeCreateStore(storeName, userId)
	if re.Err != "" {
		return Store{}, re
	}

	s, err := c.d.CreateStore(Store{
		Id:     uuid.NewString(),
		Name:   storeName,
		UserId: userId,
	})
	if err != nil {
		return Store{}, dto.RequestError{
			StatusCode: http.StatusInternalServerError,
			ErrorCode:  dto.INTERNAL_ERROR,
			Err:        err.Error(),
		}
	}

	return s, dto.RequestError{}
}

func (c client) validateBeforeCreateStore(storeName string, userId string) dto.RequestError {
	// TODO: use async
	s, err := c.d.GetStoreByName(storeName)
	if err != nil {
		return dto.RequestError{
			StatusCode: http.StatusInternalServerError,
			ErrorCode:  dto.INTERNAL_ERROR,
			Err:        err.Error(),
		}
	}
	if s != (Store{}) {
		return dto.RequestError{
			StatusCode: http.StatusBadRequest,
			ErrorCode:  dto.STORE_NAME_USED,
			Err:        errors.New("store name in use").Error(),
		}
	}

	s, err = c.d.GetStoreByUserId(userId)
	if err != nil {
		return dto.RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        err.Error(),
		}
	}
	if s != (Store{}) {
		return dto.RequestError{
			StatusCode: http.StatusBadRequest,
			ErrorCode:  dto.USER_ALREADY_HAS_STORE,
			Err:        errors.New("user already has store created").Error(),
		}
	}

	return dto.RequestError{}
}
