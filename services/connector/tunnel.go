package connector

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

// TunnelMessage represents a message sent through the tunnel
type TunnelMessage struct {
	ID      string                 `json:"id"`
	Type    TunnelMessageType      `json:"type"`
	Data    map[string]interface{} `json:"data"`
	Headers map[string]string      `json:"headers,omitempty"`
}

// TunnelMessageType represents the type of tunnel message
type TunnelMessageType string

const (
	MessageTypeHandshake  TunnelMessageType = "handshake"
	MessageTypeRequest    TunnelMessageType = "request"
	MessageTypeResponse   TunnelMessageType = "response"
	MessageTypeHeartbeat  TunnelMessageType = "heartbeat"
	MessageTypeError      TunnelMessageType = "error"
	MessageTypeDisconnect TunnelMessageType = "disconnect"
)

// TunnelConnection represents an active WebSocket tunnel connection
type TunnelConnection struct {
	ID          string
	UserID      string
	Subdomain   string
	LocalPort   int
	Conn        *websocket.Conn
	LastPing    time.Time
	CreatedAt   time.Time
	pendingReqs map[string]chan TunnelMessage
	mutex       sync.RWMutex
	ctx         context.Context
	cancel      context.CancelFunc
}

// TunnelManager manages all active tunnel connections
type TunnelManager struct {
	connections map[string]*TunnelConnection
	subdomains  map[string]*TunnelConnection // subdomain -> connection
	mutex       sync.RWMutex
	store       Store
}

// NewTunnelManager creates a new tunnel manager
func NewTunnelManager(store Store) *TunnelManager {
	tm := &TunnelManager{
		connections: make(map[string]*TunnelConnection),
		subdomains:  make(map[string]*TunnelConnection),
		store:       store,
	}

	// Start cleanup routine
	go tm.cleanupRoutine()

	return tm
}

// HandleWebSocket handles new WebSocket connections
func (tm *TunnelManager) HandleWebSocket(c *websocket.Conn, userID string) {
	defer c.Close()

	// Create new tunnel connection
	ctx, cancel := context.WithCancel(context.Background())
	tunnelConn := &TunnelConnection{
		ID:          uuid.New().String(),
		UserID:      userID,
		Subdomain:   generateSubdomain(userID),
		Conn:        c,
		LastPing:    time.Now(),
		CreatedAt:   time.Now(),
		pendingReqs: make(map[string]chan TunnelMessage),
		ctx:         ctx,
		cancel:      cancel,
	}

	// Register the connection
	tm.registerConnection(tunnelConn)
	defer tm.unregisterConnection(tunnelConn.ID)

	log.Printf("New tunnel connection established: %s for user %s with subdomain %s",
		tunnelConn.ID, userID, tunnelConn.Subdomain)

	// Send handshake
	handshake := TunnelMessage{
		ID:   uuid.New().String(),
		Type: MessageTypeHandshake,
		Data: map[string]interface{}{
			"tunnel_id":  tunnelConn.ID,
			"subdomain":  tunnelConn.Subdomain,
			"public_url": fmt.Sprintf("https://%s.syntra.dev", tunnelConn.Subdomain),
			"status":     "connected",
		},
	}

	if err := c.WriteJSON(handshake); err != nil {
		log.Printf("Failed to send handshake: %v", err)
		return
	}

	// Start message handling
	go tm.handleMessages(tunnelConn)

	// Keep connection alive with heartbeat
	tm.heartbeatLoop(tunnelConn)
}

// registerConnection registers a new tunnel connection
func (tm *TunnelManager) registerConnection(conn *TunnelConnection) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	tm.connections[conn.ID] = conn
	tm.subdomains[conn.Subdomain] = conn

	// Update store
	tunnel := &Tunnel{
		ID:           conn.ID,
		UserID:       conn.UserID,
		LocalPort:    conn.LocalPort,
		PublicURL:    fmt.Sprintf("https://%s.syntra.dev", conn.Subdomain),
		Status:       TunnelStatusActive,
		CreatedAt:    conn.CreatedAt,
		LastActivity: time.Now(),
		ExpiresAt:    time.Now().Add(24 * time.Hour),
		Analytics: &Analytics{
			StatusCodes: make(map[int]int),
			Endpoints:   make(map[string]int),
		},
	}

	tm.store.CreateTunnel(context.Background(), tunnel)
}

// unregisterConnection removes a tunnel connection
func (tm *TunnelManager) unregisterConnection(connID string) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	if conn, exists := tm.connections[connID]; exists {
		conn.cancel()
		delete(tm.connections, connID)
		delete(tm.subdomains, conn.Subdomain)

		// Update store
		tm.store.UpdateTunnelStatus(context.Background(), connID, TunnelStatusTerminated)

		log.Printf("Tunnel connection closed: %s", connID)
	}
}

