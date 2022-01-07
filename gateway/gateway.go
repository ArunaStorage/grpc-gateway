package gateway

import (
	"log"

	"context"
	"fmt"
	"net/http"

	"github.com/ScienceObjectsDB/CORE-API-Gateway/config"
	service "github.com/ScienceObjectsDB/go-api/api/services/v1"
	"github.com/ScienceObjectsDB/go-api/openapiv2"
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

	grpcEndpointHost := viper.GetString(config.BACKEND_HOST)
	grpcEndpointPort := viper.GetInt(config.BACKEND_PORT)

	opts := []grpc.DialOption{grpc.WithInsecure()}

	r := gin.Default()

	r.Any("/api/*any", gin.WrapF(gwmux.ServeHTTP))

	swagger_fs := http.FS(openapiv2.GetSwaggerEmbedded())
	r.StaticFS("/swaggerjson", swagger_fs)

	swagger_files := viper.GetString(config.SWAGGER_PATH)

	fs := http.FileSystem(http.Dir(swagger_files))

	r.GET("/oauth2-redirect.html", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "/swagger-ui/oauth2-redirect.html")
	})

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/swagger-ui/")
	})

	r.StaticFS("/swagger-ui/", fs)

	err := service.RegisterProjectServiceHandlerFromEndpoint(ctx, gwmux, fmt.Sprintf("%v:%v", grpcEndpointHost, grpcEndpointPort), opts)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	err = service.RegisterDatasetServiceHandlerFromEndpoint(ctx, gwmux, fmt.Sprintf("%v:%v", grpcEndpointHost, grpcEndpointPort), opts)
	if err != nil {
		log.Println(err.Error())

		return err
	}

	err = service.RegisterDatasetObjectsServiceHandlerFromEndpoint(ctx, gwmux, fmt.Sprintf("%v:%v", grpcEndpointHost, grpcEndpointPort), opts)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	err = service.RegisterObjectLoadServiceHandlerFromEndpoint(ctx, gwmux, fmt.Sprintf("%v:%v", grpcEndpointHost, grpcEndpointPort), opts)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	port := viper.GetInt(config.SERVER_PORT)

	return r.Run(fmt.Sprintf(":%v", port))
}
