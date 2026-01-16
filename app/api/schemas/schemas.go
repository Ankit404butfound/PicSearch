package schemas

import "time"

type UserResponse struct {
	ID        int       `json:"id"`         // The ID of the user
	Name      string    `json:"name"`       // The name of the user
	Email     string    `json:"email"`      // The email of the user
	CreatedAt time.Time `json:"created_at"` // The creation timestamp of the user
}

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`     // The name of the user
	Email    string `json:"email" binding:"required"`    // The email of the user
	Password string `json:"password" binding:"required"` // The password of the user
}

type UpdateUserRequest struct {
	Name  string `json:"name"`  // The updated name of the user
	Email string `json:"email"` // The updated email of the user
}
