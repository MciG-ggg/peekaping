package auth

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"peekaping/src/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type BruteforceMiddleware struct {
	bruteforceService BruteforceService
	logger            *zap.SugaredLogger
}

func NewBruteforceMiddleware(
	bruteforceService BruteforceService,
	logger *zap.SugaredLogger,
) *BruteforceMiddleware {
	return &BruteforceMiddleware{
		bruteforceService: bruteforceService,
		logger:            logger.Named("[bruteforce-middleware]"),
	}
}

// BruteforceProtection middleware checks for bruteforce attempts before login
func (m *BruteforceMiddleware) BruteforceProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only apply to login endpoint
		if c.Request.URL.Path != "/api/auth/login" || c.Request.Method != "POST" {
			c.Next()
			return
		}

		// Read the request body
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			m.logger.Errorw("Failed to read request body", "error", err)
			c.Next()
			return
		}

		// Restore the body for the next handler
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		// Parse the login DTO to get the email
		var loginDto LoginDto
		if err := json.Unmarshal(body, &loginDto); err != nil {
			// Let the controller handle the validation error
			c.Next()
			return
		}

		// Store the parsed DTO for the controller
		c.Set("login_dto", loginDto)

		// Get client IP and user agent
		clientIP := m.getClientIP(c)
		userAgent := c.GetHeader("User-Agent")

		// Check if the IP or email is currently blocked
		isBlocked, err := m.bruteforceService.IsBlocked(c.Request.Context(), loginDto.Email, clientIP)
		if err != nil {
			m.logger.Errorw("Failed to check bruteforce status", "error", err, "email", loginDto.Email, "ip", clientIP)
			// Continue with the request if we can't check the status
			c.Next()
			return
		}

		if isBlocked {
			m.logger.Warnw("Blocked login attempt", "email", loginDto.Email, "ip", clientIP, "user_agent", userAgent)
			c.JSON(http.StatusTooManyRequests, utils.NewFailResponse("Too many failed login attempts. Please try again later."))
			c.Abort()
			return
		}

		// Get current status to check for progressive delay
		status, err := m.bruteforceService.GetBruteforceStatus(c.Request.Context(), loginDto.Email, clientIP)
		if err != nil {
			m.logger.Errorw("Failed to get bruteforce status", "error", err, "email", loginDto.Email, "ip", clientIP)
			// Continue with the request if we can't get the status
			c.Next()
			return
		}

		// Apply progressive delay if required
		if status.RequiresProgressiveDelay && status.DelaySeconds > 0 {
			m.logger.Infow("Applying progressive delay", "email", loginDto.Email, "ip", clientIP, "delay_seconds", status.DelaySeconds)
			time.Sleep(time.Duration(status.DelaySeconds) * time.Second)
		}

		// Continue with the request
		c.Next()
	}
}

// getClientIP extracts the real client IP from the request
func (m *BruteforceMiddleware) getClientIP(c *gin.Context) string {
	// Check X-Forwarded-For header first (proxy/load balancer)
	xForwardedFor := c.GetHeader("X-Forwarded-For")
	if xForwardedFor != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		parts := strings.Split(xForwardedFor, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}

	// Check X-Real-IP header (nginx proxy)
	xRealIP := c.GetHeader("X-Real-IP")
	if xRealIP != "" {
		return xRealIP
	}

	// Check CF-Connecting-IP header (Cloudflare)
	cfConnectingIP := c.GetHeader("CF-Connecting-IP")
	if cfConnectingIP != "" {
		return cfConnectingIP
	}

	// Fall back to RemoteAddr
	ip := c.ClientIP()
	if ip != "" {
		return ip
	}

	return c.Request.RemoteAddr
}

// RecordLoginAttempt records a login attempt result
func (m *BruteforceMiddleware) RecordLoginAttempt(c *gin.Context, email string, success bool) {
	clientIP := m.getClientIP(c)
	userAgent := c.GetHeader("User-Agent")

	_, err := m.bruteforceService.CheckAndRecordAttempt(c.Request.Context(), email, clientIP, userAgent, success)
	if err != nil {
		m.logger.Errorw("Failed to record login attempt", "error", err, "email", email, "ip", clientIP, "success", success)
	}
}

// GetBruteforceStatus returns the current bruteforce status for debugging/monitoring
func (m *BruteforceMiddleware) GetBruteforceStatus(c *gin.Context, email string) (*BruteforceStatus, error) {
	clientIP := m.getClientIP(c)
	return m.bruteforceService.GetBruteforceStatus(c.Request.Context(), email, clientIP)
}
