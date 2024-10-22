package main

import "github.com/hanoys/sigma-music/internal/app/web"

// @title           Sigma Music API
// @version         1.0
// @description     This is a Sigma Music API.
// @termsOfService  http://swagger.io/terms/

// @host      localhost:80
// @BasePath  /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	web.Run()
}
