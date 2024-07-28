package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	pb "github.com/sejamuchhal/taskhub/gateway/pb/task"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TransformTask(task *pb.Task) *TaskDetails {
	// Convert protobuf Timestamp to a formatted string
	formatTimestamp := func(ts *timestamppb.Timestamp) string {
		if ts == nil {
			return ""
		}
		t := ts.AsTime()
		return t.Format(time.RFC3339)
	}

	return &TaskDetails{
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		DueDate:     formatTimestamp(task.DueDate),
		CreatedAt:   formatTimestamp(task.CreatedAt),
		UpdatedAt:   formatTimestamp(task.UpdatedAt),
	}
}

func getGRPCMetadataFromGin(c *gin.Context, logger *logrus.Entry) (metadata.MD, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		c.Abort()
		return nil, false
	}

	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email not found in context"})
		c.Abort()
		return nil, false
	}

	userIDStr, _ := userID.(string)
	emailStr, _ := email.(string)

	metadata := metadata.New(map[string]string{
		"user_id": userIDStr,
		"email":   emailStr,
	})

	return metadata, true
}

func hasPermission(c *gin.Context, accessibleRoles []string) bool {
	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role not found in context"})
		c.Abort()
		return false
	}

	userRole, _ := role.(string)
	for _, role := range accessibleRoles {
		if userRole == role {
			return true
		}
	}
	return false
}
