package controllers

import (
	"go-sse-poc/events"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func init() {
	log.Default().SetFlags(log.Llongfile) // | log.Ldate | log.Ltime)
	log.Println("starting up")
}

func UpdateNode(c echo.Context) error {
	log.Print("Fake updating a node to test SSE publishing")

	// publish something here
	m := events.NewMessage("hi", "topic")
	_ = events.GetSseServer().Publish(m)

	return c.String(http.StatusOK, "published a message, did you get it?\n")
}
