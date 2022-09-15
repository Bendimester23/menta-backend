package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"menta-backend/db"
	"menta-backend/models"
	"strings"
	"time"
)

type GroupController struct{}

func (g GroupController) Create(leaderId string, data models.CreateGroup) (*db.GroupModel, *ErrorResponse) {

	loginCode := ""

	if data.JoinWithCode {
		loginCode = data.CodePrefix + "_" + g.GenerateLoginCode()
	}

	chat, err := db.DB.ChatRoom.CreateOne(
		db.ChatRoom.Description.Set(`szoba cucc`),
	).Exec(ctx)
	if err != nil {
		log.Println(err)
		return nil, &ErrorResponse{
			Code:    500,
			Message: "db error",
		}
	}

	res, err := db.DB.Group.CreateOne(
		db.Group.Name.Set(data.Name),
		db.Group.CodeLogin.Set(data.JoinWithCode),
		db.Group.LoginCode.Set(loginCode),
		db.Group.RequiresAproval.Set(data.LimitJoins),
		db.Group.Room.Link(
			db.ChatRoom.ID.Equals(chat.ID),
		),
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

	_, err = db.DB.ChatMember.CreateOne(
		db.ChatMember.User.Link(
			db.User.ID.Equals(leaderId),
		),
		db.ChatMember.Room.Link(
			db.ChatRoom.ID.Equals(chat.ID),
		),
		db.ChatMember.Nickname.Set(""),
	).Exec(ctx)

	if err != nil {
		return nil, &ErrorResponse{
			Code:    500,
			Message: `db error`,
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
	RoomId          string         `json:"roomId"`
	RoomName        string         `json:"roomName"`
}

type SimpleMember struct {
	ID        string      `json:"id"`
	CreatedAt db.DateTime `json:"createdAt"`
	JoinedAt  db.DateTime `json:"joinedAt"`
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
		db.Group.Room.Fetch(),
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
				JoinedAt:  v.JoinedAt,
			})
		} else {
			r = append(r, SimpleMember{
				ID:        u.ID,
				CreatedAt: u.CreatedAt,
				Username:  u.Username,
				IsLeader:  v.Leader,
				IsWaiting: v.Waiting,
				JoinedAt:  v.JoinedAt,
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
		RoomId:          res.RoomID,
		RoomName:        res.Room().Description,
	}, nil
}

type examAnswerData struct {
	Answer  bool
	Answers []string
}

type examQuestionData struct {
	Name            string
	Type            string
	AutoGrade       bool
	CaseSensitiveTS bool
	MinTS           int
	MaxTS           int
}

func (g *GroupController) CreateExam(groupId string, data models.CreateExam) (*db.ExamModel, error) {
	answers := make([]examAnswerData, 0)
	questions := make([]examQuestionData, 0)

	for _, v := range data.Questions {
		answers = append(answers, examAnswerData{
			Answer:  v.Data.AnswerTF,
			Answers: v.Data.AnswersTS,
		})

		questions = append(questions, examQuestionData{
			Name:            v.Data.Name,
			Type:            v.Data.Type,
			AutoGrade:       v.Data.AutoGrade,
			CaseSensitiveTS: v.Data.CaseSensitiveTS,
			MinTS:           v.Data.MinTS,
			MaxTS:           v.Data.MaxTS,
		})
	}

	answersData, err := json.Marshal(answers)
	if err != nil {
		return nil, err
	}

	questionsData, err := json.Marshal(questions)
	if err != nil {
		return nil, err
	}

	return db.DB.Exam.CreateOne(
		db.Exam.Group.Link(
			db.Group.ID.Equals(groupId),
		),
		db.Exam.Title.Set(data.Name),
		db.Exam.Description.Set(data.Description),
		db.Exam.Questions.Set(db.JSON(questionsData)),
		db.Exam.Answers.Set(db.JSON(answersData)),
		db.Exam.StartsAt.Set(db.DateTime(time.Now())),
		db.Exam.EndsAt.Set(db.DateTime(time.Now().Add(time.Hour*72))),
		db.Exam.MaxLenght.Set(int(time.Hour.Milliseconds())),
	).Exec(ctx)
}

type simpleExam struct {
	Id                 string `json:"id"`
	Title              string `json:"title"`
	Description        string `json:"description"`
	StartsAt           string `json:"starts_at"`
	EndsAt             string `json:"ends_at"`
	MaxLenght          int    `json:"max_lenght"`
	SolutionsSubmitted int    `json:"solutions_submitted"`
}

func (g *GroupController) GetGroupExams(id string, userId string) ([]simpleExam, *ErrorResponse) {
	member, err := db.DB.GroupMember.FindFirst(
		db.GroupMember.User.Where(
			db.User.ID.Equals(userId),
		),
		db.GroupMember.Group.Where(
			db.Group.ID.Equals(id),
		),
	).Exec(ctx)

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, &ErrorResponse{
				Code:    403,
				Message: "You don't have permission to access this group!",
			}
		}
		return nil, &ErrorResponse{
			Code:    500,
			Message: "DB error!",
		}
	}

	var exams []db.ExamModel

	if member.Leader {
		exams, err = db.DB.Exam.FindMany(
			db.Exam.Group.Where(
				db.Group.ID.Equals(id),
			),
		).Exec(ctx)
	} else {
		exams, err = db.DB.Exam.FindMany(
			db.Exam.Group.Where(
				db.Group.ID.Equals(id),
			),
			db.Exam.Show.Equals(true),
		).Exec(ctx)
	}
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return make([]simpleExam, 0), nil
		}
		return nil, &ErrorResponse{
			Code:    500,
			Message: "DB error!",
		}
	}

	simplified := make([]simpleExam, 0)

	for _, v := range exams {
		simplified = append(simplified, simpleExam{
			Id:          v.ID,
			Title:       v.Title,
			Description: v.Description,
			StartsAt:    fmt.Sprintf("%d", v.StartsAt.UnixMilli()),
			EndsAt:      fmt.Sprintf("%d", v.EndsAt.UnixMilli()),
			MaxLenght:   v.MaxLenght,
			//TODO: count submitted solutions
			SolutionsSubmitted: 0,
		})
	}

	return simplified, nil
}

