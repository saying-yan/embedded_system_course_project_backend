package web_server

import (
	"github.com/gin-gonic/gin"
	. "github.com/saying-yan/embedded_system_course_project_backend/internal/logger"
	"strconv"
	"time"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		method := c.Request.Method
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		c.Next()

		statusCode := c.Writer.Status()
		latency := time.Now().Sub(start).String()
		Logger.Infof("[web] %s | %s | %d | %s", method, path, statusCode, latency)
	}
}

func DeviceIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		s := c.Param("deviceID")
		deviceID, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(200, gin.H{
				"message": "device ID error",
			})
		}
		c.Set("deviceID", deviceID)
	}
}
