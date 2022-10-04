package gateway

import (
	"log"

	"context"
	"fmt"
	"net/http"

	v1storageservices "github.com/ArunaStorage/go-api/aruna/api/storage/services/v1"
	"github.com/ArunaStorage/go-api/openapiv2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// StartETLGateway Starts the gateway server for the ETL component
func StartGateway() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	fmt.Println(viper.AllSettings())
	gwmux := runtime.NewServeMux()

	grpcEndpointHost := viper.GetString("config.gateway.grpcendpointhost")
	grpcEndpointPort := viper.GetInt("config.gateway.grpcendpointport")

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	r := gin.Default()

	defaultCors := cors.DefaultConfig()
	defaultCors.AllowAllOrigins = true
	defaultCors.AddAllowHeaders("Authorization")

	r.Use(cors.New(defaultCors))

	r.Any("/v1/*any", gin.WrapF(gwmux.ServeHTTP))

	swagger_fs := http.FS(openapiv2.GetSwaggerEmbedded())
	r.StaticFS("/swaggerjson", swagger_fs)

	swagger_files := viper.GetString("config.swagger.path")

	fs := http.FileSystem(http.Dir(swagger_files))

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/swagger-ui/")
	})

	r.StaticFS("/swagger-ui/", fs)

	err := v1storageservices.RegisterProjectServiceHandlerFromEndpoint(ctx, gwmux, fmt.Sprintf("%v:%v", grpcEndpointHost, grpcEndpointPort), opts)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	err = v1storageservices.RegisterCollectionServiceHandlerFromEndpoint(ctx, gwmux, fmt.Sprintf("%v:%v", grpcEndpointHost, grpcEndpointPort), opts)
	if err != nil {
		log.Println(err.Error())

		return err
	}

	err = v1storageservices.RegisterObjectServiceHandlerFromEndpoint(ctx, gwmux, fmt.Sprintf("%v:%v", grpcEndpointHost, grpcEndpointPort), opts)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	err = v1storageservices.RegisterObjectGroupServiceHandlerFromEndpoint(ctx, gwmux, fmt.Sprintf("%v:%v", grpcEndpointHost, grpcEndpointPort), opts)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	err = v1storageservices.RegisterEndpointServiceHandlerFromEndpoint(ctx, gwmux, fmt.Sprintf("%v:%v", grpcEndpointHost, grpcEndpointPort), opts)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	err = v1storageservices.RegisterUserServiceHandlerFromEndpoint(ctx, gwmux, fmt.Sprintf("%v:%v", grpcEndpointHost, grpcEndpointPort), opts)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	port := viper.GetInt("config.server.port")

	return r.Run(fmt.Sprintf(":%v", port))
}
