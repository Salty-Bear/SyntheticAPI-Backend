package pubsub

import (
	"context"
	"sync"
	"time"

	"github.com/Aryaman/pub-sub/sdk"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// ServiceImpl implements the PubSub interface with thread-safe operations
type ServiceImpl struct {
	Topics      map[string]*sdk.Topic
	mu          sync.RWMutex
	Uptime      time.Time
	MaxQueue    int // per-subscriber queue size
	MaxMessages int // per-topic ring buffer size
}

// NewService creates a new PubSub service instance with config
func NewService(maxQueue, maxMessages int) *ServiceImpl {
	return &ServiceImpl{
		Topics:      make(map[string]*sdk.Topic),
		Uptime:      time.Now(),
		MaxQueue:    maxQueue,
		MaxMessages: maxMessages,
	}
}

// WebSocket message types for the writer goroutine
type wsMessage struct {
	Type string
	Data interface{}
}

// HandleWebSocket processes WebSocket connections and messages
func (s *ServiceImpl) HandleWebSocket(ctx context.Context, c *websocket.Conn) {
	defer c.Close()

	var clientID string
	var currentTopic *sdk.Topic
	var currentSub *sdk.Subscriber

	// Create a channel for serializing all WebSocket writes
	writeChannel := make(chan wsMessage, 100)
	writerDone := make(chan struct{})

	// Start WebSocket writer goroutine - this is the ONLY goroutine that writes to the connection
	go func() {
		defer close(writerDone)
		for msg := range writeChannel {
			if err := c.WriteJSON(msg.Data); err != nil {
				// Connection closed or error - stop processing
				return
			}
		}
	}()

	// Helper function to safely send messages through the writer channel
	sendMessage := func(msgType string, data interface{}) {
		select {
		case writeChannel <- wsMessage{Type: msgType, Data: data}:
		default:
			// Write channel full - connection is too slow, close it
			close(writeChannel)
		}
	}

	for {
		// Parse incoming message using SDK struct
		var req sdk.WebSocketRequest
		if err := c.ReadJSON(&req); err != nil {
			// Connection closed or invalid JSON
			break
		}

		// Handle different message types according to protocol specification
		switch req.Type {
		case sdk.MessageTypeSubscribe:
			if req.Topic == "" || req.ClientID == "" {
				sendMessage("error", sdk.WebSocketResponse{
					Type:      sdk.MessageTypeError,
					RequestID: req.RequestID,
					Error: &sdk.ErrorDetail{
						Code:    sdk.ErrorCodeBadRequest,
						Message: "topic and client_id required",
					},
					Timestamp: time.Now().UTC().Format(time.RFC3339),
				})
				continue
			}

			clientID = req.ClientID

			// Check if topic exists
			s.mu.RLock()
			topic, ok := s.Topics[req.Topic]
			s.mu.RUnlock()

			if !ok {
				sendMessage("error", sdk.WebSocketResponse{
					Type:      sdk.MessageTypeError,
					RequestID: req.RequestID,
					Error: &sdk.ErrorDetail{
						Code:    sdk.ErrorCodeTopicNotFound,
						Message: "topic not found",
					},
					Timestamp: time.Now().UTC().Format(time.RFC3339),
				})
				continue
			}

			// Create subscriber with bounded queue
			sub := &sdk.Subscriber{
				Conn:         c,
				ClientID:     clientID,
				Queue:        make(chan sdk.Message, s.MaxQueue),
				QueueSize:    s.MaxQueue,
				LastActive:   time.Now(),
				CloseChannel: make(chan struct{}),
			}

			// Add subscriber to topic and handle replay if requested
			topic.Mu.Lock()
			topic.Subscribers[clientID] = sub

			// Replay last_n messages from ring buffer
			if req.LastN > 0 && len(topic.Messages) > 0 {
				n := req.LastN
				if n > len(topic.Messages) {
					n = len(topic.Messages)
				}
				// Send the last n messages
				start := len(topic.Messages) - n
				for _, msg := range topic.Messages[start:] {
					select {
					case sub.Queue <- msg:
					default:
						// Queue full during replay - disconnect slow consumer
						close(sub.CloseChannel)
						delete(topic.Subscribers, clientID)
						sendMessage("error", sdk.WebSocketResponse{
							Type:      sdk.MessageTypeError,
							RequestID: req.RequestID,
							Error: &sdk.ErrorDetail{
								Code:    sdk.ErrorCodeSlowConsumer,
								Message: "subscriber queue overflow during replay",
							},
							Timestamp: time.Now().UTC().Format(time.RFC3339),
						})
						topic.Mu.Unlock()
						goto nextMessage
					}
				}
			}
			topic.Mu.Unlock()

			currentTopic = topic
			currentSub = sub

			// Start message delivery goroutine - it will use the same writeChannel
			go s.subscriberWriter(sub, req.Topic, sendMessage)

			sendMessage("ack", sdk.WebSocketResponse{
				Type:      sdk.MessageTypeAck,
				RequestID: req.RequestID,
				Topic:     req.Topic,
				Status:    sdk.StatusOK,
				Timestamp: time.Now().UTC().Format(time.RFC3339),
			})

		case sdk.MessageTypeUnsubscribe:
			if req.Topic == "" || req.ClientID == "" {
				sendMessage("error", sdk.WebSocketResponse{
					Type:      sdk.MessageTypeError,
					RequestID: req.RequestID,
					Error: &sdk.ErrorDetail{
						Code:    sdk.ErrorCodeBadRequest,
						Message: "topic and client_id required",
					},
					Timestamp: time.Now().UTC().Format(time.RFC3339),
				})
				continue
			}

			s.mu.RLock()
			topic, ok := s.Topics[req.Topic]
			s.mu.RUnlock()

			if !ok {
				sendMessage("error", sdk.WebSocketResponse{
					Type:      sdk.MessageTypeError,
					RequestID: req.RequestID,
					Error: &sdk.ErrorDetail{
						Code:    sdk.ErrorCodeTopicNotFound,
						Message: "topic not found",
					},
					Timestamp: time.Now().UTC().Format(time.RFC3339),
				})
				continue
			}

			// Remove subscriber from topic
			topic.Mu.Lock()
			if sub, exists := topic.Subscribers[req.ClientID]; exists {
				close(sub.CloseChannel)
				delete(topic.Subscribers, req.ClientID)
			}
			topic.Mu.Unlock()

			sendMessage("ack", sdk.WebSocketResponse{
				Type:      sdk.MessageTypeAck,
				RequestID: req.RequestID,
				Topic:     req.Topic,
				Status:    sdk.StatusOK,
				Timestamp: time.Now().UTC().Format(time.RFC3339),
			})

		case sdk.MessageTypePublish:
			if req.Topic == "" || req.Message == nil {
				sendMessage("error", sdk.WebSocketResponse{
					Type:      sdk.MessageTypeError,
					RequestID: req.RequestID,
					Error: &sdk.ErrorDetail{
						Code:    sdk.ErrorCodeBadRequest,
						Message: "topic and message required",
					},
					Timestamp: time.Now().UTC().Format(time.RFC3339),
				})
				continue
			}

			s.mu.RLock()
			topic, ok := s.Topics[req.Topic]
			s.mu.RUnlock()

			if !ok {
				sendMessage("error", sdk.WebSocketResponse{
					Type:      sdk.MessageTypeError,
					RequestID: req.RequestID,
					Error: &sdk.ErrorDetail{
						Code:    sdk.ErrorCodeTopicNotFound,
						Message: "topic not found",
					},
					Timestamp: time.Now().UTC().Format(time.RFC3339),
				})
				continue
			}

			// Add server timestamp
			req.Message.TS = time.Now().UTC()

			topic.Mu.Lock()

			// Add to ring buffer for replay functionality
			if len(topic.Messages) >= topic.MaxMessages {
				// Remove oldest message (ring buffer behavior)
				topic.Messages = topic.Messages[1:]
			}
			topic.Messages = append(topic.Messages, *req.Message)

			// Fan-out message to all subscribers
			slowConsumers := make([]string, 0)
			for clientID, sub := range topic.Subscribers {
				select {
				case sub.Queue <- *req.Message:
					// Message delivered successfully
				default:
					// Queue full - mark as slow consumer for removal
					slowConsumers = append(slowConsumers, clientID)
				}
			}

			// Remove slow consumers (backpressure policy: disconnect on overflow)
			for _, clientID := range slowConsumers {
				if sub, exists := topic.Subscribers[clientID]; exists {
					close(sub.CloseChannel)
					delete(topic.Subscribers, clientID)
					// Note: We can't send error to slow consumer's connection here
					// because we don't have access to their writeChannel
				}
			}

			topic.Mu.Unlock()
			sendMessage("ack", sdk.WebSocketResponse{
				Type:      sdk.MessageTypeAck,
				RequestID: req.RequestID,
				Topic:     req.Topic,
				Status:    sdk.StatusOK,
				Timestamp: time.Now().UTC().Format(time.RFC3339),
			})

		case sdk.MessageTypePing:
			sendMessage("pong", sdk.WebSocketResponse{
				Type:      sdk.MessageTypePong,
				RequestID: req.RequestID,
				Timestamp: time.Now().UTC().Format(time.RFC3339),
			})

		default:
			sendMessage("error", sdk.WebSocketResponse{
				Type:      sdk.MessageTypeError,
				RequestID: req.RequestID,
				Error: &sdk.ErrorDetail{
					Code:    sdk.ErrorCodeBadRequest,
					Message: "unknown message type",
				},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
			})
		}
	nextMessage:
	}

	// Cleanup on connection close
	if currentTopic != nil && currentSub != nil {
		currentTopic.Mu.Lock()
		delete(currentTopic.Subscribers, currentSub.ClientID)
		currentTopic.Mu.Unlock()
	}

	// Close write channel and wait for writer to finish
	close(writeChannel)
	<-writerDone
}

// CreateTopic creates a new topic via REST API
func (s *ServiceImpl) CreateTopic(ctx context.Context, c *fiber.Ctx) error {
	var req sdk.CreateTopicRequest

	if err := c.BodyParser(&req); err != nil || req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(sdk.ErrorResponse{
			Error: "invalid request - name required",
		})
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if topic already exists
	if _, exists := s.Topics[req.Name]; exists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"status": sdk.StatusConflict,
			"topic":  req.Name,
		})
	}

	// Create new topic
	s.Topics[req.Name] = &sdk.Topic{
		Name:        req.Name,
		Subscribers: make(map[string]*sdk.Subscriber),
		Messages:    make([]sdk.Message, 0, s.MaxMessages),
		MaxMessages: s.MaxMessages,
	}

	return c.Status(fiber.StatusCreated).JSON(sdk.CreateTopicResponse{
		Status: sdk.StatusCreated,
		Topic:  req.Name,
	})
}

