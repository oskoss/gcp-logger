package consume

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"cloud.google.com/go/logging"
	"github.com/gin-gonic/gin"
	"github.com/oskoss/pa-logging/publish"
	"github.com/sirupsen/logrus"
)

type Consumer struct {
	LoggingPublisher *publish.Publisher
	LoggingErrors    chan error
}

func NewConsumer(pub *publish.Publisher, port int) {
	errChan := make(chan error)

	consumer := Consumer{
		LoggingPublisher: pub,
		LoggingErrors:    errChan,
	}
	r := gin.Default()
	r.GET("/ping", gin.WrapF(PingEndpoint))
	v1 := r.Group("/v1")
	{
		messageGroup := v1.Group("/message")
		{
			messageGroup.POST("/traffic", gin.WrapF(consumer.MessageTrafficEndpoint))
			// messageGroup.POST("/decrypt", messageDecryptEndpoint)
		}
	}
	r.Run(":" + strconv.Itoa(port)) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func PingEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "pong"})
}

func (c *Consumer) MessageTrafficEndpoint(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if bodyBytes == nil {
		logrus.WithFields(
			logrus.Fields{
				"body": r.Body,
			}).Error("Empty Body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"error": err,
				"body":  r.Body,
			}).Error("Failed read log payload")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	jsonMap := make(map[string]interface{})
	err = json.Unmarshal(bodyBytes, &jsonMap)
	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"error":   err,
				"payload": bodyBytes,
			}).Error("Failed parse log payload")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	go c.LoggingPublisher.PublishMessage(logging.Entry{Payload: jsonMap})
}
