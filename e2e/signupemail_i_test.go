package e2e_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/alpacahq/ribbit-backend/request"

	"github.com/stretchr/testify/assert"
)

func (suite *E2ETestSuite) TestSignupEmail() {

	t := suite.T()

	ts := httptest.NewServer(suite.r)
	defer ts.Close()

	urlSignup := ts.URL + "/signup"

	req := &request.EmailSignup{
		Email:    "user@example.org",
		Password: "userpassword1",
	}
	b, err := json.Marshal(req)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post(urlSignup, "application/json", bytes.NewBuffer(b))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))

	assert.Nil(t, err)
}

func (suite *E2ETestSuite) TestVerification() {
	t := suite.T()
	v := suite.v
	// verify that we can retrieve our test verification token
	assert.NotNil(t, v)

	ts := httptest.NewServer(suite.r)
	defer ts.Close()

	url := ts.URL + "/verification/" + v.Token
	fmt.Println("This is our verification url", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
	assert.Nil(t, err)

	// The second time we call our verification url, it should return not found
	resp, err = http.Get(url)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Nil(t, err)
}
