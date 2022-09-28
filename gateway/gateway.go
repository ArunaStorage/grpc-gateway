package gateway

import (
	"log"

	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ArunaStorage/Aruna-REST-Gateway/config"
	v1storageservices "github.com/ArunaStorage/go-api/aruna/api/storage/services/v1"
	"github.com/ArunaStorage/go-api/openapiv2"
	"github.com/coreos/go-oidc"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"

	"google.golang.org/grpc"
)

func CustomHeaderMatcher(key string) (string, bool) {
	switch key {
	case "Authorization":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}

// StartETLGateway Starts the gateway server for the ETL component
func StartGateway() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	fmt.Println(viper.AllSettings())
	gwmux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(CustomHeaderMatcher))

	grpcEndpointHost := viper.GetString(config.BACKEND_HOST)
	grpcEndpointPort := viper.GetInt(config.BACKEND_PORT)

	opts := []grpc.DialOption{grpc.WithInsecure()}

	r := gin.Default()

	issuer, err := url.Parse(viper.GetString("config.keycloak.url"))

	log.Printf("The issues urls is: %v", issuer.String())
	log.Printf("The issues urls-original is: %v", viper.GetString(config.KEYCLOAK_URL))
	if err != nil {
		log.Fatal(err)
	}

	clurl, err := url.Parse(viper.GetString(config.SERVER_BASE_URL))

	if err != nil {
		log.Fatal(err)
	}

	oauthz := OidcHandler{
		Router:       r,
		ClientId:     viper.GetString(config.KEYCLOAK_CLIENTID),
		ClientSecret: viper.GetString(config.KEYCLOAK_SECRET),
		Issuer:       *issuer,                                        //the URL identifier for the authorization service. for example: "https://accounts.google.com" - try adding "/.well-known/openid-configuration" to the path to make sure it's correct
		ClientUrl:    *clurl,                                         //your website's/service's URL for example: "http://localhost:8081/" or "https://mydomain.com/
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"}, //OAuth scopes. If you're unsure go with: []string{oidc.ScopeOpenID, "profile", "email"}
		Config:       nil,
	}

	oauthz.Init()

	defaultCors := cors.DefaultConfig()
	defaultCors.AllowAllOrigins = true
	defaultCors.AddAllowHeaders("Authorization")

	r.Use(cors.New(defaultCors))

	r.Any("/v1/*any", gin.WrapF(gwmux.ServeHTTP))

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

	fs_ui := http.FileSystem(http.Dir("www/ui"))
	r.StaticFS("/ui", fs_ui)

	err = v1storageservices.RegisterProjectServiceHandlerFromEndpoint(ctx, gwmux, fmt.Sprintf("%v:%v", grpcEndpointHost, grpcEndpointPort), opts)
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

	port := viper.GetInt(config.SERVER_PORT)

	return r.Run(fmt.Sprintf(":%v", port))
}
