package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/GeoIrb/itv/app"
	"github.com/GeoIrb/itv/controllers"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const defaultPort int = 9995

func main() {
	conn := app.Init()
	defer conn.Stop()

	e := echo.New()

	e.Use(middleware.Recover())
	// e.Use(middleware.Logger())
	e.Use(middleware.Gzip())
	e.Use(middleware.Secure())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.OPTIONS, echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))
	e.Use(middleware.Static("react/build"))

	taskc := controllers.NewTaskController(conn)
	reqg := e.Group("/request")
	reqg.GET("", taskc.GetTasks)
	reqg.POST("", taskc.FetchTask)
	reqg.POST("/chan", taskc.FetchTaskChan)
	reqg.DELETE("/:id", taskc.DeleteTask)

	go taskc.Worker.HandlingChan(taskc.ReqChan, taskc.ResChan)

	go func() {
		if err := e.Start(getServerPort()); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	close(taskc.ReqChan)
	close(taskc.ResChan)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func getServerPort() (port string) {
	waCnf := app.Load(app.GetPath())

	serverPort := defaultPort

	if waCnf != nil {
		if port, isExist := waCnf["port"].(int); isExist {
			serverPort = port
		}
	}

	return fmt.Sprintf(":%d", serverPort)
}
