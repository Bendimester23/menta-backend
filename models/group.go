package models

import "errors"

type CreateGroup struct {
	Name            string `json:"name" validate:"required,min=6,max=22"`
	LimitJoins      bool   `json:"needs_approval"`
	JoinWithCode    bool   `json:"join_with_code"`
	ConnectToSchool bool   `json:"connect_to_school"`
	SchoolCode      string `json:"school_code" validate:"max=10"`
	CodePrefix      string `json:"code_prefix" validate:"max=5"`
}

type CreateExam struct {
	Name        string               `json:"name" validate:"required,min=6,max=22"`
	Description string               `json:"description" validate:"min=2,max=300"`
	Questions   []CreateExamQuestion `json:"questions" validate:"required"`
}

type CreateExamQuestion struct {
	Id   int          `json:"id" validate:"required"`
	Data QuestionData `json:"data" validate:"required"`
}

type QuestionData struct {
	Name            string   `json:"name" validate:"required,min=5,max=50"`
	Type            string   `json:"type" validate:"required"`
	AutoGrade       bool     `json:"auto_grade"`
	CaseSensitiveTS bool     `json:"case_sensitive"`
	AnswerTF        bool     `json:"answer"`
	AnswersTS       []string `json:"answers"`
	MinTS           int      `json:"min"`
	MaxTS           int      `json:"max"`
}

func (r *CreateGroup) Validate() error {
	if r.ConnectToSchool {
		if len(r.SchoolCode) < 6 {
			return errors.New("school code is shorter than expected")
		}
	}
	if r.JoinWithCode {
		if len(r.CodePrefix) < 2 {
			return errors.New("school code is shorter than expected")
		}
	}
	return validation.Struct(r)
}

func (c *CreateExam) Validate() error {
	return validation.Struct(c)
}
