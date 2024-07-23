package common

import (
	"os"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Entry

func init() {
	std := logrus.StandardLogger()

	if getLogFormat() == "json" {
		std.SetFormatter(&logrus.JSONFormatter{})
	} else {
		std.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	if isDebugLoggingEnabled() {
		std.Level = logrus.DebugLevel
	} else {
		std.Level = logrus.InfoLevel
	}

	Logger = logrus.NewEntry(std).WithField("version", getVersion())
}

func isDebugLoggingEnabled() bool {
	if strings.ToLower(os.Getenv("LOG_LEVEL")) == "debug" {
		return true
	}
	debug, _ := strconv.ParseBool(os.Getenv("DEBUG"))
	return debug
}

func getLogFormat() string {
	format := os.Getenv("LOG_FORMAT")
	if format == "" {
		return "text"
	}
	return strings.ToLower(format)
}

func getVersion() string {
	v := os.Getenv("VERSION")
	if v == "" {
		return "unknown"
	}
	return v
}
