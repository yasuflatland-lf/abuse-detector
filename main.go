package main

import (
	"context"
	"net/http"
	"time"

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

type VerifyResponse struct {
	StrategyName string   `json:"strategyName" xml:"strategyName"`
	Links        []string `json:"link" xml:"link"`
	Malicious    bool     `json:"malicious" xml:"malicious"`
	StatusCode   int      `json:"statusCode" xml:"statusCode"`
	Error        error    `json:"error" xml:"error"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		// log.Fatal("Error loading .env file")
		log.Error(err)
	}

	router := NewRouter()

	// Start server
	router.Logger.Fatal(router.Start(":3000"))
}

func NewRouter() *echo.Echo {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	e.GET("/verify", Verify)

	return e
}

// Handler
func Verify(c echo.Context) error {
	url := c.QueryParam("url")

	strategies := []verify.Verify{
		verify.NewTransparencyReportVerifyStrategy(),
		// verify.NewUrlScanVerifyStrategy(),
	}

	// Set up channels
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	errCh := make(chan VerifyResponse, len(strategies))
	retCh := make(chan VerifyResponse, len(strategies))

	for _, s := range strategies {

		// Run each verification concurrent.
		// Results are verified by the response order.
		// Return the result as soon as the site is confirmed,
		// including malicious links.
		go func(strategy verify.Verify) {
			vr := &VerifyResponse{
				Links:     []string{},
				Malicious: false,
				Error:     nil,
			}

			ret, err := strategy.Do(ctx, url)

			vr.Links = ret.MaliciousLinks
			vr.StrategyName = ret.StrategyName
			vr.StatusCode = ret.StatusCode

			if err != nil {
				vr.Error = err
				errCh <- *vr
			} else {
				vr.Malicious = ret.Malicious
				retCh <- *vr
			}
		}(s)
	}

	for _, n := range strategies {
		select {
		case err := <-errCh:
			if err.Error != nil {
				cancel()
				return c.JSON(http.StatusOK, err)
			}
		case ret := <-retCh:
			if true == ret.Malicious {
				cancel()
				log.Infof("Return from [%d] %s", n, ret.StrategyName)
				return c.JSON(http.StatusOK, ret)
			}
		// Cancel is returned when either Timeout or Cancel occur
		case <-ctx.Done():
			<-errCh
			return ctx.Err()
		}
	}

	// No malicious links are found
	return c.JSON(http.StatusOK, &VerifyResponse{
		Links:     []string{},
		Malicious: false,
		Error:     nil,
	})
}
