package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/gonstoll/workouts/internal/app"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(app.Middleware.Authenticate)

		// Workouts
		r.Get("/workouts/{id}", app.Middleware.RequireUser(app.WorkoutHandler.HandleGetWorkoutByID))
		r.Post("/workouts", app.Middleware.RequireUser(app.WorkoutHandler.HandleCreateWorkout))
		r.Put("/workouts/{id}", app.Middleware.RequireUser(app.WorkoutHandler.HandleUpdateWorkoutByID))
		r.Delete("/workouts/{id}", app.Middleware.RequireUser(app.WorkoutHandler.HandleDeleteWorkoutByID))
	})

	// Health
	r.Get("/health", app.HealthCheck)

	// Users
	r.Post("/users", app.UserHandler.HandleRegisterUser)

	// Tokens
	r.Post("/token/authentication", app.TokenHandler.HandleCreateToken)

	return r
}
