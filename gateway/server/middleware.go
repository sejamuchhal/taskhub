package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	pb "github.com/sejamuchhal/taskhub/gateway/pb/auth"
)

func Authenticate(server *Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorization := ctx.GetHeader("Authorization")

		if authorization == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			ctx.Abort()
			return
		}

		tokenParts := strings.SplitN(authorization, "Bearer ", 2)

		if len(tokenParts) != 2 {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			ctx.Abort()
			return
		}

		token := tokenParts[1]

		res, err := server.AuthClient.Validate(context.Background(), &pb.ValidateRequest{
			Token: token,
		})

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			ctx.Abort()
			return
		}

		ctx.Set("user_id", res.UserId)
		ctx.Set("email", res.Email)

		// Continue to the next middleware (main handler)
		ctx.Next()
	}
}

