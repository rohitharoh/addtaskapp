package routes

import (
	"github.com/tb/task-logger/common-packages/system"
	"github.com/tb/task-logger/taskapp/controllers/api"
	"github.com/zenazn/goji"
)

func PrepareRoutes(application *system.Application) {

	//task logger

	// goji.Get("/email", application.Route(&api.Controller{}, "Email", true, []string{"admin"}))
	goji.Post("/application/service/task/add", application.Route(&api.Controller{}, "AddTask", false, nil))

}
