package request

import (
	"net/http"
	"strconv"

	"github.com/alpacahq/ribbit-backend/apperr"

	"github.com/gin-gonic/gin"
)

// Credentials stores the username and password provided in the request
type SetAccessToken struct {
	PublicToken string `json:"public_token" binding:"required"`
	AccountID   string `json:"account_id" binding:"required"`
}

// Login parses out the username and password in gin's request context, into Credentials
func SetAccessTokenbody(c *gin.Context) (*SetAccessToken, error) {
	data := new(SetAccessToken)
	if err := c.ShouldBindJSON(data); err != nil {
		apperr.Response(c, err)
		return nil, err
	}
	return data, nil
}

// Credentials stores the username and password provided in the request
type Charge struct {
	Amount    string `json:"amount" binding:"required"`
	AccountID string `json:"account_id" binding:"required"`
}

// Login parses out the username and password in gin's request context, into Credentials
func ChargeBody(c *gin.Context) (*Charge, error) {
	data := new(Charge)
	if err := c.ShouldBindJSON(data); err != nil {
		apperr.Response(c, err)
		return nil, err
	}
	return data, nil
}

func AccountID(c *gin.Context) (int, error) {
	id, err := strconv.Atoi(c.Param("account_id"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return 0, apperr.New(http.StatusBadRequest, "Account ID isn't valid")
	}
	return id, nil
}
