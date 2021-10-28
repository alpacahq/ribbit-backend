package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/alpacahq/ribbit-backend/route"
	"github.com/alpacahq/ribbit-backend/server"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// routes will be attached to s
var s server.Server

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "alpaca",
	Short: "Broker API middleware",
	Long:  `Broker MVP that uses golang gin as webserver, and go-pg library for connecting with a PostgreSQL database.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		var env string
		var ok bool
		if env, ok = os.LookupEnv("ALPACA_ENV"); !ok {
			env = "dev"
			fmt.Printf("Run server in %s mode\n", env)
		}
		err := s.Run(env)
		if err != nil {
			log.Fatal(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(customRouteServices []route.ServicesI) {
	s.RouteServices = customRouteServices
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err, "sdjsbhfjbhsfb")
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.alpaca.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".alpaca" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".alpaca")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
