package service

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/alpacahq/ribbit-backend/apperr"
	"github.com/alpacahq/ribbit-backend/repository/account"
	"github.com/alpacahq/ribbit-backend/repository/plaid"
	"github.com/alpacahq/ribbit-backend/request"

	"github.com/gin-gonic/gin"
)

func PlaidRouter(svc *plaid.Service, acc *account.Service, r *gin.RouterGroup) {
	a := Plaid{svc, acc}

	ar := r.Group("/plaid")
	ar.GET("/create_link_token", a.createLinkToken)
	ar.POST("/set_access_token", a.setAccessToken)
	ar.GET("/recipient_banks", a.accountsList)
	ar.DELETE("/recipient_banks/:bank_id", a.detachAccount)
}

// Auth represents auth http service
type Plaid struct {
	svc *plaid.Service
	acc *account.Service
}

func (a *Plaid) createLinkToken(c *gin.Context) {
	id, _ := c.Get("id")
	user := a.acc.GetProfile(c, id.(int))

	name := user.FirstName + " " + user.LastName
	linkToken, err := a.svc.CreateLinkToken(c, user.AccountID, name)
	if err != nil {
		apperr.Response(c, err)
		return
	}

	c.JSON(http.StatusOK, linkToken)
}

func (a *Plaid) setAccessToken(c *gin.Context) {
	data, err := request.SetAccessTokenbody(c)
	if err != nil {
		apperr.Response(c, err)
		return
	}

	id, _ := c.Get("id")
	user := a.acc.GetProfile(c, id.(int))

	if user.AccountID == "" {
		apperr.Response(c, apperr.New(http.StatusBadRequest, "Account not found."))
		return
	}

	response, err := a.svc.SetAccessToken(c, id.(int), user.AccountID, data)
	if err != nil {
		apperr.Response(c, apperr.New(http.StatusBadRequest, err.Error()))
		return
	}
	c.JSON(http.StatusOK, response)
}

func (a *Plaid) accountsList(c *gin.Context) {
	id, _ := c.Get("id")
	user := a.acc.GetProfile(c, id.(int))
	accountID := user.AccountID

	if accountID == "" {
		apperr.Response(c, apperr.New(http.StatusBadRequest, "Account not found."))
		return
	}

	client := &http.Client{}
	acountStatus, _ := json.Marshal(map[string]string{
		"statuses": "QUEUED,APPROVED,PENDING",
	})
	accountStatuses := bytes.NewBuffer(acountStatus)

	getAchAccountsList := os.Getenv("BROKER_API_BASE") + "/v1/accounts/" + accountID + "/ach_relationships"

	req, _ := http.NewRequest("GET", getAchAccountsList, accountStatuses)
	req.Header.Add("Authorization", os.Getenv("BROKER_TOKEN"))

	response, _ := client.Do(req)
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		apperr.Response(c, apperr.New(http.StatusInternalServerError, "Something went wrong. Try again later."))
		return
	}

	var responseObject interface{}
	json.Unmarshal(responseData, &responseObject)
	c.JSON(response.StatusCode, responseObject)
}

func (a *Plaid) detachAccount(c *gin.Context) {
	id, _ := c.Get("id")
	user := a.acc.GetProfile(c, id.(int))

	accountID := user.AccountID
	bankID := c.Param("bank_id")

	if accountID == "" {
		apperr.Response(c, apperr.New(http.StatusBadRequest, "Account not found."))
		return
	}

	deleteAccountAPIURL := os.Getenv("BROKER_API_BASE") + "/v1/accounts/" + accountID + "/ach_relationships/" + bankID

	client := &http.Client{}
	req, _ := http.NewRequest("DELETE", deleteAccountAPIURL, nil)
	req.Header.Add("Authorization", os.Getenv("BROKER_TOKEN"))
	response, _ := client.Do(req)

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		apperr.Response(c, apperr.New(http.StatusInternalServerError, "Something went wrong. Try again later."))
		return
	}

	var responseObject interface{}
	json.Unmarshal(responseData, &responseObject)
	c.JSON(response.StatusCode, responseObject)
}
