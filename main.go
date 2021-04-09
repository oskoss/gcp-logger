package main

import (
	"math/rand"
	"os"
	"strconv"
	"time"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/oskoss/pa-logging/consume"
	"github.com/oskoss/pa-logging/publish"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	projectID := os.Getenv("PROJECT_ID")
	if len(projectID) == 0 {
		logrus.WithFields(
			logrus.Fields{
				"projectID": projectID,
			}).Error("PROJECT_ID must be set")
		os.Exit(1)
	}
	rand.Seed(time.Now().UTC().UnixNano())
	loggerName := getEnv("LOGGER_NAME", petname.Generate(2, "-")+"-logger")
	loggerMaxQueueSize := getEnv("QUEUE_MAX_SIZE", "10")

	logrus.WithFields(logrus.Fields{
		"Project ID":            projectID,
		"Logger Name":           loggerName,
		"Logger Max Queue Size": loggerMaxQueueSize,
	}).Print("Starting GCP Publisher...")
	gcpPublisher, err := publish.NewGcpPublisher(loggerName, loggerMaxQueueSize, projectID)

	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"error": err,
			}).Error("Failed to start GCP publisher")
		os.Exit(1)
	}

	port := getEnv("PORT", "8080")
	portInt, err := strconv.Atoi(port)
	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"error": err,
				"PORT":  port,
			}).Error("Failed to parse PORT")
		os.Exit(1)
	}
	consume.NewConsumer(gcpPublisher, portInt)
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}
