package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/alpacahq/ribbit-backend/config"
	"github.com/alpacahq/ribbit-backend/model"
	"github.com/alpacahq/ribbit-backend/repository"
	"github.com/alpacahq/ribbit-backend/secret"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// syncAssetsCmd represents the syncAssets command
var syncAssetsCmd = &cobra.Command{
	Use:   "sync_assets",
	Short: "sync_assets sync all the assets from broker",
	Long:  `sync_assets sync all the assets from broker`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("syncAssets called")
		db := config.GetConnection()
		log, _ := zap.NewDevelopment()
		defer log.Sync()
		assetRepo := repository.NewAssetRepo(db, log, secret.New())

		client := &http.Client{}

		req, err := http.NewRequest("GET", os.Getenv("BROKER_API_BASE")+"/v1/assets", nil)
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
			log.Fatal(err.Error())
		}
		// fmt.Printf("%v", string(responseData))

		var responseObject []interface{}
		json.Unmarshal(responseData, &responseObject)

		for _, asset := range responseObject {
			asset, _ := asset.(map[string]interface{})

			newAsset := new(model.Asset)
			newAsset.ID = asset["id"].(string)
			newAsset.Class = asset["class"].(string)
			newAsset.Exchange = asset["exchange"].(string)
			newAsset.Symbol = asset["symbol"].(string)
			newAsset.Name = asset["name"].(string)
			newAsset.Status = asset["status"].(string)
			newAsset.Tradable = asset["tradable"].(bool)
			newAsset.Marginable = asset["marginable"].(bool)
			newAsset.Shortable = asset["shortable"].(bool)
			newAsset.EasyToBorrow = asset["easy_to_borrow"].(bool)
			newAsset.Fractionable = asset["fractionable"].(bool)

			if _, err := assetRepo.CreateOrUpdate(newAsset); err != nil {
				log.Fatal(err.Error())
			} else {
				// fmt.Println(asset)
			}
		}

		// m := manager.NewManager(accountRepo, roleRepo, db)
		// models := manager.GetModels()
		// m.CreateSchema(models...)
		// m.CreateRoles()
	},
}

func init() {
	rootCmd.AddCommand(syncAssetsCmd)
}
