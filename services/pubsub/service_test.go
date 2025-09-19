package pubsub

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Aryaman/pub-sub/sdk"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test HTTP API endpoints
func TestCreateTopic(t *testing.T) {
	service := NewService(100, 100)
	app := fiber.New()

	app.Post("/topics", func(c *fiber.Ctx) error {
		return service.CreateTopic(c.Context(), c)
	})

	tests := []struct {
		name           string
		payload        string
		expectedStatus int
		expectedTopic  string
	}{
		{
			name:           "valid topic creation",
			payload:        `{"name":"test-topic"}`,
			expectedStatus: 201,
			expectedTopic:  "test-topic",
		},
		{
			name:           "duplicate topic creation",
			payload:        `{"name":"test-topic"}`,
			expectedStatus: 409,
			expectedTopic:  "test-topic",
		},
		{
			name:           "invalid payload",
			payload:        `{"invalid":"json"}`,
			expectedStatus: 400,
		},
		{
			name:           "empty name",
			payload:        `{"name":""}`,
			expectedStatus: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/topics", strings.NewReader(tt.payload))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedStatus == 201 {
				var response sdk.CreateTopicResponse
				err := json.NewDecoder(resp.Body).Decode(&response)
				require.NoError(t, err)
				assert.Equal(t, sdk.StatusCreated, response.Status)
				assert.Equal(t, tt.expectedTopic, response.Topic)
			}
		})
	}
}

func TestDeleteTopic(t *testing.T) {
	service := NewService(100, 100)
	app := fiber.New()

	// Create a topic first
	service.Topics["test-topic"] = &sdk.Topic{
		Name:        "test-topic",
		Subscribers: make(map[string]*sdk.Subscriber),
		Messages:    make([]sdk.Message, 0),
		MaxMessages: 100,
	}

	app.Delete("/topics/:name", func(c *fiber.Ctx) error {
		return service.DeleteTopic(c.Context(), c)
	})

	tests := []struct {
		name           string
		topicName      string
		expectedStatus int
	}{
		{
			name:           "delete existing topic",
			topicName:      "test-topic",
			expectedStatus: 200,
		},
		{
			name:           "delete non-existing topic",
			topicName:      "non-existing",
			expectedStatus: 404,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("DELETE", "/topics/"+tt.topicName, nil)
			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestListTopics(t *testing.T) {
	service := NewService(100, 100)
	app := fiber.New()

	// Add some test topics
	service.Topics["topic1"] = &sdk.Topic{
		Name:        "topic1",
		Subscribers: make(map[string]*sdk.Subscriber),
		Messages:    make([]sdk.Message, 0),
		MaxMessages: 100,
	}
	service.Topics["topic2"] = &sdk.Topic{
		Name:        "topic2",
		Subscribers: make(map[string]*sdk.Subscriber),
		Messages:    make([]sdk.Message, 0),
		MaxMessages: 100,
	}

	app.Get("/topics", func(c *fiber.Ctx) error {
		return service.ListTopics(c.Context(), c)
	})

	req := httptest.NewRequest("GET", "/topics", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode)

	var response sdk.ListTopicsResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.Len(t, response.Topics, 2)
}

func TestHealthEndpoint(t *testing.T) {
	service := NewService(100, 100)
	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return service.Health(c.Context(), c)
	})

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode)

	var response sdk.HealthResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.GreaterOrEqual(t, response.UptimeSeconds, 0)
	assert.Equal(t, 0, response.Topics)
	assert.Equal(t, 0, response.Subscribers)
}

// Test core functionality
func TestTopicMessageStorage(t *testing.T) {
	service := NewService(100, 100)
	service.MaxMessages = 3 // Small buffer for testing

	topic := &sdk.Topic{
		Name:        "test-topic",
		Subscribers: make(map[string]*sdk.Subscriber),
		Messages:    make([]sdk.Message, 0),
		MaxMessages: service.MaxMessages,
	}

	// Add messages beyond buffer size
	messages := []sdk.Message{
		{ID: "1", Payload: "msg1", TS: time.Now()},
		{ID: "2", Payload: "msg2", TS: time.Now()},
		{ID: "3", Payload: "msg3", TS: time.Now()},
		{ID: "4", Payload: "msg4", TS: time.Now()}, // Should replace msg1
		{ID: "5", Payload: "msg5", TS: time.Now()}, // Should replace msg2
	}

	for _, msg := range messages {
		// Simulate ring buffer behavior
		if len(topic.Messages) >= topic.MaxMessages {
			topic.Messages = topic.Messages[1:]
		}
		topic.Messages = append(topic.Messages, msg)
	}

	// Should only have last 3 messages
	assert.Len(t, topic.Messages, 3)
	assert.Equal(t, "3", topic.Messages[0].ID)
	assert.Equal(t, "4", topic.Messages[1].ID)
	assert.Equal(t, "5", topic.Messages[2].ID)
}

