package request

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/GeoIrb/itv/app"
)

//ClientRequest данные запроса клиента
type ClientRequest struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

//ClientResponse данные для ответа клиенту
type ClientResponse struct {
	ID      int         `json:"id"`
	Status  int         `json:"status"`
	Headers http.Header `json:"headers"`
	Length  int         `json:"length"`
}

//TaskWorker данные программы
//env - окружение, необходимое для логирования
//idMax - максимальный id задачи
//tasks - список всех задачь, можно было использовать sync.Map, но посчитал что не будет нормально работать с atomic для idMAX
type TaskWorker struct {
	env   app.Data
	mutex sync.Mutex
	idMax int
	tasks map[int]ClientRequest
}

func NewTaskWorker(env app.Data) TaskWorker {
	return TaskWorker{
		env:   env,
		idMax: 0,
		tasks: make(map[int]ClientRequest),
	}
}

func (w *TaskWorker) Add(task ClientRequest) int {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	idTask := w.idMax
	w.idMax++
	w.tasks[idTask] = task

	return idTask
}

func (w *TaskWorker) Delete(id int) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	delete(w.tasks, id)
}

func (w *TaskWorker) FetchTask(taskChan <-chan []byte, resultChan chan<- []byte) {
	for task := range taskChan {
		go func(task []byte) {

			request := ClientRequest{}
			json.Unmarshal(task, &request)

			w.mutex.Lock()

			idTask := w.Add(request)

			response, _ := request.Do(10 * time.Second)

			resClient := ClientResponse{
				ID:      idTask,
				Status:  response.StatusCode,
				Headers: response.Header,
			}

			resByte, _ := json.Marshal(resClient)

			resultChan <- resByte

		}(task)
	}
}
