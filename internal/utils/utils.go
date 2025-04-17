package utils

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func ReadIDParam(r *http.Request) (int64, error) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		return 0, errors.New("Missing ID parameter")
	}

	workoutId, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return 0, errors.New("Invalid ID parameter")
	}

	return workoutId, nil
}
