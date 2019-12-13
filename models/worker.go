package models

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/GeoIrb/itv/app"
)

//ClientResponse данные для ответа клиенту
type ClientResponse struct {
	ID      int         `json:"id"`
	Status  int         `json:"status"`
	Headers http.Header `json:"headers"`
	Length  int         `json:"length"`
}

//RequestWorker данные программы
//env - окружение, необходимое для логирования
//idMax - максимальный id задачи
//requests - список всех задачь, можно было использовать sync.Map, но посчитал что не будет нормально работать с atomic для idMAX
type RequestWorker struct {
	env      app.Data
	mutex    sync.Mutex
	idMax    int
	requests map[int]ClientRequest
}

//NewRequestWorker создание TaskWorker
func NewRequestWorker(env app.Data) RequestWorker {
	return RequestWorker{
		env:      env,
		idMax:    0,
		requests: make(map[int]ClientRequest),
	}
}

//GetRequests получить список задач
func (w *RequestWorker) GetRequests() map[int]ClientRequest {
	return w.requests
}

//Add добавление запроса в список
func (w *RequestWorker) Add(task ClientRequest) int {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	idTask := w.idMax
	w.idMax++
	w.requests[idTask] = task

	return idTask
}

//Delete удаление запроса из списка
func (w *RequestWorker) Delete(id int) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	delete(w.requests, id)
}

//Handling обработка запроса
func (w *RequestWorker) Handling(request ClientRequest) (ClientResponse, error) {
	idTask := w.Add(request)

	response, err := request.Do(10 * time.Second)
	if err != nil {
		return ClientResponse{}, fmt.Errorf("Error request do: %v", err)
	}

	resClient := ClientResponse{
		ID:      idTask,
		Status:  response.StatusCode,
		Headers: response.Header,
	}

	return resClient, nil
}

//HandlingChan получение задач из reqChan и отправка результата в resultChan
func (w *RequestWorker) HandlingChan(reqChan <-chan ClientRequest, resultChan chan<- ClientResponse) {
	for request := range reqChan {
		go func(request ClientRequest) {
			resByte, err := w.Handling(request)
			if err != nil {
				w.env.Err("FetchTask %v", err)
				return
			}

			resultChan <- resByte
		}(request)
	}
}
