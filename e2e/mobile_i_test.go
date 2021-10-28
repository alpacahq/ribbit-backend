package e2e_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/alpacahq/ribbit-backend/request"

	"github.com/stretchr/testify/assert"
)

func (suite *E2ETestSuite) TestSignupMobile() {
	t := suite.T()
	ts := httptest.NewServer(suite.r)
	defer ts.Close()

	urlSignupMobile := ts.URL + "/mobile"

	req := &request.MobileSignup{
		CountryCode: "+65",
		Mobile:      "91919191",
	}
	b, err := json.Marshal(req)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.Post(urlSignupMobile, "application/json", bytes.NewBuffer(b))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Nil(t, err)

	// the sms code will be separately sms-ed to user's mobile phone, trigger above
	// we now test against the /mobile/verify

	url := ts.URL + "/mobile/verify"
	req2 := &request.MobileVerify{
		CountryCode: "+65",
		Mobile:      "91919191",
		Code:        "123456",
		Signup:      true,
	}
	b, err = json.Marshal(req2)
	if err != nil {
		log.Fatal(err)
	}
	resp, err = http.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Println("Verify Code")
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Nil(t, err)
}
