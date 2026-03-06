package jetstream

import (
	"context"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

const (
	// Stream names
	StreamTasks     = "TASKS"
	StreamHeartbeat = "AGENT_HEARTBEAT"

	// Subject patterns
	SubjectTasks      = "tasks.*"
	SubjectTasksAsync = "tasks.async"

	SubjectHeartbeat       = "agent.keepalive.*"
	SubjectHeartbeatPrefix = "agent.keepalive."

	// Sync tasks subject (for direct request-response)
	SubjectTasksSync = "sync.tasks"

	// Consumer names
	ConsumerTaskProcessor    = "task-processor"
	ConsumerHeartbeatMonitor = "heartbeat-monitor"

	// Queue group names
	QueueTasksSync = "queue-tasks-sync"
)

// SetupStreams creates the necessary JetStream streams
func SetupStreams(ctx context.Context, js jetstream.JetStream) error {
	// Create TASKS stream
	_, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:        StreamTasks,
		Description: "Stream for task distribution to agents",
		Subjects:    []string{SubjectTasks},
		Retention:   jetstream.WorkQueuePolicy, // Messages removed after ACK
		Storage:     jetstream.FileStorage,     // Persistent storage
		MaxAge:      24 * time.Hour,            // Keep for 24 hours max
		Discard:     jetstream.DiscardOld,
	})
	if err != nil {
		return err
	}

	// Create AGENT_HEARTBEAT stream
	_, err = js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:        StreamHeartbeat,
		Description: "Stream for agent heartbeat messages",
		Subjects:    []string{SubjectHeartbeat},
		Retention:   jetstream.LimitsPolicy,  // Keep based on limits
		Storage:     jetstream.MemoryStorage, // Memory storage (faster)
		MaxAge:      1 * time.Hour,           // Keep for 1 hour
		MaxMsgs:     1000,                    // Keep max 1000 messages
		Discard:     jetstream.DiscardOld,
	})
	if err != nil {
		return err
	}

	return nil
}

// SetupConsumers creates the necessary JetStream consumers
func SetupConsumers(ctx context.Context, js jetstream.JetStream) error {
	// Create Tasks consumer for agents
	_, err := js.CreateOrUpdateConsumer(ctx, StreamTasks, jetstream.ConsumerConfig{
		Name:          ConsumerTaskProcessor,
		Description:   "Consumer for processing tasks by agents",
		Durable:       ConsumerTaskProcessor,
		DeliverPolicy: jetstream.DeliverAllPolicy,  // Process all pending messages
		AckPolicy:     jetstream.AckExplicitPolicy, // Require explicit ACK
		AckWait:       30 * time.Second,            // Wait 30s for ACK before redelivery
		MaxDeliver:    3,                           // Retry up to 3 times
		FilterSubject: SubjectTasksAsync,
	})
	if err != nil {
		return err
	}

	// Create Heartbeat consumer for server monitoring
	_, err = js.CreateOrUpdateConsumer(ctx, StreamHeartbeat, jetstream.ConsumerConfig{
		Name:          ConsumerHeartbeatMonitor,
		Description:   "Consumer for monitoring agent heartbeats",
		Durable:       ConsumerHeartbeatMonitor,
		DeliverPolicy: jetstream.DeliverNewPolicy, // Only new messages
		AckPolicy:     jetstream.AckAllPolicy,     // Auto-acknowledge
		FilterSubject: SubjectHeartbeat,
	})
	if err != nil {
		return err
	}

	return nil
}
