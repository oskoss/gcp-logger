package publish

import (
	"context"
	"strconv"

	"cloud.google.com/go/logging"
	"github.com/sirupsen/logrus"
)

type Publisher struct {
	Logger    LoggingLogger
	QueueSize int
	FlushSize int
}

type LoggingLogger interface {
	Log(e logging.Entry)
	Flush() error
}

func NewGcpPublisher(logName string, flushSize string, projectId string) (*Publisher, error) {
	flushSizeInt, err := strconv.Atoi(flushSize)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	client, err := logging.NewClient(ctx, projectId)
	if err != nil {
		return nil, err
	}

	// Selects the log to write to.
	logger := client.Logger(logName)
	publisher := Publisher{
		Logger:    logger,
		FlushSize: flushSizeInt,
	}

	return &publisher, nil
}

func (p *Publisher) PublishMessage(e logging.Entry) {
	p.Logger.Log(e)
	p.QueueSize = p.QueueSize + 1
	if p.QueueSize >= p.FlushSize {
		logrus.WithFields(logrus.Fields{
			"QueueSize": p.QueueSize,
			"FlushSize": p.FlushSize,
		}).Debug("Flushing log queue")
		err := p.Logger.Flush()
		if err != nil {
			logrus.WithFields(
				logrus.Fields{
					"error": err,
				}).Error("Failed to publish message")
			return
		}
		p.QueueSize = 0
	}
}
