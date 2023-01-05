package models

type Order struct {
	Id            string `dynamodbav:"Id"`
	StoreId       string `dynamodbav:"StoreId"`
	Amount        string `dynamodbav:"Amount"`
	Cart          string `dynamodbav:"Cart"`
	PaymentMethod string `dynamodbav:"PaymentMethod"`
	CreatedAt     string `dynamodbav:"CreatedAt"`
	UpdatedAt     string `dynamodbav:"UpdatedAt"`
}
