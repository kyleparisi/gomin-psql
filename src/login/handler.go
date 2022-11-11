package login

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/ucarion/urlpath"
	"golang.org/x/crypto/bcrypt"
	"gomin/src/app"
	"gomin/src/framework"
	"html/template"
	"io"
	"log"
	"net/mail"
	"os"
	"strings"
)

type Login struct {
	Email    string
	Password string
}

type LoginError struct {
	Email  string `json:"email""`
	Errors struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	} `json:"errors"`
}

func GetHandler(db *sql.DB, session *sessions.Session) func(_ urlpath.Match) framework.Response {
	return func(_ urlpath.Match) framework.Response {
		t, err := template.ParseFiles(os.Getenv("APP_DIR") + "/views/login.gohtml")
		if err != nil {
			panic(err)
		}
		return framework.Response{StatusCode: 200, Template: t}
	}
}

func PostHandler(db *sql.DB, session *sessions.Session) func(_ urlpath.Match, body io.Reader) framework.Response {
	return func(_ urlpath.Match, body io.Reader) framework.Response {
		login := Login{}
		t, err := template.ParseFiles(os.Getenv("APP_DIR") + "/views/login.gohtml")
		if err != nil {
			panic(err)
		}
		err = json.NewDecoder(body).Decode(&login)
		if err != nil {
			panic(err)
		}
		hasEmail := login.Email != ""
		hasPassword := login.Password != ""
		// Input validation
		if !hasEmail || !hasPassword {
			loginError := LoginError{Errors: struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}{Email: "", Password: ""}}
			if !hasEmail {
				loginError.Errors.Email = "Please provide an email address"
			}
			if !hasPassword {
				loginError.Errors.Password = "Please provide a password"
			}
			log.Printf("LoginHandler: %+v", loginError)
			return framework.Response{StatusCode: 400, Data: loginError, Template: t}
		}
		// Email validation
		_, err = mail.ParseAddress(login.Email)
		if err != nil {
			loginError := LoginError{
				Email: login.Email,
				Errors: struct {
					Email    string `json:"email"`
					Password string `json:"password"`
				}{Email: "Not a valid email address", Password: ""}}
			log.Printf("LoginHandler: %+v", loginError)
			return framework.Response{StatusCode: 400, Data: loginError, Template: t}
		}
		// Check for existing users
		var appUser app.AppUser
		err = db.QueryRow("SELECT id, name, email, password FROM app_user where email = $1", strings.ToLower(login.Email)).Scan(&appUser.Id, &appUser.Name, &appUser.Email, &appUser.Password)
		if err != nil {
			panic(err.Error())
		}
		err = bcrypt.CompareHashAndPassword([]byte(appUser.Password), []byte(login.Password))
		if err != nil {
			fmt.Printf("LoginHandler: failed login attempt by: %s", login.Email)
			return framework.Response{StatusCode: 400, Data: struct{ Errors interface{} }{Errors: struct{ Message string }{Message: "Failed to login"}}}
		}

		return framework.Response{Redirect: "/dashboard"}
	}
}
