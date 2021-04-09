package consume_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"cloud.google.com/go/logging"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	. "github.com/oskoss/pa-logging/consume"
	"github.com/oskoss/pa-logging/publish"
)

var _ = Describe("Consume", func() {
	var server *ghttp.Server
	BeforeEach(func() {
		server = ghttp.NewServer()
	})
	AfterEach(func() {
		server.Close()
	})

	Context("#PingEndpoint when a GET request is sent to ping path", func() {
		BeforeEach(func() {
			server.AppendHandlers(
				PingEndpoint,
			)
		})
		It("returns 200 and pong", func() {
			resp, err := http.Get(server.URL() + "/ping")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))
			body, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(body)).To(Equal("{\"message\":\"pong\"}\n"))
		})
	})
	Context("with a valid consumer", func() {
		var consumer Consumer
		var mockLogger publish.MockLogger
		BeforeEach(func() {

			testPublisher := publish.Publisher{}
			mockLogger = publish.MockLogger{}
			testPublisher.Logger = &mockLogger
			consumer = Consumer{
				LoggingPublisher: &testPublisher,
			}
		})
		Context("#MessageTrafficEndpoint when a POST request is sent to v1/message/traffic path", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					consumer.MessageTrafficEndpoint,
				)
			})
			Context("with a valid taffic log message", func() {
				var trafficLogMessage map[string]interface{}
				BeforeEach(func() {
					trafficLogMessage = map[string]interface{}{"testKey": "testValue"}
				})
				It("returns 200", func() {
					jsonReq, err := json.Marshal(trafficLogMessage)
					Expect(err).Should(BeNil())

					resp, err := http.Post(server.URL()+"/v1/message/traffic", "application/json",
						bytes.NewBuffer(jsonReq))
					Expect(err).ShouldNot(HaveOccurred())
					Expect(resp.StatusCode).Should(Equal(http.StatusOK))
				})
				FIt("publishes the message", func() {
					jsonReq, err := json.Marshal(trafficLogMessage)
					Expect(err).Should(BeNil())

					resp, err := http.Post(server.URL()+"/v1/message/traffic", "application/json", bytes.NewBuffer(jsonReq))
					Expect(err).ShouldNot(HaveOccurred())
					time.Sleep(1000)
					Expect(resp.StatusCode).Should(Equal(http.StatusOK))
					Expect(mockLogger.Messages).Should(ContainElement(logging.Entry{Payload: map[string]interface{}{"testKey": "testValue"}}))

				})
			})
		})
	})
})
