package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/alpacahq/ribbit-backend/apperr"
	account "github.com/alpacahq/ribbit-backend/repository/account"
	"github.com/alpacahq/ribbit-backend/repository/assets"

	"github.com/gin-gonic/gin"
)

func AssetsRouter(svc *assets.Service, acc *account.Service, r *gin.RouterGroup) {
	a := Assets{svc, acc}

	ar := r.Group("/assets")
	ar.GET("/", a.getAssetsList)
	ar.GET("/:id", a.getAssetDetail)

}

// Auth represents auth http service
type Assets struct {
	svc *assets.Service
	acc *account.Service
}

type AssetObj struct {
	ID            string      `json:"id"`
	Class         string      `json:"class"`
	Exchange      string      `json:"exchange"`
	Symbol        string      `json:"symbol"`
	Name          string      `json:"name"`
	Status        string      `json:"status"`
	Tradable      bool        `json:"tradable"`
	Marginable    bool        `json:"marginable"`
	Shortable     bool        `json:"shortable"`
	EasyToBorrow  bool        `json:"easy_to_borrow"`
	Fractionable  bool        `json:"fractionable"`
	Ticker        interface{} `json:"ticker"`
	IsWatchlisted bool        `json:"is_watchlisted"`
}

type SymbolTicker struct {
	DailyBar struct {
		C float64   `json:"c"`
		H float64   `json:"h"`
		L float64   `json:"l"`
		O float64   `json:"o"`
		T time.Time `json:"t"`
		V int       `json:"v"`
	} `json:"dailyBar"`
	LatestQuote struct {
		Ap int       `json:"ap"`
		As int       `json:"as"`
		Ax string    `json:"ax"`
		Bp int       `json:"bp"`
		Bs int       `json:"bs"`
		Bx string    `json:"bx"`
		C  []string  `json:"c"`
		T  time.Time `json:"t"`
	} `json:"latestQuote"`
	LatestTrade struct {
		C []string  `json:"c"`
		I int       `json:"i"`
		P float64   `json:"p"`
		S int       `json:"s"`
		T time.Time `json:"t"`
		X string    `json:"x"`
		Z string    `json:"z"`
	} `json:"latestTrade"`
	MinuteBar struct {
		C float64   `json:"c"`
		H float64   `json:"h"`
		L float64   `json:"l"`
		O float64   `json:"o"`
		T time.Time `json:"t"`
		V int       `json:"v"`
	} `json:"minuteBar"`
	PrevDailyBar struct {
		C int       `json:"c"`
		H float64   `json:"h"`
		L float64   `json:"l"`
		O float64   `json:"o"`
		T time.Time `json:"t"`
		V int       `json:"v"`
	} `json:"prevDailyBar"`
	Symbol string `json:"symbol"`
}

