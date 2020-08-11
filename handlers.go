package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/skip2/go-qrcode"
)

func parseIntWithDefault(input string, def int) (int, error) {
	if input == "" {
		return def, nil
	}

	n, err := strconv.Atoi(input)
	if err != nil {
		return def, echo.NewHTTPError(http.StatusBadRequest, "parameter must be an integer")
	}
	return n, nil
}

// err400 returns 400 and is used to discourage random queries
func err400(c echo.Context) error {
	return c.String(http.StatusBadRequest, "Go away.")
}

// doc returns a documentation page
func doc(c echo.Context) error {
	doctext := `
	<h1>QR code generator</h1>
	<p>This service generates QR codes. A GET to /qr with a query
	parameter called <b>url</b> will generate a QR code as a png
	with the contents of the QR code being the body of the url string.</p>
	`
	return c.String(http.StatusOK, doctext)
}

// health returns 200 Ok and can be used by a load balancer to indicate
// that the service is stable
func health(c echo.Context) error {
	return c.String(http.StatusOK, "Ok\n")
}

// qrcodegen is a handler that returns a png image of a QR code
//
// Required query parameter is url, which is used as the body of the QR code
//
// Optional query parameters are:
// * size is a number of the pixel size of the png; default is 512.
// * level is the recovery level - options are "l" (low), "m" (medium -- default), "h" (high), "x" (max)
func qrcodegen(c echo.Context) error {
	url := c.QueryParam("url")
	if url == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "URL query parameter required")
	}

	level := qrcode.Medium
	s := c.QueryParam("level")
	switch s {
	case "l":
		level = qrcode.Low
	case "m":
		level = qrcode.Medium
	case "h":
		level = qrcode.High
	case "x":
		level = qrcode.Highest
	case "":
	// do nothing
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "level parameter must be one of l,m,h,x")
	}

	size, err := parseIntWithDefault(c.QueryParam("size"), 256)
	if err != nil {
		return err
	}
	if size < 128 || size > 1024 {
		return echo.NewHTTPError(http.StatusBadRequest, "parameter must be between 128 and 1024")
	}

	png, err := qrcode.Encode(url, level, size)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "could not encode that URL")
	}
	return c.Blob(http.StatusOK, "image/png", png)
}