func (g *GroupController) GetExamById(groupId string, userId string, examId string) (*models.CreateExam, *ErrorResponse) {
	member, err := db.DB.GroupMember.FindFirst(
		db.GroupMember.User.Where(
			db.User.ID.Equals(userId),
		),
		db.GroupMember.Group.Where(
			db.Group.ID.Equals(groupId),
		),
	).Exec(ctx)

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, &ErrorResponse{
				Code:    403,
				Message: "You don't have permission to access this group!",
			}
		}
		return nil, &ErrorResponse{
			Code:    500,
			Message: "DB error!",
		}
	}

	res, err := db.DB.Exam.FindMany(
		db.Exam.ID.Equals(examId),
		db.Exam.Group.Where(
			db.Group.ID.Equals(groupId),
		),
	).With(
		db.Exam.Solutions.Fetch(),
	).Exec(ctx)

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, &ErrorResponse{
				Code:    404,
				Message: "Exam not found!",
			}
		}
		return nil, &ErrorResponse{
			Code:    500,
			Message: "DB error!",
		}
	}

	if len(res) == 0 {
		log.Println(res, err)
		return nil, &ErrorResponse{
			Code:    404,
			Message: "Exam not found!",
		}
	}

	dbExam := res[0]

	if !dbExam.Show && !member.Leader {
		return nil, &ErrorResponse{
			Code:    403,
			Message: "You don't have permission to access this.",
		}
	}

	questions := make([]models.CreateExamQuestion, 0)

	var dbQuestions []examQuestionData
	if json.Unmarshal(dbExam.Questions, &dbQuestions) != nil {
		return nil, &ErrorResponse{
			Code:    500,
			Message: "Internal Server Error",
		}
	}

	if dbExam.ShowResults || member.Leader {
		var dbAnswers []examAnswerData
		if json.Unmarshal(dbExam.Answers, &dbAnswers) != nil {
			return nil, &ErrorResponse{
				Code:    500,
				Message: "Internal Server Error",
			}
		}

		for k, v := range dbQuestions {
			a := dbAnswers[k]
			questions = append(questions, models.CreateExamQuestion{
				Id: k,
				Data: models.QuestionData{
					Name:            v.Name,
					Type:            v.Type,
					AutoGrade:       v.AutoGrade,
					CaseSensitiveTS: v.CaseSensitiveTS,
					MinTS:           v.MinTS,
					MaxTS:           v.MaxTS,
					AnswerTF:        a.Answer,
					AnswersTS:       a.Answers,
				},
			})
		}
	} else {
		for k, v := range dbQuestions {
			questions = append(questions, models.CreateExamQuestion{
				Id: k,
				Data: models.QuestionData{
					Name:            v.Name,
					Type:            v.Type,
					AutoGrade:       v.AutoGrade,
					CaseSensitiveTS: v.CaseSensitiveTS,
					MinTS:           v.MinTS,
					MaxTS:           v.MaxTS,
				},
			})
		}
	}

	return &models.CreateExam{
		Name:        dbExam.Title,
		Description: dbExam.Description,
		Questions:   questions,
	}, nil
}

