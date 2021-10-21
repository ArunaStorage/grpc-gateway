package gateway

import (
	"log"

	"context"
	"fmt"
	"net/http"

	api "github.com/ScienceObjectsDB/go-api/api/services/v1"
	openapiv2 "github.com/ScienceObjectsDB/go-api/openapiv2"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"

	"google.golang.org/grpc"
)

// StartETLGateway Starts the gateway server for the ETL component
func StartGateway() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	gwmux := runtime.NewServeMux()

	grpcEndpointHost := viper.GetString("Config.Gateway.GRPCEndpointHost")
	grpcEndpointPort := viper.GetInt("Config.Gateway.GRPCEndpointPort")

	opts := []grpc.DialOption{grpc.WithInsecure()}

	r := gin.Default()

	r.Any("/api/*any", gin.WrapF(gwmux.ServeHTTP))

	swagger_fs := http.FS(openapiv2.GetSwaggerEmbedded())
	r.StaticFS("/swaggerjson", swagger_fs)

	swagger_files := viper.GetString("Config.Swagger.Path")

	fs := http.FileSystem(http.Dir(swagger_files))

	r.GET("/oauth2-redirect.html", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "/swagger-ui/oauth2-redirect.html")
	})

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/swagger-ui/")
	})

	r.StaticFS("/swagger-ui/", fs)

	err := api.RegisterProjectServiceHandlerFromEndpoint(ctx, gwmux, fmt.Sprintf("%v:%v", grpcEndpointHost, grpcEndpointPort), opts)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	err = api.RegisterDatasetServiceHandlerFromEndpoint(ctx, gwmux, fmt.Sprintf("%v:%v", grpcEndpointHost, grpcEndpointPort), opts)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	err = api.RegisterDatasetObjectsServiceHandlerFromEndpoint(ctx, gwmux, fmt.Sprintf("%v:%v", grpcEndpointHost, grpcEndpointPort), opts)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	err = api.RegisterObjectLoadServiceHandlerFromEndpoint(ctx, gwmux, fmt.Sprintf("%v:%v", grpcEndpointHost, grpcEndpointPort), opts)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	port := viper.GetInt("Config.Gateway.Port")

	return r.Run(fmt.Sprintf(":%v", port))
}
