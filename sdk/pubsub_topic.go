package sdk

import (
	"sync"
	"time"

	"github.com/gofiber/websocket/v2"
)

// Message represents a published message with server timestamp
type Message struct {
	ID      string      `json:"id"`
	Payload interface{} `json:"payload"`
	TS      time.Time   `json:"-"` // server timestamp, not serialized in message field
}

// Subscriber represents a client connection with buffered message queue
type Subscriber struct {
	Conn         *websocket.Conn
	ClientID     string
	Queue        chan Message
	QueueSize    int
	LastActive   time.Time
	CloseOnce    sync.Once
	CloseChannel chan struct{}
}

// Topic holds subscribers and implements ring buffer for message replay
type Topic struct {
	Name        string
	Subscribers map[string]*Subscriber // client_id -> Subscriber
	Messages    []Message              // ring buffer for last_n replay
	MaxMessages int
	Mu          sync.RWMutex // exported field
}

// TopicStats represents statistics for a single topic
type TopicStats struct {
	Messages    int `json:"messages"`
	Subscribers int `json:"subscribers"`
}

// WebSocket Protocol Structs

// WebSocketRequest represents incoming WebSocket messages from clients
type WebSocketRequest struct {
	Type      string   `json:"type"`
	Topic     string   `json:"topic,omitempty"`
	Message   *Message `json:"message,omitempty"`
	ClientID  string   `json:"client_id,omitempty"`
	LastN     int      `json:"last_n,omitempty"`
	RequestID string   `json:"request_id,omitempty"`
}

// WebSocketResponse represents outgoing WebSocket messages to clients
type WebSocketResponse struct {
	Type      string       `json:"type"`
	RequestID string       `json:"request_id,omitempty"`
	Topic     string       `json:"topic,omitempty"`
	Message   *Message     `json:"message,omitempty"`
	Status    string       `json:"status,omitempty"`
	Error     *ErrorDetail `json:"error,omitempty"`
	Timestamp string       `json:"ts,omitempty"`
	Msg       string       `json:"msg,omitempty"`
}

// ErrorDetail represents error information in responses
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// HTTP API Structs

// CreateTopicRequest represents a topic creation request
type CreateTopicRequest struct {
	Name string `json:"name"`
}

// CreateTopicResponse represents a topic creation response
type CreateTopicResponse struct {
	Status string `json:"status"`
	Topic  string `json:"topic"`
}

// DeleteTopicResponse represents a topic deletion response
type DeleteTopicResponse struct {
	Status string `json:"status"`
	Topic  string `json:"topic"`
}

// TopicInfo represents basic topic information
type TopicInfo struct {
	Name        string `json:"name"`
	Subscribers int    `json:"subscribers"`
}

// ListTopicsResponse represents the response from listing topics
type ListTopicsResponse struct {
	Topics []TopicInfo `json:"topics"`
}

// HealthResponse represents system health information
type HealthResponse struct {
	UptimeSeconds int `json:"uptime_sec"`
	Topics        int `json:"topics"`
	Subscribers   int `json:"subscribers"`
}

// StatsResponse represents system statistics
type StatsResponse struct {
	Topics map[string]TopicStats `json:"topics"`
}

// Error Response for HTTP APIs
type ErrorResponse struct {
	Error string `json:"error"`
}

// Constants for WebSocket message types
const (
	MessageTypeSubscribe   = "subscribe"
	MessageTypeUnsubscribe = "unsubscribe"
	MessageTypePublish     = "publish"
	MessageTypePing        = "ping"
	MessageTypeAck         = "ack"
	MessageTypeEvent       = "event"
	MessageTypeError       = "error"
	MessageTypePong        = "pong"
	MessageTypeInfo        = "info"
)

// Constants for error codes
const (
	ErrorCodeBadRequest    = "BAD_REQUEST"
	ErrorCodeTopicNotFound = "TOPIC_NOT_FOUND"
	ErrorCodeSlowConsumer  = "SLOW_CONSUMER"
	ErrorCodeUnauthorized  = "UNAUTHORIZED"
	ErrorCodeInternal      = "INTERNAL"
)

// Constants for HTTP status messages
const (
	StatusCreated  = "created"
	StatusDeleted  = "deleted"
	StatusOK       = "ok"
	StatusConflict = "conflict"
)
