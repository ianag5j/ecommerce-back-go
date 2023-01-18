package processor

import (
	"encoding/json"
	"ianag5j/ecommerce-back-go/create-store/pkg/dto"
	"ianag5j/ecommerce-back-go/create-store/pkg/store"
	"ianag5j/ecommerce-back-go/create-store/pkg/utils"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type (
	Processor interface {
		Process(r events.APIGatewayProxyRequest) ResponseBody
	}

	processor struct {
		s store.Client
	}

	RequestBody struct {
		StoreName string `json:"name"`
	}

	ResponseBody struct {
		dto.RequestError
		store.Store
	}
)

func New() Processor {
	return &processor{
		s: store.New(),
	}
}

func (p processor) Process(r events.APIGatewayProxyRequest) ResponseBody {
	b := RequestBody{}
	err := json.Unmarshal([]byte(r.Body), &b)

	if err != nil {
		return ResponseBody{
			dto.RequestError{
				Err:        err.Error(),
				ErrorCode:  dto.INTERNAL_ERROR,
				StatusCode: http.StatusInternalServerError,
			},
			store.Store{},
		}
	}

	s, re := p.s.CreateStore(b.StoreName, utils.GetUserId(r.Headers["authorization"]))

	if re.ErrorCode > 0 {
		return ResponseBody{
			re,
			store.Store{},
		}
	}

	return ResponseBody{
		dto.RequestError{
			StatusCode: http.StatusCreated,
		},
		s,
	}
}
