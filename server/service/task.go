package service

import (
	"encoding/json"
	"nats-go/pkg/topic"
	"nats-go/server/dto/model"
	"nats-go/server/dto/req"
	"time"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

type TaskServiceInterface interface {
	SendTask(task *req.Task) (*model.Task, error)
}

type TaskService struct {
	nc *nats.Conn
}

func NewTaskService(nc *nats.Conn) (TaskServiceInterface, error) {
	return &TaskService{
		nc: nc,
	}, nil
}

func (s *TaskService) SendTask(t *req.Task) (*model.Task, error) {
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
		return nil, err
	}

	// Request to NATS
	msg, err := s.nc.Request(topic.TopicTasks, data, 2*time.Second)
	if err != nil {
		logger.Errorf("Failed to publish task to NATS: %v", err)
		return nil, err
	}

	// Process response from NATS
	res := model.Task{}
	err = json.Unmarshal(msg.Data, &res)
	if err != nil {
		logger.Errorf("Failed to unmarshal response from NATS: %v", err)
		return nil, err
	}

	logger.Infof("Task published successfully to NATS: %s", task.ID)
	return &res, nil
}
