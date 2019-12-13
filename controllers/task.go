package controllers

import (
	"net/http"
	"strconv"

	"github.com/GeoIrb/itv/app"
	"github.com/GeoIrb/itv/models"
	"github.com/labstack/echo"
)

//TaskController контроллер для модуля
type TaskController struct {
	Worker  models.RequestWorker
	ReqChan chan models.ClientRequest
	ResChan chan models.ClientResponse
}

func NewTaskController(env app.Data) TaskController {
	return TaskController{
		Worker:  models.NewRequestWorker(env),
		ReqChan: make(chan models.ClientRequest, 1),
		ResChan: make(chan models.ClientResponse, 1),
	}
}

//FetchTask обработка запроса
func (tc *TaskController) FetchTask(context echo.Context) error {
	var req models.ClientRequest
	if err := context.Bind(&req); err != nil {
		return context.JSON(http.StatusBadRequest, models.Error{Message: err.Error()})
	}

	result, err := tc.Worker.Handling(req)
	if err != nil {
		context.JSON(http.StatusNotFound, models.Error{Message: err.Error()})
	}

	return context.JSON(http.StatusOK, result)
}

//FetchTaskChan обработка запроса через канал
func (tc *TaskController) FetchTaskChan(context echo.Context) error {
	var req models.ClientRequest
	if err := context.Bind(&req); err != nil {
		return context.JSON(http.StatusBadRequest, models.Error{Message: err.Error()})
	}

	tc.ReqChan <- req
	result := <-tc.ResChan

	return context.JSON(http.StatusOK, result)
}

//GetTasks обработка запроса
func (tc *TaskController) GetTasks(context echo.Context) error {
	return context.JSON(http.StatusOK, tc.Worker.GetRequests())
}

//DeleteTask удаление запроса
func (tc *TaskController) DeleteTask(context echo.Context) error {
	idParam := context.Param(":id")
	if id, err := strconv.Atoi(idParam); err != nil {
		return context.JSON(http.StatusBadRequest, models.Error{Message: err.Error()})
	} else {
		tc.Worker.Delete(id)
		return context.NoContent(http.StatusOK)
	}
}
