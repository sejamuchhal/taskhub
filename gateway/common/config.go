package common

import (
	"fmt"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

type Config struct {
	Logger         *logrus.Entry
	HTTPPort       int64
	AuthServiceUrl string
	TaskServiceUrl string
}

func LoadConfig() (*Config, error) {
	httpPort, err := strconv.ParseInt(getEnv("GATEWAY_PORT", "3000"), 10, 32)
	if err != nil {
		return nil, fmt.Errorf("missing required port environment variables")
	}


	authPort := getEnv("AUTH_PORT", "4040")
	authServiceUrl := fmt.Sprintf("auth:%v", authPort)

	taskPort := getEnv("TASK_PORT", "8080")
	taskServiceUrl := fmt.Sprintf("task:%v", taskPort)
	

	config := &Config{
		Logger:         Logger,
		HTTPPort:       httpPort,
		AuthServiceUrl: authServiceUrl,
		TaskServiceUrl: taskServiceUrl,
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
