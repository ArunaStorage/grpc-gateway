package gateway

import (
	"log"

	"context"
	"fmt"
	"net/http"

	"github.com/ScienceObjectsDB/CORE-API-Gateway/config"
	"github.com/ScienceObjectsDB/go-api/openapiv2"
	v1storageservices "github.com/ScienceObjectsDB/go-api/sciobjsdb/api/storage/services/v1"
	"github.com/gin-contrib/cors"
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

	defaultCors := cors.DefaultConfig()
	defaultCors.AllowAllOrigins = true
	defaultCors.AddAllowHeaders("grpc-metadata-accesstoken")

	r.Use(cors.New(defaultCors))

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

	err := v1storageservices.RegisterProjectServiceHandlerFromEndpoint(ctx, gwmux, fmt.Sprintf("%v:%v", grpcEndpointHost, grpcEndpointPort), opts)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	err = v1storageservices.RegisterDatasetServiceHandlerFromEndpoint(ctx, gwmux, fmt.Sprintf("%v:%v", grpcEndpointHost, grpcEndpointPort), opts)
	if err != nil {
		log.Println(err.Error())

		return err
	}

	err = v1storageservices.RegisterDatasetObjectsServiceHandlerFromEndpoint(ctx, gwmux, fmt.Sprintf("%v:%v", grpcEndpointHost, grpcEndpointPort), opts)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	err = v1storageservices.RegisterObjectLoadServiceHandlerFromEndpoint(ctx, gwmux, fmt.Sprintf("%v:%v", grpcEndpointHost, grpcEndpointPort), opts)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	port := viper.GetInt(config.SERVER_PORT)

	return r.Run(fmt.Sprintf(":%v", port))
}
