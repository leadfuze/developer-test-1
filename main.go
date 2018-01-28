package main

import (
	"github.com/labstack/echo"
	"github.com/yogihardi/developer-test-1/externalservice"
)

func main() {
	e := echo.New()
	server := NewServer(e, &externalservice.ClientImpl{})
	server.AddRoutes()
	server.Run(8080)
}
