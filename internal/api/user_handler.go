package api

import (
	"errors"
	"log"
	"regexp"

	"github.com/gonstoll/workouts/internal/store"
)

type password struct {
	plainText *string
	hash      []byte
}

type registerUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		logger:    logger,
	}
}

func (uh *UserHandler) validateRegisterRequest(req *registerUserRequest) error {
	if req.Username == "" {
		return errors.New("Username is required")
	}

	if len(req.Username) > 50 {
		return errors.New("Username cannot be greater than 50 characters")
	}

	if req.Email == "" {
		return errors.New("Email is required")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	if !emailRegex.MatchString(req.Email) {
		return errors.New("Invalid email format")
	}

	if req.Password == "" {
		return errors.New("Password is required")
	}

	return nil
}
