package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/Pineapple217/Corap-web/pkg/database"
	"github.com/Pineapple217/Corap-web/pkg/handler"
	"github.com/Pineapple217/Corap-web/pkg/server"
	"github.com/Pineapple217/Corap-web/pkg/static"
)

const banner = `
 ██████╗ ██████╗ ██████╗  █████╗ ██████╗ 
██╔════╝██╔═══██╗██╔══██╗██╔══██╗██╔══██╗
██║     ██║   ██║██████╔╝███████║██████╔╝
██║     ██║   ██║██╔══██╗██╔══██║██╔═══╝ 
╚██████╗╚██████╔╝██║  ██║██║  ██║██║     
 ╚═════╝ ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝      v2.0.0
https://github.com/Pineapple217/Corap-web
-----------------------------------------------------------------------------`

func main() {
	slog.SetDefault(slog.New(slog.Default().Handler()))
	fmt.Println(banner)
	os.Stdout.Sync()

	rr := static.HashPublicFS()

	db := database.NewDatabase()
	h := handler.NewHandler(db)

	server := server.NewServer()
	server.RegisterRoutes(h)
	server.ApplyMiddleware(rr)

	server.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	slog.Info("Received an interrupt signal, exiting...")

	server.Stop()
}