// DeleteTopic removes a topic and disconnects all subscribers
func (s *ServiceImpl) DeleteTopic(ctx context.Context, c *fiber.Ctx) error {
	name := c.Params("name")
	if name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(sdk.ErrorResponse{
			Error: "missing topic name",
		})
	}

	s.mu.Lock()
	topic, exists := s.Topics[name]
	if !exists {
		s.mu.Unlock()
		return c.Status(fiber.StatusNotFound).JSON(sdk.ErrorResponse{
			Error: "topic not found",
		})
	}

	// Close all subscriber channels - their writer goroutines will handle cleanup
	topic.Mu.Lock()
	for _, sub := range topic.Subscribers {
		close(sub.CloseChannel)
	}
	topic.Subscribers = make(map[string]*sdk.Subscriber)
	topic.Mu.Unlock()

	// Remove topic
	delete(s.Topics, name)
	s.mu.Unlock()

	return c.Status(fiber.StatusOK).JSON(sdk.DeleteTopicResponse{
		Status: sdk.StatusDeleted,
		Topic:  name,
	})
}

// ListTopics returns all topics with subscriber counts
func (s *ServiceImpl) ListTopics(ctx context.Context, c *fiber.Ctx) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	topics := make([]sdk.TopicInfo, 0, len(s.Topics))
	for name, topic := range s.Topics {
		topic.Mu.RLock()
		topics = append(topics, sdk.TopicInfo{
			Name:        name,
			Subscribers: len(topic.Subscribers),
		})
		topic.Mu.RUnlock()
	}

	return c.JSON(sdk.ListTopicsResponse{Topics: topics})
}

