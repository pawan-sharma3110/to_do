package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `json:"id"`
	ProfilePicture string    `json:"profile_picture"` // URL of the profile picture
	FullName       string    `json:"full_name"`
	Email          string    `json:"email"`
	PhoneNo        string    `json:"phone_no"`
	Password       string    `json:"password"`
	CreatedOn      time.Time `json:"created_on"`
}
