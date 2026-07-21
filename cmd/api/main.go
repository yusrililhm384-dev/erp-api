package main

import "github.com/labstack/echo/v5"

func main() {
	e := echo.New()

	if err := e.Start(":3000"); err != nil {
		e.Logger.Error(err.Error())
	}
}
