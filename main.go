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

	taskController := controllers.NewTaskController(conn)
	defer taskController.Kill()

	reqGroup := e.Group("/request")
	reqGroup.GET("", taskController.GetTasks)
	reqGroup.POST("", taskController.FetchTask)
	reqGroup.POST("/chan", taskController.FetchTaskChan)
	reqGroup.DELETE("/:id", taskController.DeleteTask)

	go taskController.Worker.HandlingChan(taskController.ReqChan, taskController.ResChan)

	go func() {
		if err := e.Start(getServerPort()); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func getServerPort() (port string) {
	waCnf := app.Load()

	serverPort := defaultPort

	if waCnf != nil {
		if port, isExist := waCnf["port"].(int); isExist {
			serverPort = port
		}
	}

	return fmt.Sprintf(":%d", serverPort)
}
