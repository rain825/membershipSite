package main

import (
	"net/http"

	"github.com/labstack/echo"
)

// Handler
func helloHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func main() {
	// Echo instance
	e := echo.New()

	// Routes
	e.GET("/", helloHandler)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
