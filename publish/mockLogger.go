package publish

import (
	"cloud.google.com/go/logging"
)

type MockLogger struct {
	Messages []logging.Entry
	Flushed  bool
}

func (logger *MockLogger) Log(E logging.Entry) {
	logger.Messages = append(logger.Messages, E)
}

func (logger *MockLogger) Flush() error {
	logger.Flushed = true
	return nil
}
