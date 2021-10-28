package service

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/alpacahq/ribbit-backend/apperr"
	"github.com/alpacahq/ribbit-backend/repository/account"
	"github.com/alpacahq/ribbit-backend/repository/transfer"

	"github.com/gin-gonic/gin"
)

type ErrorBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func TransferRouter(svc *transfer.Service, acc *account.Service, r *gin.RouterGroup) {
	a := Transfer{svc, acc}

	ar := r.Group("/transfer")
	ar.GET("", a.transfer)
	ar.GET("/history", a.transfer)
	ar.POST("/bank/:bank_id/deposit", a.createNewTransfer)
	ar.DELETE("/:transfer_id/delete", a.deleteTransfer)
}

// Auth represents auth http service
type Transfer struct {
	svc *transfer.Service
	acc *account.Service
}

func (a *Transfer) transfer(c *gin.Context) {
	id, _ := c.Get("id")
	user := a.acc.GetProfile(c, id.(int))

	if user.AccountID == "" {
		apperr.Response(c, apperr.New(http.StatusBadRequest, "Account not found."))
		return
	}

	limit := c.DefaultQuery("limit", "10000")
	offset := c.DefaultQuery("offset", "0")
	direction := c.DefaultQuery("direction", "")

	transferListURL := os.Getenv("BROKER_API_BASE") + "/v1/accounts/" + user.AccountID + "/transfers?limit=" + limit + "&offset=" + offset + "&direction=" + direction

	client := &http.Client{}
	transferListRequest, _ := http.NewRequest("GET", transferListURL, nil)
	transferListRequest.Header.Add("Authorization", os.Getenv("BROKER_TOKEN"))

	transferList, _ := client.Do(transferListRequest)
	transferListBody, err := ioutil.ReadAll(transferList.Body)
	if err != nil {
		apperr.Response(c, apperr.New(transferList.StatusCode, "Something went wrong. Try again later."))
		return
	}

	if transferList.StatusCode != 200 {
		errorBody := ErrorBody{}
		json.Unmarshal(transferListBody, &errorBody)
		apperr.Response(c, apperr.New(transferList.StatusCode, errorBody.Message))
		return
	}

	var transferListJSON interface{}
	json.Unmarshal(transferListBody, &transferListJSON)
	c.JSON(transferList.StatusCode, transferListJSON)
}

func (a *Transfer) createNewTransfer(c *gin.Context) {
	id, _ := c.Get("id")
	user := a.acc.GetProfile(c, id.(int))
	bankID := c.Param("bank_id")

	amount, err := strconv.ParseFloat(c.PostForm("amount"), 64)
	if err != nil {
		apperr.Response(c, apperr.New(http.StatusBadRequest, "Invalid amount."))
		return
	}

	if amount <= 0 {
		apperr.Response(c, apperr.New(http.StatusBadRequest, "Amount must be greater than 0."))
		return
	}

	if user.AccountID == "" {
		apperr.Response(c, apperr.New(http.StatusBadRequest, "Account not found."))
		return
	}

	createNewTransactionURL := os.Getenv("BROKER_API_BASE") + "/v1/accounts/" + user.AccountID + "/transfers"

	client := &http.Client{}
	createNewTransfer, _ := json.Marshal(map[string]interface{}{
		"transfer_type":   "ach",
		"relationship_id": bankID,
		"amount":          amount,
		"direction":       "INCOMING",
	})
	tranferBody := bytes.NewBuffer(createNewTransfer)

	req, err := http.NewRequest("POST", createNewTransactionURL, tranferBody)
	req.Header.Add("Authorization", os.Getenv("BROKER_TOKEN"))

	createTransfer, err := client.Do(req)
	responseData, err := ioutil.ReadAll(createTransfer.Body)
	if err != nil {
		apperr.Response(c, apperr.New(createTransfer.StatusCode, "Something went wrong. Try again later."))
		return
	}

	if createTransfer.StatusCode != 200 {
		errorBody := ErrorBody{}
		json.Unmarshal(responseData, &errorBody)
		apperr.Response(c, apperr.New(createTransfer.StatusCode, errorBody.Message))
		return
	}

	var responseObject map[string]interface{}
	json.Unmarshal(responseData, &responseObject)

	c.JSON(createTransfer.StatusCode, responseObject)
}

func (a *Transfer) deleteTransfer(c *gin.Context) {
	id, _ := c.Get("id")
	user := a.acc.GetProfile(c, id.(int))

	if user.AccountID == "" {
		apperr.Response(c, apperr.New(http.StatusBadRequest, "Account not found."))
		return
	}
	transferID := c.Param("transfer_id")

	deleteTransfersListURL := os.Getenv("BROKER_API_BASE") + "/v1/accounts/" + user.AccountID + "/transfers/" + transferID

	client := &http.Client{}
	transferDeleteRequest, _ := http.NewRequest("DELETE", deleteTransfersListURL, nil)
	transferDeleteRequest.Header.Add("Authorization", os.Getenv("BROKER_TOKEN"))

	transferDeleteResponse, _ := client.Do(transferDeleteRequest)
	transferDeleteBody, err := ioutil.ReadAll(transferDeleteResponse.Body)
	if err != nil {
		apperr.Response(c, apperr.New(transferDeleteResponse.StatusCode, "Something went wrong. Try again later."))
		return
	}

	if transferDeleteResponse.StatusCode != 200 {
		errorBody := ErrorBody{}
		json.Unmarshal(transferDeleteBody, &errorBody)
		apperr.Response(c, apperr.New(transferDeleteResponse.StatusCode, errorBody.Message))
		return
	}

	var responseObject interface{}
	json.Unmarshal(transferDeleteBody, &responseObject)
	c.JSON(transferDeleteResponse.StatusCode, responseObject)
}