var AssetsList = []AssetObj{
	AssetObj{
		ID:            "b0b6dd9d-8b9b-48a9-ba46-b9d54906e415",
		Class:         "us_equity",
		Exchange:      "NASDAQ",
		Symbol:        "AAPL",
		Name:          "Apple Inc. Common Stock",
		Status:        "active",
		Tradable:      true,
		Marginable:    true,
		Shortable:     true,
		EasyToBorrow:  true,
		Fractionable:  true,
		IsWatchlisted: false,
	},
	AssetObj{
		ID:            "f30d734c-2806-4d0d-b145-f9fade61432b",
		Class:         "us_equity",
		Exchange:      "NASDAQ",
		Symbol:        "GOOG",
		Name:          "Alphabet Inc. Class C Capital Stock",
		Status:        "active",
		Tradable:      true,
		Marginable:    true,
		Shortable:     true,
		EasyToBorrow:  true,
		Fractionable:  true,
		IsWatchlisted: false,
	},
	AssetObj{
		ID:            "69b15845-7c63-4586-b274-1cfdfe9df3d8",
		Class:         "us_equity",
		Exchange:      "NASDAQ",
		Symbol:        "GOOGL",
		Name:          "Alphabet Inc. Class A Common Stock",
		Status:        "active",
		Tradable:      true,
		Marginable:    true,
		Shortable:     true,
		EasyToBorrow:  true,
		Fractionable:  true,
		IsWatchlisted: false,
	},
	AssetObj{
		ID:            "fc6a5dcd-4a70-4b8d-b64f-d83a6dae9ba4",
		Class:         "us_equity",
		Exchange:      "NASDAQ",
		Symbol:        "FB",
		Name:          "Facebook, Inc. Class A Common Stock",
		Status:        "active",
		Tradable:      true,
		Marginable:    true,
		Shortable:     true,
		EasyToBorrow:  true,
		Fractionable:  true,
		IsWatchlisted: false,
	},
	AssetObj{
		ID:            "39a26dc1-927a-4590-b103-b8068a013e7f",
		Class:         "us_equity",
		Exchange:      "NYSE",
		Symbol:        "SPOT",
		Name:          "Spotify Technology S.A.",
		Status:        "active",
		Tradable:      true,
		Marginable:    true,
		Shortable:     true,
		EasyToBorrow:  true,
		Fractionable:  true,
		IsWatchlisted: false,
	},
	AssetObj{
		ID:            "83e52ac1-bb18-4e9f-b68d-dda5a8af3ec0",
		Class:         "us_equity",
		Exchange:      "NYSE",
		Symbol:        "SNAP",
		Name:          "Snap Inc.",
		Status:        "active",
		Tradable:      true,
		Marginable:    true,
		Shortable:     true,
		EasyToBorrow:  true,
		Fractionable:  true,
		IsWatchlisted: false,
	},
	AssetObj{
		ID:            "8ccae427-5dd0-45b3-b5fe-7ba5e422c766",
		Class:         "us_equity",
		Exchange:      "NASDAQ",
		Symbol:        "TSLA",
		Name:          "Tesla, Inc. Common Stock",
		Status:        "active",
		Tradable:      true,
		Marginable:    true,
		Shortable:     true,
		EasyToBorrow:  true,
		Fractionable:  true,
		IsWatchlisted: false,
	},
	AssetObj{
		ID:            "57c36644-876b-437c-b913-3cdb58b18fd3",
		Class:         "us_equity",
		Exchange:      "NYSE",
		Symbol:        "GE",
		Name:          "General Electric Company",
		Status:        "active",
		Tradable:      true,
		Marginable:    true,
		Shortable:     true,
		EasyToBorrow:  true,
		Fractionable:  true,
		IsWatchlisted: false,
	},
	AssetObj{
		ID:            "4f5baf1e-0e9b-4d85-b88a-d874dc4a3c42",
		Class:         "us_equity",
		Exchange:      "NYSE",
		Symbol:        "V",
		Name:          "VISA Inc.",
		Status:        "active",
		Tradable:      true,
		Marginable:    true,
		Shortable:     true,
		EasyToBorrow:  true,
		Fractionable:  true,
		IsWatchlisted: false,
	},
	AssetObj{
		ID:            "2140998d-7f62-46f2-a9b2-e44350bd4807",
		Class:         "us_equity",
		Exchange:      "NYSE",
		Symbol:        "MA",
		Name:          "Mastercard Incorporated",
		Status:        "active",
		Tradable:      true,
		Marginable:    true,
		Shortable:     true,
		EasyToBorrow:  true,
		Fractionable:  true,
		IsWatchlisted: false,
	},
	AssetObj{
		ID:            "f801f835-bfe6-4a9d-a6b1-ccbb84bfd75f",
		Class:         "us_equity",
		Exchange:      "NASDAQ",
		Symbol:        "AMZN",
		Name:          "Amazon.com, Inc. Common Stock",
		Status:        "active",
		Tradable:      true,
		Marginable:    true,
		Shortable:     true,
		EasyToBorrow:  true,
		Fractionable:  true,
		IsWatchlisted: false,
	},
	AssetObj{
		ID:            "bb2a26c0-4c77-4801-8afc-82e8142ac7b8",
		Class:         "us_equity",
		Exchange:      "NASDAQ",
		Symbol:        "NFLX",
		Name:          "Netflix, Inc. Common Stock",
		Status:        "active",
		Tradable:      true,
		Marginable:    true,
		Shortable:     true,
		EasyToBorrow:  true,
		Fractionable:  true,
		IsWatchlisted: false,
	},
	AssetObj{
		ID:            "662a919f-1455-497c-90e7-f76248e6d3a6",
		Class:         "us_equity",
		Exchange:      "NYSE",
		Symbol:        "TME",
		Name:          "Tencent Music Entertainment Group American Depositary Shares, each representing two Class A Ordinary",
		Status:        "active",
		Tradable:      true,
		Marginable:    true,
		Shortable:     true,
		EasyToBorrow:  true,
		Fractionable:  true,
		IsWatchlisted: false,
	},
}

