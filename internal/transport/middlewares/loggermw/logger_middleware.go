package loggermw

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"scribbler/pkg/logging"
)

func Logger() gin.HandlerFunc {
	return func(gctx *gin.Context) {
		// Start timer
		start := time.Now()
		path := gctx.Request.URL.Path
		raw := gctx.Request.URL.RawQuery

		// Process request
		gctx.Next()

		// Time the request
		stop := time.Now()
		latency := stop.Sub(start)

		// Get the correct request path
		if raw != "" {
			path = path + "?" + raw
		}

		logger := logging.Logger()

		logger.With(
			zap.String("severity", "INFO"),
			zap.String("client_ip", gctx.ClientIP()),
			zap.String("method", gctx.Request.Method),
			zap.String("path", path),
			zap.String("request_proto", gctx.Request.Proto),
			zap.Int("status_code", gctx.Writer.Status()),
			zap.Int64("latency_ns", latency.Nanoseconds()),
			zap.String("request_useragent", gctx.Request.UserAgent()),
			zap.String("error_message", gctx.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Int("body_size", gctx.Writer.Size()),
		).Info("request")
	}
}
