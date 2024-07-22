package common

import (
	"os"
	"strconv"
	"strings"

	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Entry

func init() {
	std := logrus.StandardLogger()
	if isDebugLoggingEnabled() {
		std.Level = logrus.DebugLevel
	}
	Logger = logrus.NewEntry(std).WithField("version", getVersion())
	Logger.Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

func isDebugLoggingEnabled() bool {
	if strings.ToLower(os.Getenv("LOG_LEVEL")) == "debug" {
		return true
	}
	debug, _ := strconv.ParseBool(os.Getenv("DEBUG"))
	return debug
}

func getVersion() string {
	v := os.Getenv("VERSION")
	if v == "" {
		return "unknown"
	}
	return v
}
