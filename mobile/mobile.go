package mobile

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/alpacahq/ribbit-backend/config"
)

// NewMobile creates a new mobile service implementation
func NewMobile(config *config.TwilioConfig) *Mobile {
	return &Mobile{config}
}

// Mobile provides a mobile service implementation
type Mobile struct {
	config *config.TwilioConfig
}

// GenerateSMSToken sends an sms token to the mobile numer
func (m *Mobile) GenerateSMSToken(countryCode, mobile string) error {
	apiURL := m.getTwilioVerifyURL()
	data := url.Values{}
	data.Set("To", countryCode+mobile)
	data.Set("Channel", "sms")
	resp, err := m.send(apiURL, data)
	if err != nil {
		return err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)
	return err
}

// CheckCode verifies if the user-provided code is approved
func (m *Mobile) CheckCode(countryCode, mobile, code string) error {
	apiURL := m.getTwilioVerifyURL()
	data := url.Values{}
	data.Set("To", countryCode+mobile)
	data.Set("Code", code)
	resp, err := m.send(apiURL, data)
	if err != nil {
		return err
	}

	// take a look at our response
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Body)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)
	return nil
}

func (m *Mobile) getTwilioVerifyURL() string {
	return "https://verify.twilio.com/v2/Services/" + m.config.Verify + "/Verifications"
}

func (m *Mobile) send(apiURL string, data url.Values) (*http.Response, error) {
	u, _ := url.ParseRequestURI(apiURL)
	urlStr := u.String()
	// http client
	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	r.SetBasicAuth(m.config.Account, m.config.Token)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	return client.Do(r)
}
