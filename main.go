// Package main provides a simple QR code service that returns a
// QR code for the string passed in the URL parameter.
package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", err400)
	e.GET("/doc", doc)
	e.GET("/health", health)
	e.GET("/qr", qrcodegen)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
