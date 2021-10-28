package e2e_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/alpacahq/ribbit-backend/model"
	"github.com/alpacahq/ribbit-backend/request"

	"github.com/stretchr/testify/assert"
)

func (suite *E2ETestSuite) TestLogin() {
	t := suite.T()

	ts := httptest.NewServer(suite.r)
	defer ts.Close()

	url := ts.URL + "/login"

	req := &request.Credentials{
		Email:    "superuser@example.org",
		Password: "testpassword",
	}
	b, err := json.Marshal(req)
	assert.Nil(t, err)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var authToken model.AuthToken
	err = json.Unmarshal(body, &authToken)
	assert.Nil(t, err)
	assert.NotNil(t, authToken)
	suite.authToken = authToken
}

func (suite *E2ETestSuite) TestRefreshToken() {
	t := suite.T()
	assert.NotNil(t, suite.authToken)

	ts := httptest.NewServer(suite.r)
	defer ts.Close()

	// delay by 1 second so that our re-generated JWT will have a 1 second difference
	time.Sleep(1 * time.Second)

	url := ts.URL + "/refresh/" + suite.authToken.RefreshToken
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var refreshToken model.RefreshToken
	err = json.Unmarshal(body, &refreshToken)
	assert.Nil(t, err)
	assert.NotNil(t, refreshToken)

	// because of a 1 second delay, our re-generated JWT will definitely be different
	assert.NotEqual(t, suite.authToken.Token, refreshToken.Token)
	assert.NotEqual(t, suite.authToken.Expires, refreshToken.Expires)
}
