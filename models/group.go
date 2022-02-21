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
