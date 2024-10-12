package server

import (
	"log/slog"
	"time"

	"github.com/Pineapple217/Corap-web/pkg/handler"

	"github.com/labstack/echo/v4"
	echoMw "github.com/labstack/echo/v4/middleware"
)

func (s *Server) ApplyMiddleware(reRoutes map[string]string) {
	slog.Info("Applying middlewares")
	s.e.Pre(echoMw.Rewrite(reRoutes))
	s.e.Use(echoMw.RequestLoggerWithConfig(echoMw.RequestLoggerConfig{
		LogStatus:  true,
		LogURI:     true,
		LogMethod:  true,
		LogLatency: true,
		LogValuesFunc: func(c echo.Context, v echoMw.RequestLoggerValues) error {
			slog.Info("request",
				"method", v.Method,
				"status", v.Status,
				"latency", v.Latency,
				"path", v.URI,
			)
			return nil

		},
	}))

	s.e.Use(echoMw.RateLimiterWithConfig(echoMw.RateLimiterConfig{
		Store: echoMw.NewRateLimiterMemoryStoreWithConfig(
			echoMw.RateLimiterMemoryStoreConfig{Rate: 10, Burst: 30, ExpiresIn: 3 * time.Minute},
		),
	}))

	s.e.Use(echoMw.GzipWithConfig(echoMw.GzipConfig{
		Level: 5,
	}))

	echo.NotFoundHandler = handler.NotFound
}
