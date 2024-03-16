package main

import (
	"go-sse-poc/controllers"
	"go-sse-poc/events"
	"log"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

func GetServer() string {
	return "hi"
}

func main() {

	// -------------------------------------
	// ECHO ROUTER INITIALISATION
	// -------------------------------------
	e := echo.New()
	e.Use(middleware.Logger())
	// e.Use(middleware.Recover()) // only works with panics, not 'fatal'
	// e.Renderer = NewTemplates()

	e.GET("/events", events.SseTest)

	e.POST("/api/nodes", controllers.UpdateNode)

	/* Start the http server */
	go e.Logger.Fatal(e.Start(port))

}
