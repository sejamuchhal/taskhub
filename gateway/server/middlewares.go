package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	pb "github.com/sejamuchhal/taskhub/gateway/pb/auth"
)

func Authenticate(server *Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Access")

		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			ctx.Abort()
			return
		}

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
		ctx.Set("role", res.Role)

		// Continue to the next middleware (main handler)
		ctx.Next()
	}
}

func RegisterPrometheusMatrics() {
	prometheus.MustRegister(latency)
}

func RecordRequestLatency(ctx *gin.Context) {
	start := time.Now()
	ctx.Next()
	elapased := time.Since(start).Seconds()
	latency.WithLabelValues(
		ctx.Request.Method,
		ctx.Request.URL.Path,
	).Observe(elapased)
}

var latency = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Namespace:  "api",
		Name:       "latency_seconds",
		Help:       "Request latency distributions.",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	},
	[]string{"method", "path"},
)