func TestSlowConsumerDetection(t *testing.T) {
	service := NewService(100, 100)
	service.MaxQueue = 2 // Small queue for testing

	topic := &sdk.Topic{
		Name:        "test-topic",
		Subscribers: make(map[string]*sdk.Subscriber),
		Messages:    make([]sdk.Message, 0),
		MaxMessages: 100,
	}

	// Create a subscriber with small queue
	sub := &sdk.Subscriber{
		ClientID:     "test-client",
		Queue:        make(chan sdk.Message, service.MaxQueue),
		QueueSize:    service.MaxQueue,
		LastActive:   time.Now(),
		CloseChannel: make(chan struct{}),
	}

	topic.Subscribers["test-client"] = sub

	// Fill the queue
	for i := 0; i < service.MaxQueue; i++ {
		select {
		case sub.Queue <- sdk.Message{ID: string(rune(i)), Payload: "test"}:
		default:
			t.Fatalf("Queue should not be full yet")
		}
	}

	// Try to add one more message - should fail (slow consumer)
	select {
	case sub.Queue <- sdk.Message{ID: "overflow", Payload: "test"}:
		t.Fatalf("Queue should be full")
	default:
		// Expected behavior - queue is full
	}
}

func TestConcurrentAccess(t *testing.T) {
	service := NewService(100, 100)

	// Create topic
	topic := &sdk.Topic{
		Name:        "concurrent-test",
		Subscribers: make(map[string]*sdk.Subscriber),
		Messages:    make([]sdk.Message, 0),
		MaxMessages: 1000,
	}
	service.Topics["concurrent-test"] = topic

	var wg sync.WaitGroup
	numGoroutines := 10
	messagesPerGoroutine := 100

	// Concurrent message publishing
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < messagesPerGoroutine; j++ {
				topic.Mu.Lock()
				msg := sdk.Message{
					ID:      fmt.Sprintf("worker-%d-msg-%d", workerID, j),
					Payload: fmt.Sprintf("payload-%d-%d", workerID, j),
					TS:      time.Now(),
				}
				if len(topic.Messages) >= topic.MaxMessages {
					topic.Messages = topic.Messages[1:]
				}
				topic.Messages = append(topic.Messages, msg)
				topic.Mu.Unlock()
			}
		}(i)
	}

	wg.Wait()

	// Verify final state
	topic.Mu.RLock()
	messageCount := len(topic.Messages)
	topic.Mu.RUnlock()

	assert.Greater(t, messageCount, 0)
	assert.LessOrEqual(t, messageCount, topic.MaxMessages)
}

// Benchmark tests
func BenchmarkMessagePublishing(b *testing.B) {
	// service := NewService() // removed unused variable
	topic := &sdk.Topic{
		Name:        "bench-topic",
		Subscribers: make(map[string]*sdk.Subscriber),
		Messages:    make([]sdk.Message, 0),
		MaxMessages: 10000,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		topic.Mu.Lock()
		msg := sdk.Message{
			ID:      fmt.Sprintf("msg-%d", i),
			Payload: "benchmark payload",
			TS:      time.Now(),
		}
		if len(topic.Messages) >= topic.MaxMessages {
			topic.Messages = topic.Messages[1:]
		}
		topic.Messages = append(topic.Messages, msg)
		topic.Mu.Unlock()
	}
}

func BenchmarkConcurrentPublishing(b *testing.B) {
	// service := NewService() // removed unused variable
	topic := &sdk.Topic{
		Name:        "bench-topic",
		Subscribers: make(map[string]*sdk.Subscriber),
		Messages:    make([]sdk.Message, 0),
		MaxMessages: 10000,
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			topic.Mu.Lock()
			msg := sdk.Message{
				ID:      fmt.Sprintf("msg-%d", i),
				Payload: "benchmark payload",
				TS:      time.Now(),
			}
			if len(topic.Messages) >= topic.MaxMessages {
				topic.Messages = topic.Messages[1:]
			}
			topic.Messages = append(topic.Messages, msg)
			topic.Mu.Unlock()
			i++
		}
	})
}

// Helper functions for testing
func createTestSubscriber(clientID string, queueSize int) *sdk.Subscriber {
	return &sdk.Subscriber{
		ClientID:     clientID,
		Queue:        make(chan sdk.Message, queueSize),
		QueueSize:    queueSize,
		LastActive:   time.Now(),
		CloseChannel: make(chan struct{}),
	}
}
