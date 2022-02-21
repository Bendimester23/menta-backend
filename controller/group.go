package controller

import (
	"log"
	"math/rand"
	"menta-backend/db"
	"menta-backend/models"
	"strings"
)

type GroupController struct{}

func (g GroupController) Create(leaderId string, data models.CreateGroup) (*db.GroupModel, *ErrorResponse) {

	loginCode := ""

	if data.JoinWithCode {
		loginCode = data.CodePrefix + "_" + g.GenerateLoginCode()
	}

	log.Println(data)

	res, err := db.DB.Group.CreateOne(
		db.Group.Name.Set(data.Name),
		db.Group.CodeLogin.Set(data.JoinWithCode),
		db.Group.LoginCode.Set(loginCode),
		db.Group.RequiresAproval.Set(data.LimitJoins),
	).Exec(ctx)

	if err != nil {
		log.Println(err)
		return nil, &ErrorResponse{
			Code:    500,
			Message: "db error",
		}
	}

	_, err = db.DB.GroupMember.CreateOne(
		db.GroupMember.Group.Link(
			db.Group.ID.Equals(res.ID),
		),
		db.GroupMember.User.Link(
			db.User.ID.Equals(leaderId),
		),
		db.GroupMember.Waiting.Set(false),
		db.GroupMember.Leader.Set(true),
	).Exec(ctx)

	if err != nil {
		log.Println(err)
		return nil, &ErrorResponse{
			Code:    500,
			Message: "db error",
		}
	}

	return res, nil
}

type SimpleGroup struct {
	ID              string         `json:"id"`
	CreatedAt       db.DateTime    `json:"createdAt"`
	UpdatedAt       db.DateTime    `json:"updatedAt"`
	Name            string         `json:"name"`
	CodeLogin       bool           `json:"codeLogin"`
	LoginCode       string         `json:"loginCode"`
	RequiresAproval bool           `json:"requiresAproval"`
	Members         []SimpleMember `json:"members"`
}

type SimpleMember struct {
	ID        string      `json:"id"`
	CreatedAt db.DateTime `json:"createdAt"`
	Username  string      `json:"username"`
	Email     string      `json:"email"`
	IsLeader  bool        `json:"is_leader"`
	IsWaiting bool        `json:"is_waiting"`
}

func (g GroupController) ById(id string, isLeader bool) (*SimpleGroup, *ErrorResponse) {
	res, err := db.DB.Group.FindFirst(
		db.Group.ID.Equals(id),
	).With(
		db.Group.Members.Fetch().With(
			db.GroupMember.User.Fetch(),
		),
	).Exec(ctx)

	if err != nil {
		log.Println(err)
		return nil, &ErrorResponse{
			Code:    500,
			Message: "db error",
		}
	}

	r := make([]SimpleMember, 0)

	for _, v := range res.Members() {
		u := v.User()
		if isLeader {
			r = append(r, SimpleMember{
				ID:        u.ID,
				CreatedAt: u.CreatedAt,
				Username:  u.Username,
				Email:     u.Email,
				IsLeader:  v.Leader,
				IsWaiting: v.Waiting,
			})
		} else {
			r = append(r, SimpleMember{
				ID:        u.ID,
				CreatedAt: u.CreatedAt,
				Username:  u.Username,
				IsLeader:  v.Leader,
				IsWaiting: v.Waiting,
			})
		}
	}

	return &SimpleGroup{
		ID:              res.ID,
		CreatedAt:       res.CreatedAt,
		UpdatedAt:       res.UpdatedAt,
		Name:            res.Name,
		CodeLogin:       res.CodeLogin,
		LoginCode:       res.LoginCode,
		RequiresAproval: res.RequiresAproval,
		Members:         r,
	}, nil
}

var firstChar = strings.Split("QWERTZUIOPASDFGHJKLYXCVBNM", "")
var normalChar = strings.Split("QWERTZUIOPASDFGHJKLYXCVBNMqwertzuiopasdfghjklyxcvbnm0123456789", "")

func (GroupController) GenerateLoginCode() string {
	first := firstChar[rand.Intn(len(firstChar))]
	end := ""
	for i := 0; i < 7; i++ {
		end = end + normalChar[rand.Intn(len(normalChar))]
	}
	return first + end
}