// handleMessages handles incoming WebSocket messages
func (tm *TunnelManager) handleMessages(conn *TunnelConnection) {
	for {
		select {
		case <-conn.ctx.Done():
			return
		default:
			var msg TunnelMessage
			if err := conn.Conn.ReadJSON(&msg); err != nil {
				log.Printf("Failed to read message from tunnel %s: %v", conn.ID, err)
				return
			}

			conn.LastPing = time.Now()

			switch msg.Type {
			case MessageTypeResponse:
				tm.handleResponse(conn, msg)
			case MessageTypeHeartbeat:
				tm.handleHeartbeat(conn)
			case MessageTypeError:
				log.Printf("Error from tunnel %s: %v", conn.ID, msg.Data)
			}
		}
	}
}

// handleResponse handles response messages from the CLI
func (tm *TunnelManager) handleResponse(conn *TunnelConnection, msg TunnelMessage) {
	conn.mutex.RLock()
	responseChan, exists := conn.pendingReqs[msg.ID]
	conn.mutex.RUnlock()

	if exists {
		responseChan <- msg

		conn.mutex.Lock()
		delete(conn.pendingReqs, msg.ID)
		conn.mutex.Unlock()
	}
}

// handleHeartbeat handles heartbeat messages
func (tm *TunnelManager) handleHeartbeat(conn *TunnelConnection) {
	heartbeat := TunnelMessage{
		ID:   uuid.New().String(),
		Type: MessageTypeHeartbeat,
		Data: map[string]interface{}{
			"timestamp": time.Now().Unix(),
		},
	}

	conn.Conn.WriteJSON(heartbeat)
}

// heartbeatLoop maintains connection with periodic heartbeats
func (tm *TunnelManager) heartbeatLoop(conn *TunnelConnection) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-conn.ctx.Done():
			return
		case <-ticker.C:
			if time.Since(conn.LastPing) > 60*time.Second {
				log.Printf("Tunnel %s timed out", conn.ID)
				return
			}
			tm.handleHeartbeat(conn)
		}
	}
}

// ForwardRequest forwards an HTTP request through the tunnel
func (tm *TunnelManager) ForwardRequest(subdomain string, req *http.Request) (*http.Response, error) {
	tm.mutex.RLock()
	conn, exists := tm.subdomains[subdomain]
	tm.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("no active tunnel for subdomain: %s", subdomain)
	}

	// Create request message
	requestID := uuid.New().String()
	requestMsg := TunnelMessage{
		ID:   requestID,
		Type: MessageTypeRequest,
		Data: map[string]interface{}{
			"method": req.Method,
			"path":   req.URL.Path,
			"query":  req.URL.RawQuery,
			"body":   "",
		},
		Headers: make(map[string]string),
	}

	// Copy headers
	for key, values := range req.Header {
		if len(values) > 0 {
			requestMsg.Headers[key] = values[0]
		}
	}

	// Read request body if present
	if req.Body != nil {
		body := make([]byte, req.ContentLength)
		req.Body.Read(body)
		requestMsg.Data["body"] = string(body)
	}

	// Create response channel
	responseChan := make(chan TunnelMessage, 1)
	conn.mutex.Lock()
	conn.pendingReqs[requestID] = responseChan
	conn.mutex.Unlock()

	// Send request through WebSocket
	if err := conn.Conn.WriteJSON(requestMsg); err != nil {
		conn.mutex.Lock()
		delete(conn.pendingReqs, requestID)
		conn.mutex.Unlock()
		return nil, fmt.Errorf("failed to send request through tunnel: %w", err)
	}

	// Wait for response with timeout
	select {
	case responseMsg := <-responseChan:
		return tm.buildHTTPResponse(responseMsg)
	case <-time.After(30 * time.Second):
		conn.mutex.Lock()
		delete(conn.pendingReqs, requestID)
		conn.mutex.Unlock()
		return nil, fmt.Errorf("tunnel request timeout")
	}
}

// buildHTTPResponse constructs an HTTP response from tunnel message
func (tm *TunnelManager) buildHTTPResponse(msg TunnelMessage) (*http.Response, error) {
	// This is a simplified implementation
	// In practice, you'd need to properly construct the HTTP response
	resp := &http.Response{
		StatusCode: 200,
		Header:     make(http.Header),
	}

	if statusCode, ok := msg.Data["status_code"].(float64); ok {
		resp.StatusCode = int(statusCode)
	}

	// Set headers
	for key, value := range msg.Headers {
		resp.Header.Set(key, value)
	}

	return resp, nil
}

// GetActiveConnections returns all active connections
func (tm *TunnelManager) GetActiveConnections() map[string]*TunnelConnection {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()

	connections := make(map[string]*TunnelConnection)
	for id, conn := range tm.connections {
		connections[id] = conn
	}
	return connections
}

// cleanupRoutine periodically cleans up stale connections
func (tm *TunnelManager) cleanupRoutine() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		tm.mutex.Lock()
		for id, conn := range tm.connections {
			if time.Since(conn.LastPing) > 2*time.Minute {
				log.Printf("Cleaning up stale connection: %s", id)
				conn.cancel()
				delete(tm.connections, id)
				delete(tm.subdomains, conn.Subdomain)
			}
		}
		tm.mutex.Unlock()
	}
}

// generateSubdomain generates a unique subdomain for a user
func generateSubdomain(userID string) string {
	// Generate a shorter, user-friendly subdomain
	shortID := uuid.New().String()[:8]
	return fmt.Sprintf("%s-%s", userID, shortID)
}