func (g *GroupController) DeleteExam(userId string, groupId string, examId string) *ErrorResponse {
	member, err := db.DB.GroupMember.FindFirst(
		db.GroupMember.User.Where(
			db.User.ID.Equals(userId),
		),
		db.GroupMember.Group.Where(
			db.Group.ID.Equals(groupId),
		),
	).Exec(ctx)

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return &ErrorResponse{
				Code:    403,
				Message: "You don't have permission to access this group!",
			}
		}
		return &ErrorResponse{
			Code:    500,
			Message: "DB error!",
		}
	}

	if !member.Leader {
		return &ErrorResponse{
			Code:    403,
			Message: "You don't have permission to do this!",
		}
	}

	_, err = db.DB.Exam.FindMany(
		db.Exam.ID.Equals(examId),
		db.Exam.Group.Where(
			db.Group.ID.Equals(groupId),
		),
	).Delete().Exec(ctx)

	if err != nil {
		return &ErrorResponse{
			Code:    500,
			Message: "DB error!",
		}
	}

	return nil
}

func (g *GroupController) DeleteGroup(userId string, groupId string) *ErrorResponse {
	member, err := db.DB.GroupMember.FindFirst(
		db.GroupMember.User.Where(
			db.User.ID.Equals(userId),
		),
		db.GroupMember.Group.Where(
			db.Group.ID.Equals(groupId),
		),
	).Exec(ctx)

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return &ErrorResponse{
				Code:    403,
				Message: "You don't have permission to access this group!",
			}
		}
		log.Println(err.Error())
		return &ErrorResponse{
			Code:    500,
			Message: "DB error!",
		}
	}

	if !member.Leader {
		return &ErrorResponse{
			Code:    403,
			Message: "You don't have permission to do this!",
		}
	}

	group, err := db.DB.Group.FindFirst(
		db.Group.ID.Equals(groupId),
	).With(
		db.Group.Room.Fetch(),
	).Exec(ctx)

	if err != nil {
		log.Println(err.Error())
		return &ErrorResponse{
			Code:    500,
			Message: "DB error!",
		}
	}

	_, err = db.DB.Exam.FindMany(
		db.Exam.Group.Where(
			db.Group.ID.Equals(groupId),
		),
	).Delete().Exec(ctx)

	if err != nil && !errors.Is(err, db.ErrNotFound) {
		log.Println(err.Error())
		return &ErrorResponse{
			Code:    500,
			Message: "DB error!",
		}
	}

	_, err = db.DB.GroupMember.FindMany(
		db.GroupMember.Group.Where(
			db.Group.ID.Equals(groupId),
		),
	).Delete().Exec(ctx)

	if err != nil && !errors.Is(err, db.ErrNotFound) {
		log.Println(err.Error())
		return &ErrorResponse{
			Code:    500,
			Message: "DB error!",
		}
	}

	_, err = db.DB.Group.FindMany(
		db.Group.ID.Equals(groupId),
	).Delete().Exec(ctx)

	if err != nil && !errors.Is(err, db.ErrNotFound) {
		log.Println(err.Error())
		return &ErrorResponse{
			Code:    500,
			Message: "DB error!",
		}
	}

	_, err = db.DB.ChatMember.FindMany(
		db.ChatMember.Room.Where(
			db.ChatRoom.ID.Equals(group.RoomID),
		),
	).Delete().Exec(ctx)

	if err != nil && !errors.Is(err, db.ErrNotFound) {
		log.Println(err.Error())
		return &ErrorResponse{
			Code:    500,
			Message: "DB error!",
		}
	}

	_, err = db.DB.ChatRoom.FindMany(
		db.ChatRoom.ID.Equals(group.RoomID),
	).Delete().Exec(ctx)

	if err != nil && !errors.Is(err, db.ErrNotFound) {
		log.Println(err.Error())
		return &ErrorResponse{
			Code:    500,
			Message: "DB error!",
		}
	}

	return nil
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
