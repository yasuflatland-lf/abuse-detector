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
	StrategyName string   `json:"strategyName" xml:"strategyName"`
	Links        []string `json:"link" xml:"link"`
	Malicious    bool     `json:"malicious" xml:"malicious"`
	Error        error    `json:"error" xml:"error"`
}

// Handler
func Verify(c echo.Context) error {
	url := c.QueryParam("url")

	strategies := []verify.Verify{
		verify.NewTransparencyReportVerifyStrategy(),
		// verify.NewUrlScanVerifyStrategy(),
	}

	// Set up channels
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	errCh := make(chan VerifyResponse, len(strategies))
	retCh := make(chan VerifyResponse, len(strategies))

	for _, strategy := range strategies {

		// Run each verification concurrent.
		// Results are verified by the response order.
		// Return the result as soon as the site is confirmed,
		// including malicious links.
		go func() {
			vr := &VerifyResponse{
				Links:     []string{},
				Malicious: false,
				Error:     nil,
			}

			ret, err := strategy.Do(ctx, url)

			if err != nil {
				vr.Error = err
				vr.StrategyName = ret.StrategyName
				errCh <- *vr
			} else {
				vr.Links = ret.MaliciousLinks
				vr.Malicious = ret.Malicious
				vr.StrategyName = ret.StrategyName
				retCh <- *vr

				//if true == ret.Malicious {
				//	log.Info("Malicious link found. Interrupt processing...")
				//
				//	// Cancel all other process once
				//	// one of verification confirmed the URL is malicious.
				//	cancel()
				//}
			}
		}()
	}

	for _, n := range strategies {
		select {
		case err := <-errCh:
			return c.JSON(http.StatusOK, err)

		case ret := <-retCh:
			if true == ret.Malicious {
				return c.JSON(http.StatusOK, ret)
			}
			log.Infof("Return from [%d] %s", n, ret.StrategyName)
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
