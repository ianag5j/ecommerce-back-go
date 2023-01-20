package processor

import (
	"ianag5j/ecommerce-back-go/get-store/pkg/dto"
	"ianag5j/ecommerce-back-go/get-store/pkg/store"
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
	ui := r.RequestContext.Authorizer["userId"]
	s, re := p.s.GetStoreByUserId(ui.(string))
	if re.ErrorCode > 0 {
		return ResponseBody{
			re,
			store.Store{},
		}
	}

	if s.Id == "" {
		return ResponseBody{
			dto.RequestError{
				StatusCode: http.StatusNoContent,
			},
			s,
		}
	}

	return ResponseBody{
		dto.RequestError{
			StatusCode: http.StatusOK,
		},
		s,
	}
}
