package redis

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// Publisher wraps a go-redis client and provides typed event publishing.
type Publisher struct {
	client *redis.Client
}

// NewPublisher creates a Publisher and verifies connectivity.
func NewPublisher(addr, password string, db int) (*Publisher, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Publisher{client: client}, nil
}

// Close releases the underlying Redis connection.
func (p *Publisher) Close() error {
	return p.client.Close()
}

// AssignmentCreatedEvent is published when a new assignment is created.
type AssignmentCreatedEvent struct {
	Event        string `json:"event"`
	AssignmentID string `json:"assignment_id"`
	ClassID      string `json:"class_id"`
	Title        string `json:"title"`
	Deadline     string `json:"deadline"`
	TeacherID    string `json:"teacher_id"`
	Timestamp    string `json:"timestamp"`
}

// AssignmentSubmittedEvent is published when a student submits an assignment.
type AssignmentSubmittedEvent struct {
	Event        string `json:"event"`
	SubmissionID string `json:"submission_id"`
	AssignmentID string `json:"assignment_id"`
	StudentID    string `json:"student_id"`
	ClassID      string `json:"class_id"`
	Timestamp    string `json:"timestamp"`
}

// AssignmentGradedEvent is published when a teacher grades a submission.
type AssignmentGradedEvent struct {
	Event        string  `json:"event"`
	SubmissionID string  `json:"submission_id"`
	AssignmentID string  `json:"assignment_id"`
	StudentID    string  `json:"student_id"`
	Grade        float64 `json:"grade"`
	Timestamp    string  `json:"timestamp"`
}

const (
	channelAssignmentCreated   = "lms:assignment:created"
	channelAssignmentSubmitted = "lms:assignment:submitted"
	channelAssignmentGraded    = "lms:assignment:graded"
)

// PublishAssignmentCreated sends an ASSIGNMENT_CREATED event to Redis Pub/Sub.
func (p *Publisher) PublishAssignmentCreated(ctx context.Context, evt AssignmentCreatedEvent) {
	p.publish(ctx, channelAssignmentCreated, evt)
}

// PublishAssignmentSubmitted sends an ASSIGNMENT_SUBMITTED event to Redis Pub/Sub.
func (p *Publisher) PublishAssignmentSubmitted(ctx context.Context, evt AssignmentSubmittedEvent) {
	p.publish(ctx, channelAssignmentSubmitted, evt)
}

// PublishAssignmentGraded sends an ASSIGNMENT_GRADED event to Redis Pub/Sub.
func (p *Publisher) PublishAssignmentGraded(ctx context.Context, evt AssignmentGradedEvent) {
	p.publish(ctx, channelAssignmentGraded, evt)
}

func (p *Publisher) publish(ctx context.Context, channel string, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[redis] failed to marshal event for channel %s: %v", channel, err)
		return
	}

	if err := p.client.Publish(ctx, channel, data).Err(); err != nil {
		log.Printf("[redis] failed to publish to channel %s: %v", channel, err)
	}
}
