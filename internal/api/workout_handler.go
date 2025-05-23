package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gonstoll/workouts/internal/middleware"
	"github.com/gonstoll/workouts/internal/store"
	"github.com/gonstoll/workouts/internal/utils"
)

type WorkoutHandler struct {
	workoutStore store.WorkoutStore
	logger       *log.Logger
}

func NewWorkoutHandler(workoutStore store.WorkoutStore, logger *log.Logger) *WorkoutHandler {
	return &WorkoutHandler{workoutStore: workoutStore, logger: logger}
}

func (wh *WorkoutHandler) HandleGetWorkoutByID(w http.ResponseWriter, r *http.Request) {
	workoutId, err := utils.ReadIDParam(r)
	if err != nil {
		wh.logger.Printf("[ERROR] ReadIDParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid workout id"})
		return
	}

	user := middleware.GetUser(r)

	workout, err := wh.workoutStore.GetWorkoutByID(workoutId, user.ID)
	if err != nil {
		wh.logger.Printf("[ERROR] GetWorkoutByID: %v", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "Workout not found"})
		return
	}

	if workout == nil {
		wh.logger.Printf("[ERROR] GetWorkoutByID - no workout: %v", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "Workout not found"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": workout})
}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		wh.logger.Printf("[ERROR] Decoding on HandleCreateWorkout: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request"})
		return
	}

	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser == store.AnonymousUser {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "You must be logged in"})
		return
	}

	workout.UserID = currentUser.ID

	createdWorkotut, err := wh.workoutStore.CreateWorkout(&workout)
	if err != nil {
		wh.logger.Printf("[ERROR] CreateWorkout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to create workout"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"workout": createdWorkotut})
}

func (wh *WorkoutHandler) HandleUpdateWorkoutByID(w http.ResponseWriter, r *http.Request) {
	workoutId, err := utils.ReadIDParam(r)
	if err != nil {
		wh.logger.Printf("[ERROR] ReadIDParam %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid workout id"})
		return
	}

	user := middleware.GetUser(r)

	existingWorkout, err := wh.workoutStore.GetWorkoutByID(workoutId, user.ID)
	if err != nil {
		wh.logger.Printf("[ERROR] GetWorkoutByID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}

	if existingWorkout == nil {
		wh.logger.Printf("[ERROR] GetWorkoutByID - no workout: %v", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "Workout not found"})
		return
	}

	// NOTE: Struct to decode the request into and perform validation before encoding to
	// the main sturct.
	// Using pointers as the zero value of poitners is `nil` and we can check against that
	var updateWorkoutRequest struct {
		Title           *string              `json:"title"`
		Description     *string              `json:"description"`
		DurationMinutes *int                 `json:"duration_minutes"`
		CaloriesBurned  *int                 `json:"calories_burned"`
		Entries         []store.WorkoutEntry `json:"entries"`
	}

	err = json.NewDecoder(r.Body).Decode(&updateWorkoutRequest)
	if err != nil {
		wh.logger.Printf("[ERROR] Decoding HandleUpdateWorkoutByID on request struct: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request"})
		return
	}

	// Validation and updating main struct
	if updateWorkoutRequest.Title != nil {
		existingWorkout.Title = *updateWorkoutRequest.Title
	}
	if updateWorkoutRequest.Description != nil {
		existingWorkout.Description = *updateWorkoutRequest.Description
	}
	if updateWorkoutRequest.DurationMinutes != nil {
		existingWorkout.DurationMinutes = *updateWorkoutRequest.DurationMinutes
	}
	if updateWorkoutRequest.CaloriesBurned != nil {
		existingWorkout.CaloriesBurned = *updateWorkoutRequest.CaloriesBurned
	}
	if updateWorkoutRequest.Entries != nil {
		existingWorkout.Entries = updateWorkoutRequest.Entries
	}

	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser == store.AnonymousUser {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "You must be logged in"})
		return
	}

	workoutOwner, err := wh.workoutStore.GetWorkoutOwner(workoutId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "Workout does not exist"})
			return
		}
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}

	if workoutOwner != currentUser.ID {
		utils.WriteJSON(w, http.StatusForbidden, utils.Envelope{"error": "You are not authorized to update this workout"})
		return
	}

	err = wh.workoutStore.UpdateWorkout(existingWorkout)
	if err != nil {
		wh.logger.Printf("[ERROR] UpdateWorkout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": existingWorkout})
}

func (wh *WorkoutHandler) HandleDeleteWorkoutByID(w http.ResponseWriter, r *http.Request) {
	workoutId, err := utils.ReadIDParam(r)
	if err != nil {
		wh.logger.Printf("[ERROR] ReadIDParam %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid workout id"})
		return
	}

	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser == store.AnonymousUser {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "You must be logged in"})
		return
	}

	workoutOwner, err := wh.workoutStore.GetWorkoutOwner(workoutId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "Workout does not exist"})
			return
		}
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}

	if workoutOwner != currentUser.ID {
		utils.WriteJSON(w, http.StatusForbidden, utils.Envelope{"error": "You are not authorized to update this workout"})
		return
	}

	err = wh.workoutStore.DeleteWorkout(workoutId)
	if err == sql.ErrNoRows {
		wh.logger.Printf("[ERROR] DeleteWorkout NoRows %v", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "Workout not found"})
		return
	}
	if err != nil {
		wh.logger.Printf("[ERROR] DeleteWorkout %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, utils.Envelope{})
}
