package order

type (
	Client interface {
		GetOrders(storeId string) ([]OrderModel, error)
	}

	client struct {
		d Database
	}
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

func New() Client {
	return &client{
		d: NewDatabase(),
	}
}

func (cli client) GetOrders(storeId string) ([]OrderModel, error) {
	return cli.d.GetOrders(storeId)
}
