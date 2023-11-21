package api

import (
	"encoding/json"

	_ "fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tb/task-logger/common-packages/system"
	cache "github.com/tb/task-logger/taskapp/cache"
	"github.com/tb/task-logger/taskapp/models"
	"github.com/tb/task-logger/taskapp/services"
	_ "github.com/tb/task-logger/taskapp/validations"
	"github.com/zenazn/goji/web"
	"net/http"
	"strconv"
)

type Controller struct {
	TaskController
}



type TaskController interface {

	AddTask(c web.C, w http.ResponseWriter, r *http.Request, logger *log.Entry) ([]byte, error)
}

func (controller *Controller) AddTask( c web.C, w http.ResponseWriter, r *http.Request, logger *log.Entry) ([]byte, error) {
	decoder := json.NewDecoder(r.Body)
	var addTaskInput models.AddTaskInput
	err := decoder.Decode(&addTaskInput)
	if err != nil {
		logger.Error(err)
		return nil, system.InvalidPayloadError
	}
	logger.Info(addTaskInput)
	response, err := services.AddTask(logger, addTaskInput, c.Env["emailId"].(string))
	return response, err
}


