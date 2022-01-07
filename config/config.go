package config

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	BACKEND_HOST = "Gateway.GRPCEndpointHost"
	BACKEND_PORT = "Gateway.GRPCEndpointPort"

	SWAGGER_PATH = "Swagger.Path"

	SERVER_PORT = "Server.Port"
)

func HandleConfigFile(configFile string) {
	SetDefaults()

	viper.SetConfigFile(configFile)
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}

func SetDefaults() {
	viper.SetDefault(BACKEND_HOST, "127.0.0.1")
	viper.SetDefault(BACKEND_PORT, 50051)
	viper.SetDefault(SWAGGER_PATH, "www/swagger")
	viper.SetDefault(SERVER_PORT, 8080)
}
