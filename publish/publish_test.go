package publish_test

import (
	"cloud.google.com/go/logging"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/oskoss/pa-logging/publish"
)

var _ = Describe("Publish", func() {
	Context("with a valid Publisher", func() {
		l := MockLogger{}
		p := Publisher{
			Logger:    &l,
			FlushSize: 2,
		}
		Context("#PublishMessage when a valid message is recieved", func() {
			It("should write the log to GCP Cloud Logging", func() {
				logEntry := logging.Entry{Payload: "test"}
				p.PublishMessage(logEntry)
				Expect(len(l.Messages)).Should(Equal(1))
				Expect(l.Messages[0].Payload).Should(Equal("test"))
				Expect(l.Flushed).Should(BeFalse())
			})
			Context("and the queue is larger than the size specified", func() {
				It("should flush the log messages to GCP Cloud Logging", func() {
					logEntry := logging.Entry{Payload: "test2"}
					p.PublishMessage(logEntry)
					Expect(len(l.Messages)).Should(Equal(2))
					Expect(l.Messages[1].Payload).Should(Equal("test2"))
					Expect(l.Flushed).Should(BeTrue())
				})
			})
		})
	})
})
