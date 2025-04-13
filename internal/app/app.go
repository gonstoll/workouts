package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gonstoll/workouts/internal/api"
)

type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	workoutHander := api.NewWorkoutHandler()

	app := &Application{
		Logger:         logger,
		WorkoutHandler: workoutHander,
	}

	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Status is available\n")
}
