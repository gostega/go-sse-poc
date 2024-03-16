package main

import (
	"fmt"
	"go-sse-poc/controllers"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tmaxmax/go-sse"
)

var portRaw string = "1988"
var port string = ":" + portRaw

func init() {

	/*
		Initialise logging
	*/
	log.Default().SetFlags(log.Llongfile) // | log.Ldate | log.Ltime)
	// log.Default().SetPrefix(utils.GetCurrentFunction() + ": ")

	/* Print opening message */
	log.Printf("Starting main go program at %v\n", time.Now())

}

const (
	topicRandomNumbers = "numbers"
	topicStateChanges  = "evt_statechange"
)

func sseTest(c echo.Context) error {
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

func main() {

	// -------------------------------------
	// ECHO ROUTER INITIALISATION
	// -------------------------------------
	e := echo.New()
	e.Use(middleware.Logger())
	// e.Use(middleware.Recover()) // only works with panics, not 'fatal'
	// e.Renderer = NewTemplates()

	e.GET("/events", sseTest)

	e.POST("/api/nodes", controllers.UpdateNode)

	/* Start the http server */
	go e.Logger.Fatal(e.Start(port))

}
