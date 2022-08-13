package controller

import (
	"errors"
	"menta-backend/db"

	"github.com/gofiber/fiber/v2"
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
				Message: `user not found`,
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

type EvenSimplerGroup struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Leader    SimpleMember `json:"leader"`
	IsLeader  bool         `json:"is_leader"`
	IsWaiting bool         `json:"is_waiting"`
}

func (u UserController) GetGroups(id string) (*[]EvenSimplerGroup, *ErrorResponse) {
	res, err := db.DB.GroupMember.FindMany(
		db.GroupMember.UserID.Equals(id),
	).With(
		db.GroupMember.Group.Fetch().With(
			db.Group.Members.Fetch(
				db.GroupMember.Leader.Equals(true),
			).With(
				db.GroupMember.User.Fetch(),
			),
		),
	).Exec(ctx)

	if err != nil {
		return nil, &ErrorResponse{
			Code:    500,
			Message: `db error`,
		}
	}

	groups := make([]EvenSimplerGroup, 0)

	for _, v := range res {
		g := v.Group()
		l := g.Members()[0].User()
		groups = append(groups, EvenSimplerGroup{
			ID:        g.ID,
			IsWaiting: v.Waiting,
			IsLeader:  v.Leader,
			Name:      g.Name,
			Leader: SimpleMember{
				ID:        l.ID,
				CreatedAt: l.CreatedAt,
				Username:  l.Username,
				Email:     ``,
				IsLeader:  true,
				IsWaiting: false,
			},
		})
	}

	return &groups, nil
}

func (u UserController) JoinGroup(userId, groupCode string) *ErrorResponse {
	group, err := db.DB.Group.FindFirst(
		db.Group.CodeLogin.Equals(true),
		db.Group.LoginCode.Equals(groupCode),
	).Exec(ctx)

	if err != nil {
		return &ErrorResponse{
			Code:    500,
			Message: `db error`,
		}
	}

	member, _ := db.DB.GroupMember.FindFirst(
		db.GroupMember.GroupID.Equals(group.ID),
		db.GroupMember.UserID.Equals(userId),
	).Exec(ctx)

	if member != nil {
		return &ErrorResponse{
			Code:    fiber.StatusConflict,
			Message: `already in group`,
		}
	}

	_, err = db.DB.GroupMember.CreateOne(
		db.GroupMember.Group.Link(
			db.Group.ID.Equals(group.ID),
		),
		db.GroupMember.User.Link(
			db.User.ID.Equals(userId),
		),
		db.GroupMember.Waiting.Set(group.RequiresAproval),
		db.GroupMember.Leader.Set(false),
	).Exec(ctx)

	if err != nil {
		return &ErrorResponse{
			Code:    500,
			Message: `db error`,
		}
	}

	return nil
}
