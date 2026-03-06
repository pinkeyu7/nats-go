package service

import (
	"context"
	"encoding/json"
	"nats-go/pkg/jetstream"
	"nats-go/server/dto/model"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/nats-io/nats.go"
	js "github.com/nats-io/nats.go/jetstream"
)

type TaskServiceInterface interface {
	Close()
}

type TaskService struct {
	nc       *nats.Conn
	js       js.JetStream
	consumer js.Consumer
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewTaskService(nc *nats.Conn, jsContext js.JetStream) TaskServiceInterface {
	ts := &TaskService{
		nc: nc,
		js: jsContext,
	}

	ts.ctx, ts.cancel = context.WithCancel(context.Background())

	var err error
	// Get the consumer
	ts.consumer, err = jsContext.Consumer(ts.ctx, jetstream.StreamTasks, jetstream.ConsumerTaskProcessor)
	if err != nil {
		logger.Fatalf("Error getting consumer: %v", err)
	}

	// Start consuming messages
	go ts.consumeAsyncMessages()
	go ts.consumeSyncMessages()

	logger.Infof("Agent is listening for tasks via JetStream...")

	return ts
}

func (s *TaskService) consumeAsyncMessages() {
	logger.Infof("Starting to consume asynchronous tasks from JetStream...")

	// Consume messages from the stream
	cons, err := s.consumer.Consume(func(msg js.Msg) {
		var task model.Task
		if err := json.Unmarshal(msg.Data(), &task); err != nil {
			logger.Errorf("Error unmarshaling task: %v", err)
			// Acknowledge with NAK to indicate processing failure
			if err := msg.Nak(); err != nil {
				logger.Errorf("Error sending NAK: %v", err)
			}
			return
		}

		logger.Infof("Received task: ID=%s, Name=%s, Description=%s",
			task.ID, task.Name, task.Description)

		// Process the task
		task = s.processTask(&task)

		// Acknowledge successful processing
		if err := msg.Ack(); err != nil {
			logger.Errorf("Error acknowledging message: %v", err)
			return
		}

		logger.Infof("Task %s acknowledged successfully", task.ID)
	})

	if err != nil {
		logger.Fatalf("Error consuming messages: %v", err)
	}

	logger.Info("Agent is consuming asynchronous tasks from JetStream...")

	// Wait for context cancellation
	<-s.ctx.Done()
	logger.Infof("Agent is stopping consumption of asynchronous tasks")
	cons.Stop()
}

func (s *TaskService) consumeSyncMessages() {
	logger.Infof("Starting to consume synchronous tasks from NATS...")

	// Subscribe to tasks subject
	sub, err := s.nc.QueueSubscribe(jetstream.SubjectTasksSync, jetstream.QueueTasksSync, func(msg *nats.Msg) {
		logger.Infof("Received message: %s", msg.Data)
		var task model.Task
		if err := json.Unmarshal(msg.Data, &task); err != nil {
			logger.Infof("Error unmarshaling task: %v", err)
			return
		}

		// Process the task here
		t := s.processTask(&task)

		// Encode task to JSON
		data, err := json.Marshal(t)
		if err != nil {
			logger.Errorf("Failed to marshal task: %v", err)
			return
		}

		// Send the response back
		err = msg.Respond(data)
		if err != nil {
			logger.Infof("Error responding to task: %v", err)
			return
		}
	})
	if err != nil {
		logger.Fatalf("Error subscribing to tasks subject: %v", err)
	}

	logger.Info("Agent is subscribed to synchronous tasks subject...")

	// Wait for context cancellation
	<-s.ctx.Done()
	logger.Infof("Agent is unsubscribed from tasks subject")
	_ = sub.Unsubscribe()
}

func (s *TaskService) Close() {
	logger.Info("Closing TaskService...")
	s.cancel()
}

func (s *TaskService) processTask(task *model.Task) model.Task {
	// Implement your task processing logic here
	logger.Infof("Processing task %s: name: %s, description: %s", task.ID, task.Name, task.Description)
	// Add your business logic here
	logger.Infof("Task %s completed successfully", task.ID)

	task.Result = "Task processed successfully"

	return *task
}
