package controller

import (
	"errors"
	"menta-backend/db"
)

type UserController struct{}

type UserResponse struct {
	ID        string      `json:"id"`
	CreatedAt db.DateTime `json:"createdAt"`
	UpdatedAt db.DateTime `json:"updatedAt"`
	Username  string      `json:"username"`
	Email     string      `json:"email"`
	IsTeacher bool        `json:"isTeacher"`
}

func (u UserController) GetUserById(id string) (*UserResponse, *ErrorResponse) {
	res, err := db.DB.User.FindFirst(
		db.User.ID.Equals(id),
	).Exec(ctx)

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, &ErrorResponse{
				Code:    404,
				Message: `user not fount`,
			}
		}
		return nil, &ErrorResponse{
			Code:    500,
			Message: `db error`,
		}
	}
	return &UserResponse{
		ID:        res.ID,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
		Username:  res.Username,
		Email:     res.Email,
		IsTeacher: res.IsTeacher,
	}, nil
}
