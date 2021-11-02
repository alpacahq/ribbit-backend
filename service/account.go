package service

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/alpacahq/ribbit-backend/apperr"
	"github.com/alpacahq/ribbit-backend/model"
	"github.com/alpacahq/ribbit-backend/repository/account"
	"github.com/alpacahq/ribbit-backend/request"

	"github.com/bradfitz/slice"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9/orm"
	shortuuid "github.com/lithammer/shortuuid/v3"
)

// AccountService represents the account http service
type AccountService struct {
	svc *account.Service
	db  orm.DB
}

// AccountRouter sets up all the controller functions to our router
func AccountRouter(svc *account.Service, db orm.DB, r *gin.RouterGroup) {
	a := AccountService{
		svc: svc,
		db:  db,
	}
	pr := r.Group("/profile")
	pr.GET("", a.profile)
	pr.POST("/avatar", a.uploadAvatar)
	pr.DELETE("/avatar", a.deleteAvatar)
	pr.GET("/shareable-link", a.getShareableProfileLik)
	pr.PATCH("", a.updateProfile)

	cr := r.Group("/countries")
	cr.GET("", a.countriesList)

	cr1 := cr.Group("/:country_code")
	cr1.GET("/states", a.statesList)
	cr1.GET("/states/:state_code/cities", a.citiesList)

	clr := r.Group("/clock")
	clr.GET("", a.clock)

	acr := r.Group("/account")
	acr.GET("", a.getAccount)
	acr.POST("/sign", a.sign)
	acr.GET("/portfolio/history", a.portfolioHistory)
	acr.GET("/trading-profile", a.tradingProfile)
	acr.GET("/stats", a.stats)

	rr := r.Group("/referral")
	rr.GET("", a.getShareableProfileLik)

	ar := r.Group("/users")
	ar.POST("", a.create)
	ar.PATCH("/:id/password", a.changePassword)

	ac := r.Group("/orders")
	ac.GET("", a.getOrders)
	ac.POST("", a.createOrder)
	ac.GET("/:order_id", a.getOrderDetails)
	ac.PATCH("/:order_id", a.replaceOrder)
	ac.DELETE("", a.cancelAllOrders)
	ac.DELETE("/:order_id", a.cancelOrder)

	pz := r.Group("/positions")
	pz.GET("", a.getPositions)
	pz.GET("/:symbol", a.getOneOpenPosition)
	pz.DELETE("", a.closePositions)
	pz.DELETE("/:symbol", a.closeOnePosition)

	mrk := r.Group("/market")
	mrk.GET("/tickers", a.getMarketTickers)
	mrk.GET("/tickers/:symbol", a.getMarketTickerBySymbol)
	mrk.GET("/stocks/:symbol/trades", a.getMarketTradesBySymbol)
	mrk.GET("/stocks/:symbol/trades/latest", a.getMarketLatestTradeBySymbol)
	mrk.GET("/stocks/:symbol/quotes", a.getMarketQuotesBySymbol)
	mrk.GET("/stocks/:symbol/quotes/latest", a.getMarketLatestQuoteBySymbol)
	mrk.GET("/stocks/:symbol/bars", a.getMarketBarsBySymbol)

	watchlist := r.Group("/watchlist")
	watchlist.GET("", a.getWatchList)
	watchlist.POST("", a.addAssetInWatchList)
	watchlist.DELETE("/:symbol", a.removeAssetFromWatchList)

	cl := r.Group("/calendar")
	cl.GET("", a.getCalendar)

}

