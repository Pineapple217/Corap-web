package server

import (
	"log/slog"

	"github.com/Pineapple217/Corap-web/pkg/handler"
	"github.com/Pineapple217/Corap-web/pkg/static"

	"github.com/labstack/echo/v4"
)

func (server *Server) RegisterRoutes(hdlr *handler.Handler) {
	slog.Info("Registering routes")
	e := server.e

	s := e.Group("/static")

	s.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Add("Cache-Control", "public, max-age=31536000, immutable")
			return next(c)
		}
	})
	s.StaticFS("/", echo.MustSubFS(static.PublicFS, "public"))

	// e.GET("/index.xml", hdlr.RSSFeed)
	e.GET("/devices", hdlr.DevicesHome)
	e.GET("/devices/table", hdlr.DevicesTable)
	e.GET("/devices/table_ana", hdlr.DevicesAnalysisTable)
	e.GET("/devices/:deveui/history", hdlr.DeviceHistory)
	e.GET("/devices/:deveui", hdlr.DeviceHome)

	e.GET("/stats", hdlr.StatsHome)

	// e.GET("robot.txt", hdlr.RobotTxt)
	// e.GET("/site.webmanifest", hdlr.Manifest)

	//TODO better caching with http headers

	// e.GET("/", hdlr.AuthForm)

	// e.Static("/m", config.UploadDir)
}
