package app

import (
	"context"
	"flag"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ungame/command-time-track/app/cors"
	"github.com/ungame/command-time-track/app/exit"
	"github.com/ungame/command-time-track/app/handlers"
	"github.com/ungame/command-time-track/app/httpext"
	"github.com/ungame/command-time-track/app/ioext"
	"github.com/ungame/command-time-track/app/middlewares"
	"github.com/ungame/command-time-track/app/observer"
	"github.com/ungame/command-time-track/app/repository"
	"github.com/ungame/command-time-track/app/service"
	"github.com/ungame/command-time-track/db"
	"log"
	"net/http"
	"time"
)

var port int

func init() {
	flag.IntVar(&port, "p", 15555, "set port")
	flag.Parse()
}

func Run() {
	conn := db.New()

	closerGroup := ioext.NewCloserGroup(func() { ioext.Close(conn) })

	exit.OnExit(closerGroup.Close)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	closerGroup.Add(cancel)

	if err := conn.PingContext(ctx); err != nil {
		log.Panicln("unable to ping database:", err.Error())
	}

	var (
		activitiesRepository = repository.NewActivitiesRepository(context.Background(), conn)
		activitiesObserver   = observer.NewActivitiesObserver()
		activitiesService    = service.NewActivitiesService(activitiesRepository, activitiesObserver)
		activitiesHandler    = handlers.NewActivitiesHandler(activitiesService)
	)

	router := mux.NewRouter().StrictSlash(true)
	router.Use(middlewares.Logger)
	router.Path("/metrics").Handler(promhttp.Handler())
	activitiesHandler.Register(router)

	log.Printf("Listening http://localhost:%d\n\n", port)

	log.Fatalln(http.ListenAndServe(httpext.Port(port).Addr(), cors.Apply(router)))
}
