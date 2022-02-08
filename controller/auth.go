package controller

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"menta-backend/db"
	"menta-backend/mail"
	"menta-backend/models"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

type AuthController struct{}

type ErrorResponse struct {
	Code    int
	Message string
}

var ctx = context.Background()
var VerifyTokenSecret = []byte("BEEF74FAC331925E971CB464BA44F")
var AuthTokenSecret = []byte("BEEF74FAC331925E971CB464BA44F")
var RefreshTokenSecret = []byte("FEF8AF8930EAB2190B3A62A95D943A400DFF9E80A3A7215093E264CF83F0D4CD")

func (a AuthController) RegisterUser(data models.Register) (*db.UserModel, *ErrorResponse) {
	e, err := db.DB.User.FindMany(
		db.User.Or(
			db.User.Email.Equals(data.Email),
			db.User.Username.Equals(data.Username),
		),
	).Exec(ctx)
	if err != nil {
		return nil, &ErrorResponse{
			Code:    500,
			Message: "db error",
		}
	}

	if len(e) > 0 {
		return nil, &ErrorResponse{
			Code:    http.StatusConflict,
			Message: "name or email confilct",
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.MinCost)
	if err != nil {
		return nil, &ErrorResponse{
			Code:    500,
			Message: "Hash error",
		}
	}

	code := calcVerifyCode()

	res, err := db.DB.User.CreateOne(
		db.User.Username.Set(data.Username),
		db.User.Email.Set(data.Email),
		db.User.Password.Set(string(hash)),
		db.User.EmailCode.Set(code),
	).Exec(ctx)

	m := gomail.NewMessage()

	m.SetHeader(`From`, `noreply@zeus.bendi.cf`)
	m.SetHeader(`To`, data.Email)
	m.SetHeader("Subject", `E-Mail cím megerősítése`)
	m.SetBody("text/html", "A megerősítő kód: <b>"+code+"</b>!")

	log.Println(mail.SendMessage(m))

	if err != nil {
		return nil, &ErrorResponse{
			Code:    500,
			Message: "Db error",
		}
	}
	return res, nil
}

func (a AuthController) NeedsValidation(id string) bool {
	res, err := db.DB.User.FindFirst(
		db.User.ID.Equals(id),
	).Exec(ctx)
	if err != nil {
		return false
	}
	return !res.Verified
}

func (a AuthController) CreateVerifyToken(id string) (string, *ErrorResponse) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": id,
	})

	str, err := token.SignedString(VerifyTokenSecret)
	if err != nil {
		return "", &ErrorResponse{
			Code:    500,
			Message: "Error signing token",
		}
	}
	return str, nil
}

func (a AuthController) VerifyEmail(code, id string) *ErrorResponse {
	res, err := db.DB.User.FindFirst(
		db.User.ID.Equals(id),
		db.User.EmailCode.Equals(code),
	).Exec(ctx)
	if err != nil {
		return &ErrorResponse{
			Code:    404,
			Message: "bad code",
		}
	}

	if res.Verified {
		return &ErrorResponse{
			Code:    404,
			Message: "allready verified",
		}
	}

	_, err = db.DB.User.FindMany(
		db.User.ID.Equals(id),
		db.User.EmailCode.Equals(code),
	).Update(
		db.User.Verified.Set(true),
	).Exec(ctx)

	if err != nil {
		return &ErrorResponse{
			Code:    500,
			Message: "db error",
		}
	}

	return nil
}

func (a AuthController) Login(data models.Login) (*db.UserModel, *ErrorResponse) {
	res, err := db.DB.User.FindFirst(
		db.User.Email.Equals(data.Email),
	).Exec(ctx)
	if err != nil {
		return nil, &ErrorResponse{
			Code:    500,
			Message: "db error",
		}
	}

	if bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(data.Password)) != nil {
		return nil, &ErrorResponse{
			Code:    401,
			Message: "wrong password",
		}
	}

	return res, nil
}

func (a AuthController) CreateToken(user *db.UserModel) (string, *ErrorResponse) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"created":  time.Now().Unix(),
		"username": user.Username,
		"type":     "access",
	})

	str, err := token.SignedString(AuthTokenSecret)
	if err != nil {
		return "", &ErrorResponse{
			Code:    500,
			Message: "signing error",
		}
	}

	return str, nil
}

func (a AuthController) CreateRefreshToken(id string) (string, *ErrorResponse) {
	db.DB.RefreshToken.FindMany(
		db.RefreshToken.Owner.Where(
			db.User.ID.Equals(id),
		),
	).Delete().Exec(ctx)

	refresh, err := db.DB.RefreshToken.CreateOne(
		db.RefreshToken.Owner.Link(
			db.User.ID.Equals(id),
		),
	).Exec(ctx)

	if err != nil {
		return "", &ErrorResponse{
			Code:    500,
			Message: "signing error",
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":         id,
		"created":    time.Now().Unix(),
		"type":       "refresh",
		"refresh_id": refresh.ID,
	})

	str, err := token.SignedString(AuthTokenSecret)
	if err != nil {
		return "", &ErrorResponse{
			Code:    500,
			Message: "signing error",
		}
	}

	return str, nil
}

var randomChars = strings.Split("0123456789", "")

const codeLenght = 6

func calcVerifyCode() string {
	tmp := ""
	for x := 0; x < codeLenght; x++ {
		tmp = fmt.Sprintf("%s%s", tmp, getRandomChar())
	}
	return tmp
}

func getRandomChar() string {
	return randomChars[rand.Intn(len(randomChars)-1)]
}
