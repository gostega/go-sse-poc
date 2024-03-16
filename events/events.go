package events

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/tmaxmax/go-sse"
)

const (
	topicRandomNumbers = "numbers"
	topicStateChanges  = "evt_statechange"
)

func GetSseServer() *sse.Server {
	return sseHandler
}

func NewMessage(message, topic string) *sse.Message {
	m := &sse.Message{}
	m.ID = sse.ID(topic)
	m.AppendData(message)
	// _ = sseHandler.Publish(m)
	return m
}

func SseTest(c echo.Context) error {
	// s := &sse.Server{}

	// seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	// const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	go func() {

		// randID := seededRand.Intn(10)
		// b := make([]byte, 10)

		// for range time.Tick(time.Second) {
		// 	m := &sse.Message{}
		// 	m.ID = sse.ID(fmt.Sprintf("thread:%d", randID))
		// 	for i := range b {
		// 		b[i] = charset[seededRand.Intn(len(charset))]
		// 	}

		// 	m.AppendData(string(b))
		// 	_ = sseHandler.Publish(m)
		// }
		m := &sse.Message{}
		m.ID = sse.ID(topicStateChanges)
		m.AppendData("Connected")
		_ = sseHandler.Publish(m)

	}()

	sseHandler.ServeHTTP(c.Response(), c.Request())

	return nil
}

var sseHandler = &sse.Server{
	Provider: &sse.Joe{
		ReplayProvider: &sse.ValidReplayProvider{
			TTL:        time.Minute * 5,
			GCInterval: time.Minute,
			AutoIDs:    false,
		},
	},
	Logger: nil,
	OnSession: func(s *sse.Session) (sse.Subscription, bool) {
		topics := s.Req.URL.Query()["topic"]
		for _, topic := range topics {
			if topic != topicRandomNumbers && topic != topicStateChanges {
				fmt.Fprintf(s.Res, "invalid topic %q; supported are %q, %q", topic, topicRandomNumbers, topicStateChanges)
				s.Res.WriteHeader(http.StatusBadRequest)
				return sse.Subscription{}, false
			}
		}
		if len(topics) == 0 {
			// Provide default topics, if none are given.
			topics = []string{topicRandomNumbers, topicStateChanges}
		}

		return sse.Subscription{
			Client:      s,
			LastEventID: s.LastEventID,
			Topics:      append(topics, sse.DefaultTopic), // the shutdown message is sent on the default topic
		}, true
	},
}
