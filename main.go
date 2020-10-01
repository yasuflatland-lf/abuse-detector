package main

import (
	"net/http"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/op/go-logging"
	"studio.design/studio-abuse-detector/pkg/verify"
)

var log = logging.MustGetLogger("verify")

var logFmt = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} PID=%{pid} MOD=%{module} PKG=%{shortpkg} %{shortfile} FUNC=%{shortfunc} ▶ %{level:.4s} %{id:03x} %{color:reset} %{message}`,
)

func main() {
	err := godotenv.Load()
	if err != nil {
		// log.Fatal("Error loading .env file")
		log.Error(err)
	}

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	e.GET("/verify", Verify)

	// Start server
	e.Logger.Fatal(e.Start(":3000"))
}

type VerifyResponse struct {
	Link      string `json:"link" xml:"link"`
	Malicious bool   `json:"malicious" xml:"malicious"`
	Error     error  `json:"error" xml:"error"`
}

// Handler
func Verify(c echo.Context) error {
	url := c.QueryParam("url")
	vr := &VerifyResponse{
		Link:      "",
		Malicious: false,
		Error:     nil,
	}

	ret, link, err := verify.Do(url)

	if err != nil {
		log.Error(err)
		vr.Error = err
	} else {
		vr.Link = link
		vr.Malicious = ret
	}

	return c.JSON(http.StatusOK, vr)
}