func (a *AccountService) create(c *gin.Context) {
	r, err := request.AccountCreate(c)
	if err != nil {
		return
	}
	user := &model.User{
		Username:            r.Username,
		Password:            r.Password,
		Email:               r.Email,
		FirstName:           r.FirstName,
		LastName:            r.LastName,
		RoleID:              r.RoleID,
		AccountID:           r.AccountID,
		AccountNumber:       r.AccountNumber,
		AccountCurrency:     r.AccountCurrency,
		AccountStatus:       r.AccountStatus,
		DOB:                 r.DOB,
		City:                r.City,
		State:               r.State,
		Country:             r.Country,
		TaxIDType:           r.TaxIDType,
		TaxID:               r.TaxID,
		FundingSource:       r.FundingSource,
		EmploymentStatus:    r.EmploymentStatus,
		InvestingExperience: r.InvestingExperience,
		PublicShareholder:   r.PublicShareholder,
		AnotherBrokerage:    r.AnotherBrokerage,
		DeviceID:            r.DeviceID,
		ProfileCompletion:   r.ProfileCompletion,
		ReferralCode:        r.ReferralCode,
	}
	if err := a.svc.Create(c, user); err != nil {
		apperr.Response(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func (a *AccountService) changePassword(c *gin.Context) {
	p, err := request.PasswordChange(c)
	if err != nil {
		return
	}
	if err := a.svc.ChangePassword(c, p.OldPassword, p.NewPassword, p.ID); err != nil {
		apperr.Response(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (a *AccountService) profile(c *gin.Context) {
	id, _ := c.Get("id")
	user := a.svc.GetProfile(c, id.(int))
	if user != nil {
		c.JSON(http.StatusOK, user)
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Couldn't fetch profile.",
	})
}

func (a *AccountService) uploadAvatar(c *gin.Context) {
	id, _ := c.Get("id")
	user := a.svc.GetProfile(c, id.(int))
	const MAX_UPLOAD_SIZE = 1024 * 1024 * 5 // 5MB
	if user != nil {
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			apperr.Response(c, apperr.New(apperr.BadRequest.Status, fmt.Sprintf("file err : %s", err.Error())))
			return
		}
		// Check size
		if err := c.Request.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
			apperr.Response(c, apperr.New(apperr.BadRequest.Status, "The uploaded file is too big. Please choose an file that's less than 5MB in size"))
			return
		}
		// Check file type
		// if strings.Split(header.Header.Get("Content-Type"), "/")[0] != "image" {
		// 	apperr.Response(c, apperr.New(apperr.BadRequest.Status, "The uploaded file is not an image"))
		// 	return
		// }
		buff := make([]byte, 512)
		if _, err = file.Read(buff); err != nil {
			apperr.Response(c, apperr.New(apperr.BadRequest.Status, fmt.Sprintf("%s", err.Error())))
			return
		}
		if filetype := http.DetectContentType(buff); filetype != "image/jpg" && filetype != "image/jpeg" && filetype != "image/png" {
			apperr.Response(c, apperr.New(apperr.BadRequest.Status, "The provided file format is not allowed. Please upload a JPEG or PNG image"))
			return
		}
		if _, err := file.Seek(0, io.SeekStart); err != nil {
			apperr.Response(c, apperr.New(apperr.BadRequest.Status, fmt.Sprintf("%s", err.Error())))
			return
		}

		filename := filepath.Base(header.Filename)
		if filename == "." {
			apperr.Response(c, apperr.New(http.StatusUnprocessableEntity, "File has malformed file format"))
		}
		newFileName := strconv.Itoa(user.ID) + "-" + filename
		// delete old file
		if user.Avatar != "" {
			if err := os.Remove("public/users/" + user.Avatar); err != nil {
				log.Fatal(err)
			}
		}

		out, err := os.Create("public/users/" + newFileName)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			log.Fatal(err)
		}
		if err := a.svc.UpdateAvatar(c, newFileName, user.ID); err != nil {
			apperr.Response(c, err)
			return
		}
		c.JSON(http.StatusOK, map[string]interface{}{
			"avatar": newFileName,
		})
		return
	}
	apperr.Response(c, apperr.New(apperr.BadRequest.Status, "Internal Server Error, please try again."))
}

func (a *AccountService) deleteAvatar(c *gin.Context) {
	id, _ := c.Get("id")
	user := a.svc.GetProfile(c, id.(int))
	if user != nil {
		if user.Avatar != "" {
			if err := os.Remove("public/users/" + user.Avatar); err != nil {
				log.Fatal(err)
			}
		}
		if err := a.svc.UpdateAvatar(c, "", user.ID); err != nil {
			apperr.Response(c, err)
			return
		}
		apperr.Response(c, apperr.New(http.StatusOK, "Congratulations! Your avatar has been deleted successfully."))
		return
	}
	apperr.Response(c, apperr.New(http.StatusBadRequest, "Internal Server Error, please try again."))
}

type ShareableProfileLink struct {
	URL  string `json:"url"`
	Code string `json:"code"`
}

func (a *AccountService) getShareableProfileLik(c *gin.Context) {
	id, _ := c.Get("id")
	user := a.svc.GetProfile(c, id.(int))
	if user != nil {
		c.JSON(http.StatusOK, ShareableProfileLink{
			URL:  "https://alpaca.com/profile/" + user.ReferralCode,
			Code: user.ReferralCode,
		})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Couldn't fetch shareable link.",
	})
}

func (a *AccountService) updateProfile(c *gin.Context) {
	p, err := request.UpdateProfile(c)
	if err != nil {
		return
	}
	user, err := a.svc.UpdateProfile(c, p)
	if err != nil {
		apperr.Response(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

type Country struct {
	Name      string `json:"name"`
	ShortCode string `json:"short_code"`
}

func countryAlreadyExists(a []Country, x Country) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func (a *AccountService) countriesList(c *gin.Context) {
	path, _ := os.Getwd()
	q := c.Query("q")

	csvFile, err := os.Open(path + "/countries.csv")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened CSV file")
	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		fmt.Println(err)
	}
	countries := []Country{}
	for _, line := range csvLines {
		country := Country{
			ShortCode: line[0],
			Name:      line[1],
		}
		if !countryAlreadyExists(countries, country) {
			if len(q) > 0 {
				if strings.Contains(strings.ToLower(country.Name), strings.ToLower(q)) {
					countries = append(countries, country)
				}
			} else {
				countries = append(countries, country)
			}
		}
	}

	c.JSON(http.StatusOK, countries[1:])
	return

	// c.JSON(http.StatusOK, []Country{{
	// 	ShortCode: "USA",
	// 	Name:      "United States",
	// }})
}

type State struct {
	Name      string `json:"name"`
	ShortCode string `json:"short_code"`
}

func stateAlreadyExists(a []State, x State) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func (a *AccountService) statesList(c *gin.Context) {
	path, _ := os.Getwd()
	countryCode := c.Param("country_code")
	q := c.Query("q")

	if countryCode == "USA" {
		csvFile, err := os.Open(path + "/uscities.csv")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Successfully Opened CSV file")
		defer csvFile.Close()

		csvLines, err := csv.NewReader(csvFile).ReadAll()
		if err != nil {
			fmt.Println(err)
		}
		states := []State{}
		for _, line := range csvLines {
			state := State{
				ShortCode: line[2],
				Name:      line[3],
			}
			if !stateAlreadyExists(states, state) {
				if len(q) > 0 {
					if strings.Contains(strings.ToLower(state.Name), strings.ToLower(q)) {
						states = append(states, state)
					}
				} else {
					states = append(states, state)
				}
			}
		}
		slice.Sort(states[:], func(i, j int) bool {
			return states[i].Name < states[j].Name
		})

		c.JSON(http.StatusOK, states[1:])
		return
	}
	c.JSON(http.StatusNotFound, gin.H{
		"message": "400 Not found",
	})
}

type City struct {
	Name  string `json:"name"`
	ASCII string `json:"ascii"`
	LAT   string `json:"lat"`
	LNG   string `json:"lng"`
}

func (a *AccountService) citiesList(c *gin.Context) {
	path, _ := os.Getwd()
	countryCode := c.Param("country_code")
	stateCode := c.Param("state_code")
	q := c.Query("q")

	if countryCode == "USA" {
		csvFile, err := os.Open(path + "/uscities.csv")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Successfully Opened CSV file")
		defer csvFile.Close()

		csvLines, err := csv.NewReader(csvFile).ReadAll()
		if err != nil {
			fmt.Println(err)
		}
		citys := []City{}
		for _, line := range csvLines {
			city := City{
				Name:  line[0],
				ASCII: line[1],
				LAT:   line[6],
				LNG:   line[7],
			}
			if line[2] == stateCode {
				if len(q) > 0 {
					if strings.Contains(strings.ToLower(city.Name), strings.ToLower(q)) {
						citys = append(citys, city)
					}
				} else {
					citys = append(citys, city)
				}

			}
		}

		c.JSON(http.StatusOK, citys[1:])
		return
	}
	c.JSON(http.StatusNotFound, gin.H{
		"message": "400 Not found",
	})
}

type Referral struct {
	Code            string `json:"code"`
	Url             string `json:"url"`
	ReferredSignups int    `json:"referred_signups"`
}

func (a *AccountService) getReferralUrl(c *gin.Context) {
	u := shortuuid.New()

	c.JSON(http.StatusOK, Referral{
		Code:            u,
		Url:             "https://alpa.ca/" + u,
		ReferredSignups: 0,
	})
}

type BrokerContact struct {
	Email   string   `json:"email_address"`
	Phone   string   `json:"phone_number"`
	Address []string `json:"street_address"`
	City    string   `json:"city"`
	State   string   `json:"state"`
	Country string   `json:"country"`
}
type BrokerIdentity struct {
	FirstName             string   `json:"given_name"`
	LastName              string   `json:"family_name"`
	DateOfBirth           string   `json:"date_of_birth"`
	TaxID                 string   `json:"tax_id"`
	TaxIDType             string   `json:"tax_id_type"`
	CountryOfCitizenship  string   `json:"country_of_citizenship"`
	CountryOfBirth        string   `json:"country_of_birth"`
	CountryOfTaxResidence string   `json:"country_of_tax_residence"`
	FundingSource         []string `json:"funding_source"`
}
type BrokerDisclosures struct {
	IsControlPerson             bool `json:"is_control_person"`
	IsAffiliatedExchangeOrFinra bool `json:"is_affiliated_exchange_or_finra"`
	IsPoliticallyExposed        bool `json:"is_politically_exposed"`
	ImmediateFamilyExposed      bool `json:"immediate_family_exposed"`
}
type BrokerAgreement struct {
	Agreement string `json:"agreement"`
	SignedAt  string `json:"signed_at"`
	IPAddress string `json:"ip_address"`
}

// type BrokerDocument struct {
// 	DocumentType    string `json:"document_type"`
// 	DocumentSubType string `json:"document_sub_type"`
// 	Content         string `json:"content"`
// 	MimeType        string `json:"mime_type"`
// }
// type BrokerTrustedContact struct {
// 	Code            string `json:"code"`
// 	Url             string `json:"url"`
// 	ReferredSignups int    `json:"referred_signups"`
// }
type BrokerAccount struct {
	Contact     BrokerContact     `json:"contact"`
	Identity    BrokerIdentity    `json:"identity"`
	Disclosures BrokerDisclosures `json:"disclosures"`
	Agreements  []BrokerAgreement `json:"agreements"`
}
type BrokerAccountResponse struct {
	AccountNumber string `json:"account_number"`
	CreatedAt     string `json:"created_at"`
	Currency      string `json:"currency"`
	ID            string `json:"id"`
	LastEquity    string `json:"last_equity"`
	Status        string `json:"status"`
}

func getBrokerAccount(u *model.User) BrokerAccount {
	account := BrokerAccount{
		Contact: BrokerContact{
			Email:   u.Email,
			Phone:   u.Mobile,
			Address: []string{u.Address},
			City:    u.City,
			State:   u.State,
			Country: "USA",
		},
		Identity: BrokerIdentity{
			FirstName:             u.FirstName,
			LastName:              u.LastName,
			DateOfBirth:           u.DOB,
			TaxID:                 u.TaxID,
			TaxIDType:             u.TaxIDType,
			CountryOfCitizenship:  "USA",
			CountryOfBirth:        "USA",
			CountryOfTaxResidence: "USA",
			FundingSource:         strings.Split(u.FundingSource, ","),
		},
		Disclosures: BrokerDisclosures{
			IsControlPerson:             false,
			IsAffiliatedExchangeOrFinra: false,
			IsPoliticallyExposed:        false,
			ImmediateFamilyExposed:      false,
		},
		Agreements: []BrokerAgreement{
			{
				Agreement: "margin_agreement",
				SignedAt:  time.Now().Format(time.RFC3339),
				IPAddress: "127.0.0.1",
			},
			{
				Agreement: "account_agreement",
				SignedAt:  time.Now().Format(time.RFC3339),
				IPAddress: "127.0.0.1",
			},
			{
				Agreement: "customer_agreement",
				SignedAt:  time.Now().Format(time.RFC3339),
				IPAddress: "127.0.0.1",
			},
		},
	}

	return account
}

func (a *AccountService) sign(c *gin.Context) {

	id, _ := c.Get("id")
	user := a.svc.GetProfile(c, id.(int))
	if user != nil {
		brokerAccount := getBrokerAccount(user)
		requestBytes, _ := json.Marshal(brokerAccount)
		fmt.Println(string(requestBytes))

		client := &http.Client{}
		req, err := http.NewRequest("POST", os.Getenv("BROKER_API_BASE")+"/v1/accounts", bytes.NewReader(requestBytes))
		if err != nil {
			fmt.Print(err.Error())
		}

		req.Header.Add("Authorization", os.Getenv("BROKER_TOKEN"))
		response, err := client.Do(req)

		if err != nil {
			fmt.Print(err.Error())
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(responseData))
		if strings.Contains(string(responseData), "account_number") {
			brokerResponse := BrokerAccountResponse{}
			json.Unmarshal(responseData, &brokerResponse)

			reqUser := request.Update{
				ID:              user.ID,
				AccountID:       &brokerResponse.ID,
				AccountCurrency: &brokerResponse.Currency,
				AccountNumber:   &brokerResponse.AccountNumber,
				AccountStatus:   &brokerResponse.Status,
			}

			user2, err := a.svc.UpdateProfile(c, &reqUser)
			if err != nil {
				apperr.Response(c, err)
				return
			}
			c.JSON(http.StatusOK, user2)
			return

		} else {
			var brokerResponse interface{}
			json.Unmarshal(responseData, &brokerResponse)
			c.JSON(response.StatusCode, brokerResponse)
			return
		}
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Couldn't sign the account.",
	})
}

func (a *AccountService) clock(c *gin.Context) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", os.Getenv("BROKER_API_BASE")+"/v1/clock", nil)
	if err != nil {
		fmt.Print(err.Error())
	}

	req.Header.Add("Authorization", os.Getenv("BROKER_TOKEN"))
	response, err := client.Do(req)

	if err != nil {
		fmt.Print(err.Error())
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject interface{}
	json.Unmarshal(responseData, &responseObject)

	c.JSON(response.StatusCode, responseObject)
}

func (a *AccountService) getOrders(c *gin.Context) {
	id, _ := c.Get("id")
	user := a.svc.GetProfile(c, id.(int))
	if user != nil && user.AccountID != "" {
		client := &http.Client{}
		accountID := user.AccountID

		req, err := http.NewRequest("GET", os.Getenv("BROKER_API_BASE")+"/v1/trading/accounts/"+accountID+"/orders", nil)
		if err != nil {
			fmt.Print(err.Error())
		}

		req.Header.Add("Authorization", os.Getenv("BROKER_TOKEN"))
		response, err := client.Do(req)

		if err != nil {
			fmt.Print(err.Error())
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		var responseObject interface{}
		json.Unmarshal(responseData, &responseObject)

		c.JSON(response.StatusCode, responseObject)
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Couldn't fetch orders",
	})
}

func (a *AccountService) createOrder(c *gin.Context) {
	id, _ := c.Get("id")
	user := a.svc.GetProfile(c, id.(int))
	if user != nil && user.AccountID != "" {
		client := &http.Client{}
		accountID := user.AccountID

		req, err := http.NewRequest("POST", os.Getenv("BROKER_API_BASE")+"/v1/trading/accounts/"+accountID+"/orders", c.Request.Body)
		if err != nil {
			fmt.Print(err.Error())
		}

		req.Header.Add("Authorization", os.Getenv("BROKER_TOKEN"))
		response, err := client.Do(req)

		if err != nil {
			fmt.Print(err.Error())
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		var responseObject interface{}
		json.Unmarshal(responseData, &responseObject)

		c.JSON(response.StatusCode, responseObject)
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Couldn't create order.",
	})
}

func (a *AccountService) getOrderDetails(c *gin.Context) {
	id, _ := c.Get("id")
	user := a.svc.GetProfile(c, id.(int))
	if user != nil && user.AccountID != "" {
		client := &http.Client{}
		accountID := user.AccountID

		orderID := c.Param("order_id")

		req, err := http.NewRequest("GET", os.Getenv("BROKER_API_BASE")+"/v1/trading/accounts/"+accountID+"/orders/"+orderID, nil)
		if err != nil {
			fmt.Print(err.Error())
		}

		req.Header.Add("Authorization", os.Getenv("BROKER_TOKEN"))
		response, err := client.Do(req)

		if err != nil {
			fmt.Print(err.Error())
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		var responseObject interface{}
		json.Unmarshal(responseData, &responseObject)

		c.JSON(response.StatusCode, responseObject)
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Couldn't get order details.",
	})
}

func (a *AccountService) replaceOrder(c *gin.Context) {
	id, _ := c.Get("id")
	user := a.svc.GetProfile(c, id.(int))
	if user != nil && user.AccountID != "" {
		client := &http.Client{}
		accountID := user.AccountID
		orderID := c.Param("order_id")

		req, err := http.NewRequest("PATCH", os.Getenv("BROKER_API_BASE")+"/v1/trading/accounts/"+accountID+"/orders/"+orderID, c.Request.Body)
		if err != nil {
			fmt.Print(err.Error())
		}

		req.Header.Add("Authorization", os.Getenv("BROKER_TOKEN"))
		response, err := client.Do(req)

		if err != nil {
			fmt.Print(err.Error())
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		var responseObject interface{}
		json.Unmarshal(responseData, &responseObject)

		c.JSON(response.StatusCode, responseObject)
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Couldn't replace the order.",
	})
}

func (a *AccountService) cancelAllOrders(c *gin.Context) {
	id, _ := c.Get("id")
	user := a.svc.GetProfile(c, id.(int))
	if user != nil && user.AccountID != "" {
		client := &http.Client{}
		accountID := user.AccountID

		req, err := http.NewRequest("DELETE", os.Getenv("BROKER_API_BASE")+"/v1/trading/accounts/"+accountID+"/orders", nil)
		if err != nil {
			fmt.Print(err.Error())
		}

		req.Header.Add("Authorization", os.Getenv("BROKER_TOKEN"))
		response, err := client.Do(req)

		if err != nil {
			fmt.Print(err.Error())
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		var responseObject interface{}
		json.Unmarshal(responseData, &responseObject)

		c.JSON(response.StatusCode, responseObject)
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Couldn't cancel orders.",
	})
}

func (a *AccountService) cancelOrder(c *gin.Context) {
	id, _ := c.Get("id")
	user := a.svc.GetProfile(c, id.(int))
	if user != nil && user.AccountID != "" {
		client := &http.Client{}
		accountID := user.AccountID
		orderID := c.Param("order_id")

		req, err := http.NewRequest("DELETE", os.Getenv("BROKER_API_BASE")+"/v1/trading/accounts/"+accountID+"/orders/"+orderID, nil)
		if err != nil {
			fmt.Print(err.Error())
		}

		req.Header.Add("Authorization", os.Getenv("BROKER_TOKEN"))
		response, err := client.Do(req)

		if err != nil {
			fmt.Print(err.Error())
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		var responseObject interface{}
		json.Unmarshal(responseData, &responseObject)

		c.JSON(response.StatusCode, responseObject)
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Couldn't cancel order.",
	})
}

func (a *AccountService) portfolioHistory(c *gin.Context) {
	id, _ := c.Get("id")
	user := a.svc.GetProfile(c, id.(int))
	if user != nil && user.AccountID != "" {
		client := &http.Client{}
		accountID := user.AccountID

		req, err := http.NewRequest("GET", os.Getenv("BROKER_API_BASE")+"/v1/trading/accounts/"+accountID+"/account/portfolio/history?period="+url.QueryEscape(c.Query("period"))+"&timeframe="+url.QueryEscape(c.Query("timeframe"))+"&date_end="+url.QueryEscape(c.Query("date_end"))+"&extended_hours="+url.QueryEscape(c.Query("extended_hours"))+"", nil)
		if err != nil {
			fmt.Print(err.Error())
		}

		req.Header.Add("Authorization", os.Getenv("BROKER_TOKEN"))
		response, err := client.Do(req)

		if err != nil {
			fmt.Print(err.Error())
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		var responseObject interface{}
		json.Unmarshal(responseData, &responseObject)

		c.JSON(response.StatusCode, responseObject)
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Couldn't get profile history.",
	})
}

func (a *AccountService) tradingProfile(c *gin.Context) {
	id, _ := c.Get("id")
	user := a.svc.GetProfile(c, id.(int))

	if user.AccountID == "" {
		apperr.Response(c, apperr.New(http.StatusBadRequest, "Account not found."))
		return
	}

	client := &http.Client{}
	accountID := user.AccountID

	req, err := http.NewRequest("GET", os.Getenv("BROKER_API_BASE")+"/v1/trading/accounts/"+accountID+"/account", nil)
	req.Header.Add("Authorization", os.Getenv("BROKER_TOKEN"))

	response, _ := client.Do(req)
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		apperr.Response(c, apperr.New(response.StatusCode, "Something went wrong. Try again later."))
		return
	}

	var responseObject interface{}
	json.Unmarshal(responseData, &responseObject)
	c.JSON(response.StatusCode, responseObject)
}

func (a *AccountService) stats(c *gin.Context) {
	id, _ := c.Get("id")

	peopleInvited, _ := a.db.Model(&model.UserReward{}).Where(`referred_by = ?`, id.(int)).Count()

	referralReward := new(model.UserReward)
	referreReward := new(model.UserReward)
	_, err := a.db.Model((*model.UserReward)(nil)).QueryOne(referralReward, `
		SELECT SUM(reward_value) reward_value from user_rewards where referred_by = ? AND reward_transfer_status = ?;`, id, true)

	_, err1 := a.db.Model((*model.UserReward)(nil)).QueryOne(referreReward, `
	SELECT SUM(reward_value) reward_value from user_rewards where user_id = ? AND reward_transfer_status = ?;`, id, true)
	var totalReward float32 = 0
	if err != nil {
		totalReward = 0
	}

	if err1 != nil {
		totalReward = totalReward + 0
	}

	totalReward = referralReward.RewardValue + referreReward.RewardValue

	c.JSON(http.StatusOK, gin.H{
		"reward_earned":  totalReward,
		"people_invited": peopleInvited,
	})
}

type AssetsResponse struct {
	Class        string `json:"class"`
	EasyToBorrow bool   `json:"easy_to_borrow"`
	Exchange     string `json:"exchange"`
	ID           string `json:"id"`
	Marginable   bool   `json:"marginable"`
	Shortable    bool   `json:"shortable"`
	Status       string `json:"status"`
	Symbol       string `json:"symbol"`
	Tradable     bool   `json:"tradable"`
}
type WatchList struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	AccountID string `json:"account_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Assets    []AssetsResponse
}

func (a *AccountService) getWatchList(c *gin.Context) {
	id, _ := c.Get("id")
	user := a.svc.GetProfile(c, id.(int))
	assets := []interface{}{}

	if user.AccountID == "" {
		apperr.Response(c, apperr.New(http.StatusBadRequest, "Account not found."))
		return
	}

	if user.WatchlistID == "" {
		// var watchlist []string
		c.JSON(http.StatusOK, gin.H{
			"message": "You didn't watchlisted any asset yet.",
			"assets":  []interface{}{},
		})
		return
	}

	client := &http.Client{}
	accountID := user.AccountID

	req, err := http.NewRequest("GET", os.Getenv("BROKER_API_BASE")+"/v1/trading/accounts/"+accountID+"/watchlists/"+user.WatchlistID, nil)
	req.Header.Add("Authorization", os.Getenv("BROKER_TOKEN"))

	response, _ := client.Do(req)
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		apperr.Response(c, apperr.New(response.StatusCode, "Something went wrong. Try again later."))
		return
	}
	var watchlistResponseObject map[string]interface{}
	json.Unmarshal(responseData, &watchlistResponseObject)

	// Get symbol names
	var symbolNames []string
	for _, ass := range watchlistResponseObject["assets"].([]interface{}) {
		ass, _ := ass.(map[string]interface{})
		symbolNames = append(symbolNames, ass["symbol"].(string))
		assets = append(assets, ass)
	}
	// fmt.Println(assets)
	// fmt.Println(symbolNames)

	// fetch market data of assets
	if len(symbolNames) > 0 {
		req2, err := http.NewRequest("GET", os.Getenv("BROKER_API_DATA_BASE")+"/v2/stocks/snapshots?symbols="+url.QueryEscape(strings.Join(symbolNames[:], ",")), nil)
		if err != nil {
			fmt.Print(err.Error())
		}

		req2.Header.Add("Authorization", os.Getenv("BROKER_TOKEN"))
		response2, err := client.Do(req2)

		if err != nil {
			fmt.Print(err.Error())
		}

		responseData2, err := ioutil.ReadAll(response2.Body)
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Printf("%v", string(responseData))

		var responseObject map[string]interface{}
		json.Unmarshal(responseData2, &responseObject)

		for index := range assets {
			assets[index].(map[string]interface{})["ticker"] = responseObject[assets[index].(map[string]interface{})["symbol"].(string)]
		}
		watchlistResponseObject["assets"] = assets
	}
	c.JSON(response.StatusCode, watchlistResponseObject)
}

type Asset struct {
	Symbol string `json:"symbol"`
}

func (a *AccountService) addAssetInWatchList(c *gin.Context) {
	id, _ := c.Get("id")
	user := a.svc.GetProfile(c, id.(int))

	accountID := user.AccountID
	if accountID == "" {
		apperr.Response(c, apperr.New(http.StatusBadRequest, "Account not found."))
		return
	}

	var assets Asset
	c.BindJSON(&assets)

	if assets.Symbol == "" {
		apperr.Response(c, apperr.New(http.StatusBadRequest, "Symbol is required."))
		return
	}

	var symbolList [1]string
	symbolList[0] = assets.Symbol

	if user.WatchlistID == "" {
		client := &http.Client{}
		watchlistJson, _ := json.Marshal(map[string]interface{}{
			"symbols": symbolList,
			"name":    "Watchlist assets",
		})
		watchlistBody := bytes.NewBuffer(watchlistJson)

		fmt.Println(watchlistBody)

		req, err := http.NewRequest("POST", os.Getenv("BROKER_API_BASE")+"/v1/trading/accounts/"+accountID+"/watchlists", watchlistBody)
		req.Header.Add("Authorization", os.Getenv("BROKER_TOKEN"))

		response, _ := client.Do(req)
		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			apperr.Response(c, apperr.New(response.StatusCode, "Something went wrong. Try again later."))
			return
		}

		var responseObject WatchList
		json.Unmarshal(responseData, &responseObject)

		reqUser := request.Update{
			ID:          id.(int),
			AccountID:   &accountID,
			WatchlistID: &responseObject.ID,
		}

		_, error := a.svc.UpdateProfile(c, &reqUser)
		if error != nil {
			apperr.Response(c, apperr.New(http.StatusBadRequest, "Something went wrong. Try again"))
			return
		}

		c.JSON(response.StatusCode, responseObject)
	} else {
		client := &http.Client{}
		watchlistJson, _ := json.Marshal(map[string]interface{}{
			"symbol": assets.Symbol,
		})
		watchlistBody := bytes.NewBuffer(watchlistJson)

		req, err := http.NewRequest("POST", os.Getenv("BROKER_API_BASE")+"/v1/trading/accounts/"+accountID+"/watchlists/"+user.WatchlistID, watchlistBody)
		req.Header.Add("Authorization", os.Getenv("BROKER_TOKEN"))

		response, _ := client.Do(req)
		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			apperr.Response(c, apperr.New(response.StatusCode, "Something went wrong. Try again later."))
			return
		}

		if response.StatusCode != 200 {
			errorBody := ErrorBody{}
			json.Unmarshal(responseData, &errorBody)
			apperr.Response(c, apperr.New(response.StatusCode, errorBody.Message))
			return
		}

		var responseObject interface{}
		json.Unmarshal(responseData, &responseObject)
		c.JSON(response.StatusCode, responseObject)
	}

}

func (a *AccountService) removeAssetFromWatchList(c *gin.Context) {
	id, _ := c.Get("id")
	user := a.svc.GetProfile(c, id.(int))

	if user.AccountID == "" {
		apperr.Response(c, apperr.New(http.StatusBadRequest, "Account not found."))
		return
	}

	if user.WatchlistID == "" {
		apperr.Response(c, apperr.New(http.StatusBadRequest, "Asset not found."))
	}

	symbol := c.Param("symbol")

	client := &http.Client{}
	accountID := user.AccountID

	req, err := http.NewRequest("DELETE", os.Getenv("BROKER_API_BASE")+"/v1/trading/accounts/"+accountID+"/watchlists/"+user.WatchlistID+"/"+symbol, nil)
	req.Header.Add("Authorization", os.Getenv("BROKER_TOKEN"))

	response, _ := client.Do(req)
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		apperr.Response(c, apperr.New(response.StatusCode, "Something went wrong. Try again later."))
		return
	}

	if response.StatusCode != 200 {
		errorBody := ErrorBody{}
		json.Unmarshal(responseData, &errorBody)
		apperr.Response(c, apperr.New(response.StatusCode, errorBody.Message))
		return
	}

	var responseObject interface{}
	json.Unmarshal(responseData, &responseObject)
	c.JSON(response.StatusCode, responseObject)
}

func (a *AccountService) getPositions(c *gin.Context) {
	id, _ := c.Get("id")
	assets := []interface{}{}
	user := a.svc.GetProfile(c, id.(int))
	if user != nil && user.AccountID != "" {
		client := &http.Client{}
		accountID := user.AccountID

		req, err := http.NewRequest("GET", os.Getenv("BROKER_API_BASE")+"/v1/trading/accounts/"+accountID+"/positions", nil)
		if err != nil {
			fmt.Print(err.Error())
		}

		req.Header.Add("Authorization", os.Getenv("BROKER_TOKEN"))
		response, err := client.Do(req)

		if err != nil {
			fmt.Print(err.Error())
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		var responseObject interface{}
		json.Unmarshal(responseData, &responseObject)

		// Get symbol names
		var symbolNames []string
		for _, ass := range responseObject.([]interface{}) {
			ass, _ := ass.(map[string]interface{})
			symbolNames = append(symbolNames, ass["symbol"].(string))

			for _, ass2 := range AssetsList {
				if ass2.Symbol == ass["symbol"] {
					ass["name"] = ass2.Name
					ass["ticker"] = gin.H{}
					ass["is_watchlisted"] = false
				}
			}
			assets = append(assets, ass)

		}

		if len(symbolNames) > 0 {
			client := &http.Client{}

			req, err := http.NewRequest("GET", os.Getenv("BROKER_API_DATA_BASE")+"/v2/stocks/snapshots?symbols="+url.QueryEscape(strings.Join(symbolNames[:], ",")), nil)
			if err != nil {
				fmt.Print(err.Error())
			}

			req.Header.Add("Authorization", os.Getenv("BROKER_TOKEN"))
			response, err := client.Do(req)

			if err != nil {
				fmt.Print(err.Error())
			}

			responseData, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			}
			// fmt.Printf("%v", string(responseData))

			var responseObject map[string]interface{}
			json.Unmarshal(responseData, &responseObject)

			for index := range assets {
				assets[index].(map[string]interface{})["ticker"] = responseObject[assets[index].(map[string]interface{})["symbol"].(string)]
			}

			// Watchlisted flag
			if user.WatchlistID != "" {
				req2, err := http.NewRequest("GET", os.Getenv("BROKER_API_BASE")+"/v1/trading/accounts/"+user.AccountID+"/watchlists/"+user.WatchlistID, nil)
				req2.Header.Add("Authorization", os.Getenv("BROKER_TOKEN"))

				response2, _ := client.Do(req2)
				responseData2, err := ioutil.ReadAll(response2.Body)
				if err != nil {
					apperr.Response(c, apperr.New(response2.StatusCode, "Something went wrong. Try again later."))
					return
				}
				json.Unmarshal(responseData2, &responseObject)
				for index := range assets {
					isWatchlisted := false
					for _, ass := range responseObject["assets"].([]interface{}) {
						ass, _ := ass.(map[string]interface{})
						if ass["symbol"] == assets[index].(map[string]interface{})["symbol"] {
							isWatchlisted = true
							break
						}
					}

					assets[index].(map[string]interface{})["is_watchlisted"] = isWatchlisted
				}
			}
		}

		// fmt.Println(assets)
		// fmt.Println(symbolNames)

		c.JSON(response.StatusCode, assets)
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Couldn't get positions.",
	})
}

func (a *AccountService) getOneOpenPosition(c *gin.Context) {
	id, _ := c.Get("id")
	user := a.svc.GetProfile(c, id.(int))
	if user != nil && user.AccountID != "" {
		client := &http.Client{}
		accountID := user.AccountID

		symbol := c.Param("symbol")

		req, err := http.NewRequest("GET", os.Getenv("BROKER_API_BASE")+"/v1/trading/accounts/"+accountID+"/positions/"+symbol, nil)
		if err != nil {
			fmt.Print(err.Error())
		}

		req.Header.Add("Authorization", os.Getenv("BROKER_TOKEN"))
		response, err := client.Do(req)

		if err != nil {
			fmt.Print(err.Error())
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		var responseObject interface{}
		json.Unmarshal(responseData, &responseObject)

		c.JSON(response.StatusCode, responseObject)
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Couldn't get the position.",
	})
}

func (a *AccountService) closePositions(c *gin.Context) {
	id, _ := c.Get("id")
	user := a.svc.GetProfile(c, id.(int))
	if user != nil && user.AccountID != "" {
		client := &http.Client{}
		accountID := user.AccountID

		req, err := http.NewRequest("DELETE", os.Getenv("BROKER_API_BASE")+"/v1/trading/accounts/"+accountID+"/positions", nil)
		if err != nil {
			fmt.Print(err.Error())
		}

		req.Header.Add("Authorization", os.Getenv("BROKER_TOKEN"))
		response, err := client.Do(req)

		if err != nil {
			fmt.Print(err.Error())
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		var responseObject interface{}
		json.Unmarshal(responseData, &responseObject)

		c.JSON(response.StatusCode, responseObject)
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Couldn't close the position.",
	})
}

func (a *AccountService) closeOnePosition(c *gin.Context) {
	id, _ := c.Get("id")
	user := a.svc.GetProfile(c, id.(int))
	if user != nil && user.AccountID != "" {
		client := &http.Client{}
		accountID := user.AccountID
		symbol := c.Param("symbol")

		req, err := http.NewRequest("DELETE", os.Getenv("BROKER_API_BASE")+"/v1/trading/accounts/"+accountID+"/positions/"+symbol, nil)
		if err != nil {
			fmt.Print(err.Error())
		}

		req.Header.Add("Authorization", os.Getenv("BROKER_TOKEN"))
		response, err := client.Do(req)

		if err != nil {
			fmt.Print(err.Error())
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		var responseObject interface{}
		json.Unmarshal(responseData, &responseObject)

		c.JSON(response.StatusCode, responseObject)
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Couldn't close the positions.",
	})
}

func (a *AccountService) getCalendar(c *gin.Context) {
	id, _ := c.Get("id")
	user := a.svc.GetProfile(c, id.(int))
	if user != nil && user.AccountID != "" {
		client := &http.Client{}
		// accountID := user.AccountID

		req, err := http.NewRequest("GET", os.Getenv("BROKER_API_BASE")+"/v2/calendar", nil)
		if err != nil {
			fmt.Print(err.Error())
		}

		req.Header.Add("Authorization", os.Getenv("BROKER_TOKEN"))
		response, err := client.Do(req)

		if err != nil {
			fmt.Print(err.Error())
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		var responseObject interface{}
		json.Unmarshal(responseData, &responseObject)

		c.JSON(response.StatusCode, responseObject)
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Bad request",
	})
}

func (a *AccountService) getAccount(c *gin.Context) {
	id, _ := c.Get("id")
	user := a.svc.GetProfile(c, id.(int))
	if user != nil && user.AccountID != "" {
		client := &http.Client{}
		accountID := user.AccountID

		req, err := http.NewRequest("GET", os.Getenv("BROKER_API_BASE")+"/v1/trading/accounts/"+accountID+"/account", nil)
		if err != nil {
			fmt.Print(err.Error())
		}

		req.Header.Add("Authorization", os.Getenv("BROKER_TOKEN"))
		response, err := client.Do(req)

		if err != nil {
			fmt.Print(err.Error())
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		var responseObject interface{}
		json.Unmarshal(responseData, &responseObject)

		c.JSON(response.StatusCode, responseObject)
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Couldn't get account details.",
	})
}

func (a *AccountService) getMarketTradesBySymbol(c *gin.Context) {
	id, _ := c.Get("id")
	user := a.svc.GetProfile(c, id.(int))
	if user != nil && user.AccountID != "" {
		client := &http.Client{}
		// accountID := user.AccountID
		symbol := c.Param("symbol")

		req, err := http.NewRequest("GET", os.Getenv("BROKER_API_DATA_BASE")+"/v2/stocks/"+symbol+"/trades?start="+url.QueryEscape(c.Query("start"))+"&end="+url.QueryEscape(c.Query("end"))+"&limit="+url.QueryEscape(c.Query("limit"))+"&page_token="+url.QueryEscape(c.Query("page_token"))+"", nil)
		if err != nil {
			fmt.Print(err.Error())
		}

		req.Header.Add("Authorization", os.Getenv("BROKER_TOKEN"))
		response, err := client.Do(req)

		if err != nil {
			fmt.Print(err.Error())
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		var responseObject interface{}
		json.Unmarshal(responseData, &responseObject)

		c.JSON(response.StatusCode, responseObject)
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Couldn't get market data.",
	})
}

func (a *AccountService) getMarketLatestTradeBySymbol(c *gin.Context) {
	id, _ := c.Get("id")
	user := a.svc.GetProfile(c, id.(int))
	if user != nil && user.AccountID != "" {
		client := &http.Client{}
		// accountID := user.AccountID
		symbol := c.Param("symbol")

		req, err := http.NewRequest("GET", os.Getenv("BROKER_API_DATA_BASE")+"/v2/stocks/"+symbol+"/trades/latest", nil)
		if err != nil {
			fmt.Print(err.Error())
		}

		req.Header.Add("Authorization", os.Getenv("BROKER_TOKEN"))
		response, err := client.Do(req)

		if err != nil {
			fmt.Print(err.Error())
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		var responseObject interface{}
		json.Unmarshal(responseData, &responseObject)

		c.JSON(response.StatusCode, responseObject)
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Couldn't get market data.",
	})
}

func (a *AccountService) getMarketQuotesBySymbol(c *gin.Context) {
	id, _ := c.Get("id")
	user := a.svc.GetProfile(c, id.(int))
	if user != nil && user.AccountID != "" {
		client := &http.Client{}
		// accountID := user.AccountID
		symbol := c.Param("symbol")

		req, err := http.NewRequest("GET", os.Getenv("BROKER_API_DATA_BASE")+"/v2/stocks/"+symbol+"/quotes?start="+url.QueryEscape(c.Query("start"))+"&end="+url.QueryEscape(c.Query("end"))+"&limit="+url.QueryEscape(c.Query("limit"))+"&page_token="+url.QueryEscape(c.Query("page_token"))+"", nil)
		if err != nil {
			fmt.Print(err.Error())
		}

		req.Header.Add("Authorization", os.Getenv("BROKER_TOKEN"))
		response, err := client.Do(req)

		if err != nil {
			fmt.Print(err.Error())
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		var responseObject interface{}
		json.Unmarshal(responseData, &responseObject)

		c.JSON(response.StatusCode, responseObject)
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Couldn't get market data.",
	})
}

func (a *AccountService) getMarketLatestQuoteBySymbol(c *gin.Context) {
	id, _ := c.Get("id")
	user := a.svc.GetProfile(c, id.(int))
	if user != nil && user.AccountID != "" {
		client := &http.Client{}
		// accountID := user.AccountID
		symbol := c.Param("symbol")

		req, err := http.NewRequest("GET", os.Getenv("BROKER_API_DATA_BASE")+"/v2/stocks/"+symbol+"/quotes/latest", nil)
		if err != nil {
			fmt.Print(err.Error())
		}

		req.Header.Add("Authorization", os.Getenv("BROKER_TOKEN"))
		response, err := client.Do(req)

		if err != nil {
			fmt.Print(err.Error())
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		var responseObject interface{}
		json.Unmarshal(responseData, &responseObject)

		c.JSON(response.StatusCode, responseObject)
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Couldn't get market data.",
	})
}

func (a *AccountService) getMarketBarsBySymbol(c *gin.Context) {
	id, _ := c.Get("id")
	user := a.svc.GetProfile(c, id.(int))
	if user != nil {
		client := &http.Client{}
		// accountID := user.AccountID
		symbol := c.Param("symbol")

		req, err := http.NewRequest("GET", os.Getenv("BROKER_API_DATA_BASE")+"/v2/stocks/"+symbol+"/bars?start="+url.QueryEscape(c.Query("start"))+"&end="+url.QueryEscape(c.Query("end"))+"&limit="+url.QueryEscape(c.Query("limit"))+"&page_token="+url.QueryEscape(c.Query("page_token"))+"&timeframe="+url.QueryEscape(c.Query("timeframe"))+"", nil)
		if err != nil {
			fmt.Print(err.Error())
		}

		req.Header.Add("Authorization", os.Getenv("BROKER_TOKEN"))
		response, err := client.Do(req)

		if err != nil {
			fmt.Print(err.Error())
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		var responseObject interface{}
		json.Unmarshal(responseData, &responseObject)

		c.JSON(response.StatusCode, responseObject)
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Couldn't get market data.",
	})
}

func (a *AccountService) getMarketTickers(c *gin.Context) {
	id, _ := c.Get("id")
	user := a.svc.GetProfile(c, id.(int))
	if user != nil && user.AccountID != "" {
		client := &http.Client{}
		// accountID := user.AccountID

		req, err := http.NewRequest("GET", os.Getenv("BROKER_API_DATA_BASE")+"/v2/stocks/snapshots?symbols="+url.QueryEscape(c.Query("symbols")), nil)
		if err != nil {
			fmt.Print(err.Error())
		}

		req.Header.Add("Authorization", os.Getenv("BROKER_TOKEN"))
		response, err := client.Do(req)

		if err != nil {
			fmt.Print(err.Error())
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		var responseObject interface{}
		json.Unmarshal(responseData, &responseObject)

		c.JSON(response.StatusCode, responseObject)
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Couldn't get market data.",
	})
}

func (a *AccountService) getMarketTickerBySymbol(c *gin.Context) {
	id, _ := c.Get("id")
	user := a.svc.GetProfile(c, id.(int))
	if user != nil && user.AccountID != "" {
		client := &http.Client{}
		// accountID := user.AccountID
		symbol := c.Param("symbol")

		req, err := http.NewRequest("GET", os.Getenv("BROKER_API_DATA_BASE")+"/v2/stocks/"+symbol+"/snapshot", nil)
		if err != nil {
			fmt.Print(err.Error())
		}

		req.Header.Add("Authorization", os.Getenv("BROKER_TOKEN"))
		response, err := client.Do(req)

		if err != nil {
			fmt.Print(err.Error())
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		var responseObject interface{}
		json.Unmarshal(responseData, &responseObject)

		c.JSON(response.StatusCode, responseObject)
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Couldn't get market data.",
	})
}
