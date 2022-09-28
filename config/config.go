package config

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	BACKEND_HOST      = "Gateway.GRPCEndpointHost"
	BACKEND_PORT      = "Gateway.GRPCEndpointPort"
	SWAGGER_PATH      = "Swagger.Path"
	SERVER_PORT       = "Server.Port"
	SERVER_BASE_URL   = "Server.BaseUrl"
	KEYCLOAK_URL      = "KeyCloak.Url"
	KEYCLOAK_CLIENTID = "KeyCloak.ClientId"
	KEYCLOAK_SECRET   = "KeyCloak.Secret"
)

func HandleConfigFile() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/config")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	fmt.Println(viper.AllSettings())
}
