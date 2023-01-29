package middleware

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (m Middleware) LoggingMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		m.logger.WithFields(log.Fields{"remote": ctx.Request.RemoteAddr, "method": ctx.Request.Method, "uri": ctx.Request.RequestURI})

		ctx.Next()
	}
}
