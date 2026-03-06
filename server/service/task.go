package service

import (
	"encoding/json"
	"nats-go/server/dto/model"
	"nats-go/server/dto/req"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

type TaskServiceInterface interface {
	SendTask(task *req.Task) error
}

type TaskService struct {
	nc *nats.Conn
}

func NewTaskService(nc *nats.Conn) (TaskServiceInterface, error) {
	return &TaskService{
		nc: nc,
	}, nil
}

func (s *TaskService) SendTask(t *req.Task) error {
	task := &model.Task{
		ID:          uuid.New().String(),
		Name:        t.Name,
		Description: t.Description,
	}

	logger.Infof("Send task: %+v", task)

	// Encode task to JSON
	data, err := json.Marshal(task)
	if err != nil {
		logger.Errorf("Failed to marshal task: %v", err)
		return err
	}

	// Publish to NATS
	if err := s.nc.Publish("tasks", data); err != nil {
		logger.Errorf("Failed to publish task to NATS: %v", err)
		return err
	}

	logger.Infof("Task published successfully to NATS: %s", task.ID)
	return nil
}