// var baseURL string = "https://broker-api.sandbox.alpaca.markets"
// var baseDataURL string = "https://data.sandbox.alpaca.markets"
// var token string = "Basic Q0tWTU9JRUNPQk9PWVMwMEZKRVQ6MW9HcXFkdGNSWkRNeUhuWkg1d1N4dE94SEswbXZDWnlyOUdZUHlSUw=="

func (a *Assets) getAssetsList(c *gin.Context) {
	q := c.Query("q")
	_assets := []AssetObj{}
	id, _ := c.Get("id")
	user := a.acc.GetProfile(c, id.(int))

	if len(q) > 0 {
		searchedAssets, err := a.svc.SearchAssets(q)
		if err != nil {
			c.JSON(http.StatusOK, _assets)
			return
		}
		for _, searchedAsset := range searchedAssets {
			var _ass AssetObj
			_ass.ID = searchedAsset.ID
			_ass.Class = searchedAsset.Class
			_ass.Exchange = searchedAsset.Exchange
			_ass.Symbol = searchedAsset.Symbol
			_ass.Name = searchedAsset.Name
			_ass.Status = searchedAsset.Status
			_ass.Tradable = searchedAsset.Tradable
			_ass.Marginable = searchedAsset.Marginable
			_ass.Shortable = searchedAsset.Shortable
			_ass.EasyToBorrow = searchedAsset.EasyToBorrow
			_ass.Fractionable = searchedAsset.Fractionable
			_ass.IsWatchlisted = searchedAsset.IsWatchlisted
			_assets = append(_assets, _ass)
		}
	} else {
		for _, asset := range AssetsList {
			if len(q) > 0 {
				if strings.Contains(strings.ToLower(asset.Symbol), strings.ToLower(q)) {
					_assets = append(_assets, asset)
				}
			} else {
				_assets = append(_assets, asset)
			}
		}
	}

	// get symbol names list
	var symbolNames []string
	for _, asset := range _assets {
		symbolNames = append(symbolNames, asset.Symbol)
	}

	// fetch market data of _assets
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

		for index := range _assets {
			_assets[index].Ticker = responseObject[_assets[index].Symbol]
		}

		// Watchlisted flag
		if user.AccountID != "" && user.WatchlistID != "" {
			req2, err := http.NewRequest("GET", os.Getenv("BROKER_API_BASE")+"/v1/trading/accounts/"+user.AccountID+"/watchlists/"+user.WatchlistID, nil)
			req2.Header.Add("Authorization", os.Getenv("BROKER_TOKEN"))

			response2, _ := client.Do(req2)
			responseData2, err := ioutil.ReadAll(response2.Body)
			if err != nil {
				apperr.Response(c, apperr.New(response2.StatusCode, "Something went wrong. Try again later."))
				return
			}
			json.Unmarshal(responseData2, &responseObject)
			for index := range _assets {
				isWatchlisted := false
				for _, ass := range responseObject["assets"].([]interface{}) {
					ass, _ := ass.(map[string]interface{})
					if ass["symbol"] == _assets[index].Symbol {
						isWatchlisted = true
						break
					}
				}

				_assets[index].IsWatchlisted = isWatchlisted
			}
		}
	}

	c.JSON(http.StatusOK, _assets)
}

func (a *Assets) getAssetDetail(c *gin.Context) {
	id := c.Param("id")

	for i := range AssetsList {
		if AssetsList[i].ID == id {
			// fetch market data of assets
			client := &http.Client{}

			req, err := http.NewRequest("GET", os.Getenv("BROKER_API_BASE")+"/v2/stocks/snapshots?symbols="+url.QueryEscape(AssetsList[i].Symbol), nil)
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

			AssetsList[i].Ticker = responseObject[AssetsList[i].Symbol]
			c.JSON(200, AssetsList[i])
		}
	}
	c.JSON(http.StatusNotFound, gin.H{
		"message": "400 Not found",
	})
}
