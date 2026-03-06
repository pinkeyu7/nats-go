package service

import (
	"encoding/json"
	"nats-go/pkg/topic"
	"nats-go/server/dto/model"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/nats-io/nats.go"
)

type TaskServiceInterface interface {
	Close()
}

type TaskService struct {
	nc  *nats.Conn
	sub *nats.Subscription
}

func NewTaskService(nc *nats.Conn) TaskServiceInterface {
	ts := &TaskService{
		nc: nc,
	}

	var err error
	// Subscribe to tasks subject
	ts.sub, err = nc.QueueSubscribe(topic.TopicTasks, topic.QueueTasks, func(msg *nats.Msg) {
		var task model.Task
		if err := json.Unmarshal(msg.Data, &task); err != nil {
			logger.Infof("Error unmarshaling task: %v", err)
			return
		}

		logger.Infof("Received task: ID=%s, Name=%s, Description=%s",
			task.ID, task.Name, task.Description)

		// Process the task here
		task = ts.processTask(&task)

		// Encode task to JSON
		data, err := json.Marshal(task)
		if err != nil {
			logger.Errorf("Failed to marshal task: %v", err)
			return
		}

		// Send the response back
		err = msg.Respond(data)
		if err != nil {
			logger.Infof("Error responding to task: %v", err)
		}
	})
	if err != nil {
		logger.Fatalf("Error subscribing to tasks: %v", err)
	}

	logger.Infof("Agent is listening for tasks on 'tasks' subject...")

	return ts
}

func (s *TaskService) Close() {
	if err := s.sub.Unsubscribe(); err != nil {
		logger.Infof("Error unsubscribing: %v", err)
	}
}

func (s *TaskService) processTask(task *model.Task) model.Task {
	// Implement your task processing logic here
	logger.Infof("Processing task %s: name: %s, description: %s", task.ID, task.Name, task.Description)
	// Add your business logic here
	logger.Infof("Task %s completed successfully", task.ID)

	task.Result = "Task processed successfully"

	return *task
}
