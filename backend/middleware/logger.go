package middleware

import (
	"montelukast/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware(c *gin.Context) {
	log := logger.Log
	startTime := time.Now()
	c.Next()
	endTime := time.Now()
	latency := endTime.Sub(startTime).String()
	reqMethod := c.Request.Method
	reqHost := c.Request.Host
	reqURI := c.Request.RequestURI
	statusCode := c.Writer.Status()
	clientIP := c.ClientIP()

	fields := map[string]any{
		"method":    reqMethod,
		"uri":       reqURI,
		"status":    statusCode,
		"latency":   latency,
		"client_ip": clientIP,
		"host":      reqHost,
	}
	if lastErr := c.Errors.Last(); lastErr != nil {
		log.WithFields(fields).Error(c.Errors[0])
		return
	}

	log.WithFields(fields).Infof("REQUEST %s %s SUCCESS", reqMethod, reqURI)
}