// Health returns system health information
func (s *ServiceImpl) Health(ctx context.Context, c *fiber.Ctx) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	totalSubscribers := 0
	for _, topic := range s.Topics {
		topic.Mu.RLock()
		totalSubscribers += len(topic.Subscribers)
		topic.Mu.RUnlock()
	}

	uptimeSeconds := int(time.Since(s.Uptime).Seconds())

	return c.JSON(sdk.HealthResponse{
		UptimeSeconds: uptimeSeconds,
		Topics:        len(s.Topics),
		Subscribers:   totalSubscribers,
	})
}

// Stats returns detailed statistics per topic
func (s *ServiceImpl) Stats(ctx context.Context, c *fiber.Ctx) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := make(map[string]sdk.TopicStats)
	for name, topic := range s.Topics {
		topic.Mu.RLock()
		stats[name] = sdk.TopicStats{
			Messages:    len(topic.Messages),
			Subscribers: len(topic.Subscribers),
		}
		topic.Mu.RUnlock()
	}

	return c.JSON(sdk.StatsResponse{Topics: stats})
}

// subscriberWriter delivers messages from subscriber queue to WebSocket
// It uses the sendMessage function to ensure thread-safe writes
func (s *ServiceImpl) subscriberWriter(sub *sdk.Subscriber, topic string, sendMessage func(string, interface{})) {
	for {
		select {
		case msg := <-sub.Queue:
			sub.LastActive = time.Now()
			sendMessage("event", sdk.WebSocketResponse{
				Type:      sdk.MessageTypeEvent,
				Topic:     topic,
				Message:   &msg,
				Timestamp: msg.TS.Format(time.RFC3339),
			})
		case <-sub.CloseChannel:
			return
		}
	}
}
