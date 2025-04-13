package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/gonstoll/workouts/internal/app"
	"github.com/gonstoll/workouts/internal/routes"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "Server port")
	flag.Parse()

	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}
	defer app.DB.Close()

	r := routes.SetupRoutes(app)

	server := http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf(":%d", port),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.Logger.Printf("App is running on port %d", port)

	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatal(err)
	}
}
