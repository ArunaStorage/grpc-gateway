# CORE-API-Gateway

API Gateway for the CORE-API based on gRPC-Gateway

## Configuration parameters

The configuration for the CORE-Gateway has to be described as yaml file. It has to be places under ./config.

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
