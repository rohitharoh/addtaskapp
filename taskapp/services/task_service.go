package services

import (
	"encoding/json"
	"fmt"
	"github.com/pborman/uuid"
	"github.com/prometheus/common/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"github.com/tb/task-logger/common-packages/system"
	cache "github.com/tb/task-logger/taskapp/cache"
	"github.com/tb/task-logger/taskapp/models"
	"github.com/tb/task-logger/taskapp/validations"

	"gopkg.in/mgo.v2/bson"
	"time"
)
type TaskService interface {
	AddTask(logger *logrus.Entry, createTaskInput models.AddTaskInput, emailId string) ([]byte, error)
}

func AddTask(logger *logrus.Entry, createTaskInput models.AddTaskInput, emailId string) ([]byte, error) {

	isValid := Validationspackage.ValidateEmail(emailId)
	if !isValid {
		return nil, system.InvalidEmailErr
	}

	err := Validationspackage.ValidateAddTaskInput(logger, createTaskInput)
	if err != nil {
		return nil, err
	}

	taskObj := models.Task{
		Id:          uuid.New(),
		Title:       createTaskInput.Title,
		ScheduledOn: createTaskInput.ScheduledOn,
		Description: createTaskInput.Description,
		EmailId:     emailId,
		Status:      system.TASK_STATUS_PENDING,
		CreatedOn:   time.Now(),
		ModifiedOn:  time.Now(),
	}

	fmt.Println("_id", taskObj.Id)

//cache.NewRedisCache("127.0.0.1", 0, system.REDIS_DEFAULT_EXPIRATION_TIME).Flush("")
	client := cache.NewRedisCache("127.0.0.1", 0, system.REDIS_DEFAULT_EXPIRATION_TIME)

	client.Set(system.TASKS_COLLECTION+ ":" + taskObj.Id, &taskObj)

	collectionName := system.TASKS_COLLECTION
	databaseName := system.GetDatabaseName(collectionName)
	sessionDb := system.TbAppContext.MongoDBSessionMap[databaseName].Clone()
	defer sessionDb.Close()
	collection := sessionDb.DB(databaseName).C(collectionName)
	err = collection.Insert(&taskObj)
	if err != nil {
		return nil, err
	}
	response := make(map[string]interface{})
	response["message"] = "Task created successfully"
	finalResponse, _ := json.Marshal(response)
	return finalResponse, nil

}

