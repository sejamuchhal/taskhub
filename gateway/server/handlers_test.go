package server

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func mockGinContext() *gin.Context {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c
}


func TestServer_Health(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		s    *Server
		args args
	}{ {
		name: "Healthy server",
		s: &Server{
			Logger: logrus.NewEntry(logrus.New()),
		},
		args: args{
			c: mockGinContext(),
		},
	},
		
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Health(tt.args.c)
		})
	}
}
