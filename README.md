# Aruna API-Gateway

API Gateway for the [**Aruna API**](https://github.com/ArunaStorage/api) based on gRPC-Gateway. This translates RESTful http calls into native gRPC API calls. 

## Configuration parameters

The configuration for the Aruna API gRPC-Gateway has to be described as yaml file and it has to be placed under `./config`.

If you want to use the config for a local Aruna testing deployment you have to switch the config files:
```bash
mv config/config.yaml config/config-prod.yaml
mv config/config-local.yaml config/config.yaml
go run main.go
```

### Server parameters

| Name          | Description | Value  |
| ------------- | ----------- | ------ |
| `Server.Port` | Server port | `8080` |

### SWAGGER parameters

| Name           | Description                       | Value         |
| -------------- | --------------------------------- | ------------- |
| `Swagger.Path` | Path to the swagger/openAPI files | `www/swagger` |

### Backend parameters

| Name           | Description         | Value       |
| -------------- | ------------------- | ----------- |
| `BACKEND.HOST` | Backend server host | `127.0.0.1` |
| `BACKEND.PORT` | Backend server port | `50051`     |
