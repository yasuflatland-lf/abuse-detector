package detector

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	e.GET("/verify", verify)

	// Start server
	e.Logger.Fatal(e.Start(":3000"))
}

// Handler
func verify(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

