package service

import (
	"nats-go/server/dto/model"
	"nats-go/server/dto/req"

	"github.com/bytedance/gopkg/util/logger"
)

type TaskServiceInterface interface {
	SendTask(task *req.Task) error
}

type TaskService struct {
}

func NewTaskService() TaskServiceInterface {
	return &TaskService{}
}

func (s *TaskService) SendTask(t *req.Task) error {
	task := &model.Task{
		Name:        t.Name,
		Description: t.Description,
	}

	logger.Infof("Send task: %+v", task)

	return nil
}
