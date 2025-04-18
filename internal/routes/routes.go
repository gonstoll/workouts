package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/gonstoll/workouts/internal/app"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()
	// Health
	r.Get("/health", app.HealthCheck)

	// Workouts
	r.Get("/workouts/{id}", app.WorkoutHandler.HandleGetWorkoutByID)
	r.Post("/workouts", app.WorkoutHandler.HandleCreateWorkout)
	r.Put("/workouts/{id}", app.WorkoutHandler.HandleUpdateWorkoutByID)
	r.Delete("/workouts/{id}", app.WorkoutHandler.HandleDeleteWorkoutByID)

	// Users
	r.Post("/users", app.UserHandler.HandleRegisterUser)

	// Tokens
	r.Post("/token/authentication", app.TokenHandler.HandleCreateToken)

	return r
}
