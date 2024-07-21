package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authenticate(server *Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			server.logger.Error("No Authorization Header Provided")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No Authorization Header Provided"})
			c.Abort()
			return
		}

		claims, err := server.tokenHandler.VerifyToken(clientToken)
		if err != nil {
			server.logger.WithError(err).Error("Invalid Token")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
