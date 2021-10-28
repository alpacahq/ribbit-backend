package plaid

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/alpacahq/ribbit-backend/request"

	"github.com/alpacahq/ribbit-backend/model"
	"github.com/go-pg/pg/v9/orm"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/plaid/plaid-go/plaid"
)

var baseURL string = os.Getenv("BROKER_API_BASE")
var token string = os.Getenv("BROKER_TOKEN")

var (
	PLAID_CLIENT_ID     = os.Getenv("PLAID_CLIENT_ID")
	PLAID_SECRET        = os.Getenv("PLAID_SECRET")
	PLAID_ENV           = os.Getenv("PLAID_ENV")
	PLAID_PRODUCTS      = os.Getenv("PLAID_PRODUCTS")
	PLAID_COUNTRY_CODES = os.Getenv("PLAID_COUNTRY_CODES")
	PLAID_REDIRECT_URI  = os.Getenv("PLAID_REDIRECT_URI")
)

var environments = map[string]plaid.Environment{
	"sandbox":     plaid.Sandbox,
	"development": plaid.Development,
	"production":  plaid.Production,
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file. Did you copy .env.example to .env and fill it out?")
	}
}

var client = func() *plaid.Client {
	clientOptions := plaid.ClientOptions{
		ClientID:    PLAID_CLIENT_ID,
		Secret:      PLAID_SECRET,
		Environment: environments[PLAID_ENV],
		HTTPClient:  &http.Client{},
	}
	client, err := plaid.NewClient(clientOptions)
	if err != nil {
		panic(fmt.Errorf("unexpected error while initializing plaid client %w", err))
	}
	return client
}()

// NewAuthService creates new auth service
func NewPlaidService(userRepo model.UserRepo, accountRepo model.AccountRepo, jwt JWT, db orm.DB, log *zap.Logger) *Service {
	return &Service{userRepo, accountRepo, jwt, db, log}
}

// Service represents the auth application service
type Service struct {
	userRepo    model.UserRepo
	accountRepo model.AccountRepo
	jwt         JWT
	db          orm.DB
	log         *zap.Logger
}

// JWT represents jwt interface
type JWT interface {
	GenerateToken(*model.User) (string, string, error)
}

func (s *Service) CreateLinkToken(c context.Context, accountID string, name string) (*model.PlaidAuthToken, error) {
	countryCodes := strings.Split(PLAID_COUNTRY_CODES, ",")
	products := strings.Split(PLAID_PRODUCTS, ",")

	configs := plaid.LinkTokenConfigs{
		User: &plaid.LinkTokenUser{
			ClientUserID: accountID,
		},
		ClientName:   "Ribbit",
		Products:     products,
		CountryCodes: countryCodes,
		Language:     "en",
	}

	resp, err := client.CreateLinkToken(configs)
	if err != nil {
		return nil, err
	}

	return &model.PlaidAuthToken{
		LinkToken: resp.LinkToken,
	}, nil
}

func (s *Service) SetAccessToken(c context.Context, id int, accountID string, e *request.SetAccessToken) (interface{}, error) {
	response, err := client.ExchangePublicToken(e.PublicToken)
	if err != nil {
		return nil, err
	}

	auth, err := client.GetAuth(response.AccessToken)
	if err != nil {
		return nil, err
	}

	var bank_account_number, bank_routing_number, account_owner_name, bank_account_type, bank_account_name string
	bank_account_type = "CHECKING"
	if len(auth.Numbers.ACH) > 0 {
		for _, account := range auth.Numbers.ACH {
			if e.AccountID == account.AccountID {
				bank_routing_number = account.Routing
				bank_account_number = account.Account
			}
		}
	}

	if bank_routing_number == "" || bank_account_number == "" {
		return nil, errors.New("Bank routing/account number not found")
	}

	identity, err := client.GetIdentity(response.AccessToken)
	if err != nil {
		return nil, err
	}

	if len(identity.Accounts) > 0 {
		for _, ele := range identity.Accounts {
			if ele.AccountID == e.AccountID {
				bank_account_name = ele.Name
				if len(ele.Owners) > 0 {
					for _, owners := range ele.Owners {
						account_owner_name = owners.Names[0]
					}
				}
			}
		}
	}

	client := &http.Client{}
	attachAccount, _ := json.Marshal(map[string]string{
		"account_owner_name":  account_owner_name,
		"bank_account_type":   bank_account_type,
		"bank_account_number": bank_account_number,
		"bank_routing_number": bank_routing_number,
		"nickname":            bank_account_name,
	})
	attachAccountBody := bytes.NewBuffer(attachAccount)

	req, err := http.NewRequest("POST", baseURL+"/v1/accounts/"+accountID+"/ach_relationships", attachAccountBody)
	if err != nil {
		return nil, errors.New("Something went wrong. Try again later.")
	}

	req.Header.Add("Authorization", token)
	cardBody, err := client.Do(req)

	if err != nil {
		return nil, errors.New("Something went wrong. Try again later.")
	}

	responseData, err := ioutil.ReadAll(cardBody.Body)
	if err != nil {
		return nil, errors.New("Something went wrong. Try again later.")
	}

	var responseObject interface{}
	json.Unmarshal(responseData, &responseObject)

	return responseObject, nil
}
