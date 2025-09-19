package connector

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// TunnelProxy handles HTTP proxying through tunnels
type TunnelProxy struct {
	service Service
}

// NewTunnelProxy creates a new tunnel proxy instance
func NewTunnelProxy(service Service) *TunnelProxy {
	return &TunnelProxy{
		service: service,
	}
}

// ProxyHandler handles incoming requests and forwards them through tunnels
func (tp *TunnelProxy) ProxyHandler(c *fiber.Ctx) error {
	// Extract tunnel ID from the URL path
	tunnelID := c.Params("tunnelId")
	if tunnelID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "tunnel ID is required",
		})
	}

	// Get tunnel information
	tunnel, err := tp.service.GetTunnel(c.Context(), tunnelID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "tunnel not found",
		})
	}

	if tunnel.Status != TunnelStatusActive {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": fmt.Sprintf("tunnel is %s", tunnel.Status),
		})
	}

	// Create target URL for localhost
	targetURL := fmt.Sprintf("http://localhost:%d", tunnel.LocalPort)
	target, err := url.Parse(targetURL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "invalid target URL",
		})
	}

	// Record request start time for analytics
	startTime := time.Now()

	// Create reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(target)

	// Convert Fiber context to standard HTTP
	req := c.Request()
	resp := c.Response()

	// Convert fasthttp.Request to http.Request
	httpReq, err := convertFastHTTPRequest(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to convert request",
		})
	}

	// Set the request context
	httpReq = httpReq.WithContext(c.Context())

	// Create a custom response writer to capture response details
	respWriter := &responseWriter{
		response: resp,
		status:   200,
	}

	// Forward the request
	proxy.ServeHTTP(respWriter, httpReq)

	// Calculate metrics
	responseTime := time.Since(startTime)
	endpoint := string(req.URI().Path())
	bytesTransferred := int64(len(resp.Body()))

	// Update tunnel analytics
	go tp.service.UpdateTunnelAnalytics(
		context.Background(),
		tunnelID,
		respWriter.status,
		responseTime,
		endpoint,
		bytesTransferred,
	)

	return nil
}

// responseWriter implements http.ResponseWriter for fasthttp integration
type responseWriter struct {
	response *fiber.Response
	status   int
}

func (rw *responseWriter) Header() http.Header {
	header := make(http.Header)
	rw.response.Header.VisitAll(func(key, value []byte) {
		header.Add(string(key), string(value))
	})
	return header
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	n, err := rw.response.BodyWriter().Write(data)
	return n, err
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.response.SetStatusCode(statusCode)
}

// convertFastHTTPRequest converts fasthttp.Request to http.Request
func convertFastHTTPRequest(req *fiber.Request) (*http.Request, error) {
	method := string(req.Header.Method())
	uri := string(req.URI().FullURI())

	httpReq, err := http.NewRequest(method, uri, strings.NewReader(string(req.Body())))
	if err != nil {
		return nil, err
	}

	// Copy headers
	req.Header.VisitAll(func(key, value []byte) {
		httpReq.Header.Add(string(key), string(value))
	})

	return httpReq, nil
}

// ValidatePort checks if a port is valid and available
func ValidatePort(port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}

	// Common system ports that should be restricted
	restrictedPorts := []int{22, 23, 25, 53, 80, 110, 143, 443, 993, 995}
	for _, restricted := range restrictedPorts {
		if port == restricted {
			return fmt.Errorf("port %d is restricted", port)
		}
	}

	return nil
}

// ParsePortFromString converts string to int and validates
func ParsePortFromString(portStr string) (int, error) {
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0, fmt.Errorf("invalid port number: %s", portStr)
	}

	err = ValidatePort(port)
	if err != nil {
		return 0, err
	}

	return port, nil
}

// GenerateTunnelConfig creates default tunnel configuration
func GenerateTunnelConfig(baseURL string) TunnelConfig {
	return TunnelConfig{
		MaxTunnelsPerUser: 5,
		TunnelExpiration:  24 * time.Hour, // 24 hours
		AllowedPorts:      []int{3000, 3001, 4000, 5000, 8000, 8080, 8081, 9000},
		BaseURL:           baseURL,
	}
}

// CleanupJob runs periodic cleanup of expired tunnels
func CleanupJob(service Service, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		ctx := context.Background()
		err := service.CleanupExpiredTunnels(ctx)
		if err != nil {
			// Log error (you might want to use a proper logger here)
			fmt.Printf("Error cleaning up expired tunnels: %v\n", err)
		}
	}
}
