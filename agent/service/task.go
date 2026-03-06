package service

import (
	"encoding/json"
	"log"
	"nats-go/server/dto/model"

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
	ts.sub, err = nc.Subscribe("tasks", func(msg *nats.Msg) {
		var task model.Task
		if err := json.Unmarshal(msg.Data, &task); err != nil {
			log.Printf("Error unmarshaling task: %v", err)
			return
		}

		log.Printf("Received task: ID=%s, Name=%s, Description=%s",
			task.ID, task.Name, task.Description)

		// Process the task here
		ts.processTask(&task)
	})
	if err != nil {
		log.Fatalf("Error subscribing to tasks: %v", err)
	}

	log.Println("Agent is listening for tasks on 'tasks' subject...")

	return ts
}

func (s *TaskService) Close() {
	if err := s.sub.Unsubscribe(); err != nil {
		log.Printf("Error unsubscribing: %v", err)
	}
}

func (s *TaskService) processTask(task *model.Task) {
	// Implement your task processing logic here
	log.Printf("Processing task %s: name: %s, description: %s", task.ID, task.Name, task.Description)
	// Add your business logic here
	log.Printf("Task %s completed successfully", task.ID)
}
