package service

import (
	"context"
	"encoding/json"
	"nats-go/pkg/jetstream"
	"nats-go/server/dto/model"
	"nats-go/server/dto/req"
	"time"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	js "github.com/nats-io/nats.go/jetstream"
)

type TaskServiceInterface interface {
	SendAsyncTask(task *req.Task) (*model.Task, error)
	SendSyncTask(task *req.Task) (*model.Task, error)
}

type TaskService struct {
	nc *nats.Conn
	js js.JetStream
}

func NewTaskService(nc *nats.Conn, jsContext js.JetStream) (TaskServiceInterface, error) {
	return &TaskService{
		nc: nc,
		js: jsContext,
	}, nil
}

func (s *TaskService) SendAsyncTask(t *req.Task) (*model.Task, error) {
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

	// Publish to JetStream
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ack, err := s.js.Publish(ctx, jetstream.SubjectTasksAsync, data)
	if err != nil {
		logger.Errorf("Failed to publish task to JetStream: %v", err)
		return nil, err
	}

	logger.Infof("Task published successfully to JetStream: %s (sequence: %d)", task.ID, ack.Sequence)
	return task, nil
}

func (s *TaskService) SendSyncTask(t *req.Task) (*model.Task, error) {
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
	msg, err := s.nc.Request(jetstream.SubjectTasksSync, data, 2*time.Second)
	if err != nil {
		logger.Errorf("Failed to publish task to NATS: %v", err)
		return nil, err
	}

	// Unmarshal response
	var res model.Task
	if err := json.Unmarshal(msg.Data, &res); err != nil {
		logger.Errorf("Failed to unmarshal response: %v", err)
		return nil, err
	}

	logger.Infof("Task published successfully to NATS: %s", res.ID)

	return &res, nil
}
