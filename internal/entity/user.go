package entity

import (
	"github.com/google/uuid"
	"github.com/prawirdani/go-midtrans-example/internal/model"
	"github.com/prawirdani/go-midtrans-example/pkg/errors"
)

var ErrUserNotFound = errors.NotFound("user not found")

type User struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
}

func NewUser(request model.UserCreateRequest) User {
	return User{
		ID:        uuid.New(),
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		Phone:     request.Phone,
	}
}
