package config

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Load returns Configuration struct
func Load(env string) *Configuration {
	_, filePath, _, _ := runtime.Caller(0)
	configName := "config." + env + ".yaml"
	configPath := filePath[:len(filePath)-9] + "files" + string(filepath.Separator)

	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	var config Configuration
	viper.Unmarshal(&config)
	setGinMode(config.Server.Mode)

	return &config
}

// Configuration holds data necessery for configuring application
type Configuration struct {
	Server *Server `yaml:"server"`
}

// Server holds data necessary for server configuration
type Server struct {
	Mode string `yaml:"mode"`
}

func setGinMode(mode string) {
	switch mode {
	case "release":
		gin.SetMode(gin.ReleaseMode)
		break
	case "test":
		gin.SetMode(gin.TestMode)
		break
	default:
		gin.SetMode(gin.DebugMode)
	}
}
