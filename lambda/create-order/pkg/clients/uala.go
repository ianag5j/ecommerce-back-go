package uala

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	credential "ianag5j/ecommerce-back-go/create-order/pkg/credential/models"
	"ianag5j/ecommerce-back-go/create-order/pkg/order"
	store "ianag5j/ecommerce-back-go/create-order/pkg/store/models"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type (
	Uala interface {
		CreateOrder(order order.Order, store store.Store) (UalaOrder, error)
	}

	uala struct {
		clientConfig clientConfig
	}

	clientConfig struct {
		accessToken string
		userName    string
		url         string
	}
)

type (
	authRequest struct {
		UserName       string `json:"user_name"`
		ClientId       string `json:"client_id"`
		ClientSecretId string `json:"client_secret_id"`
		GrantType      string `json:"grant_type"`
	}

	authResponse struct {
		AccessToken string `json:"access_token"`
	}
)

type (
	checkoutRequest struct {
		Amount          string `json:"amount"`
		Description     string `json:"description"`
		UserName        string `json:"userName"`
		CallbackFail    string `json:"callback_fail"`
		CallbackSuccess string `json:"callback_success"`
		NotificationUrl string `json:"notification_url"`
	}

	UalaOrder struct {
		Uuid  string `json:"uuid"`
		Links struct {
			CheckoutLink string `json:"checkoutLink"`
			Success      string `json:"success"`
			Failed       string `json:"failed"`
		} `json:"links"`
	}
)

func New(c credential.Credential) Uala {
	ac, _ := createToken(c)
	u := "https://checkout.stage.ua.la/1"
	if os.Getenv("ENVIROMENT") == "prod" {
		u = "https://checkout.prod.ua.la/1"
	}
	return &uala{
		clientConfig: clientConfig{
			accessToken: ac,
			userName:    c.ExternalUserName,
			url:         u,
		},
	}
}

func createToken(c credential.Credential) (string, error) {
	au := "https://auth.stage.ua.la/1/auth/token"
	if os.Getenv("ENVIROMENT") == "prod" {
		au = "https://auth.prod.ua.la/1/auth/token"
	}

	j, _ := json.Marshal(authRequest{
		UserName:       c.ExternalUserName,
		ClientId:       c.ExternalClientId,
		ClientSecretId: c.ExternalClientSecret,
		GrantType:      "client_credentials",
	})
	fmt.Println(string(j))
	resp, _ := http.Post(au, "application/json", bytes.NewBuffer(j))
	fmt.Println(resp.Status)

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("error on create token")
	}

	at := authResponse{}
	json.NewDecoder(resp.Body).Decode(&at)
	return at.AccessToken, nil
}

func (u uala) CreateOrder(o order.Order, s store.Store) (UalaOrder, error) {
	uo := UalaOrder{}
	j, _ := json.Marshal(checkoutRequest{
		Amount:          strconv.FormatFloat(o.Amount, 'g', 3, 64),
		Description:     s.Name + " order",
		UserName:        u.clientConfig.userName,
		CallbackFail:    os.Getenv("FRONT_BASE_URL") + "/fail",
		CallbackSuccess: os.Getenv("FRONT_BASE_URL") + "/success",
		NotificationUrl: os.Getenv("LAMBDA_URL") + "/v2/uala-webhook/" + o.Id,
	})

	req, err := http.NewRequest("POST", u.clientConfig.url+"/checkout", bytes.NewBuffer(j))
	if err != nil {
		return uo, errors.New("error on create new request")
	}

	req.Header.Add("Authorization", "Bearer "+u.clientConfig.accessToken)
	req.Header.Add("Content-Type", "application/json")
	fmt.Println(req)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return uo, errors.New(fmt.Sprintf("error on response. err: %e", err.Error()))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(b), string(j))
		return uo, errors.New(fmt.Sprintf("error on create uala order, status code: %s", resp.Status))
	}
	json.NewDecoder(resp.Body).Decode(&uo)
	fmt.Println(uo)
	return uo, nil
}
