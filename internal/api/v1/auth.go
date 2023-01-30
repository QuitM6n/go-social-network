package v1

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"social_network/internal/api/v1/models"
	"social_network/internal/pkg/jwt"
	DB "social_network/internal/repository"
	"social_network/utils"
)

func Login(wrt http.ResponseWriter, req *http.Request) {
	wrt.Header().Set("Content-Type", "text/html; charset=utf-8")

	if req.Method == http.MethodPost {
		username := req.FormValue("name")
		email := req.FormValue("email")
		password := req.FormValue("password")

		apiUser := models.SignInRequest{}
		apiUser.Username = username
		apiUser.Email = email
		apiUser.Password = password

		ctx := context.Background()
		user, err := DB.GetUserByEmail(ctx, email)
		if err != nil {
			wrt.WriteHeader(http.StatusNotFound)
			log.Println(err, " :User not found")
			utils.ExecTemplate(wrt, "C:/Users/Ruslan/Desktop/go-social-network/static/login.html", "User not found")
			return
		}
		compare := utils.CheckPasswordHash(user.Password, apiUser.Password)

		if !compare {
			wrt.WriteHeader(http.StatusForbidden)
			utils.ExecTemplate(wrt, "C:/Users/Ruslan/Desktop/go-social-network/static/login.html", "password is wrong")
			fmt.Println("Here error")
			return
		}

		validToken, err := jwt.GenerateJWT(apiUser.Email, apiUser.Role)

		if err != nil {
			log.Println(err, ": Invalid Token")
			utils.ExecTemplate(wrt, "C:/Users/Ruslan/Desktop/go-social-network/static/login.html", "Invalid Generate token")
			return
		}

		var token models.SignInResponse
		token.Role = apiUser.Role
		token.Email = apiUser.Email
		token.AccessToken = validToken

		http.Redirect(wrt, req, "/registration", http.StatusSeeOther)
	}

	utils.ExecTemplate(wrt, "C:/Users/Ruslan/Desktop/go-social-network/static/login.html", nil)

}

func Logout(wrt http.ResponseWriter, req *http.Request) {
	wrt.Header().Set("Content-Type", "text/html; charset=utf-8")

}

func SignUp(wrt http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		name := req.FormValue("name")
		email := req.FormValue("email")
		password := req.FormValue("password")
		confirm_pswd := req.FormValue("confirm_pswd") // password to confirm

		fmt.Println("Form Password: ", password)

		if !utils.IsName(name) {
			utils.ExecTemplate(wrt, "C:/Users/Ruslan/Desktop/go-social-network/static/signup.html", "Wrong name entered,or user with so name already exists")
			return
		} else if !utils.IsEmail(email) {
			utils.ExecTemplate(wrt, "C:/Users/Ruslan/Desktop/go-social-network/static/signup.html", "Wrong email entered,or so email already exists")
			return
		} else if !utils.IsPassword(password) {
			utils.ExecTemplate(wrt, "C:/Users/Ruslan/Desktop/go-social-network/static/signup.html", "Wrong password entered")
			return
		}

		hash, err := utils.HashPassword(password)

		if err != nil {
			log.Println(err, ": Failed to hashing password")
			return
		}

		user := models.User{}
		user.Name = name
		user.Password = hash
		user.Email = email
		user.ConfirmPassword = confirm_pswd

		if password == "" || password != user.ConfirmPassword {
			wrt.WriteHeader(http.StatusNotFound)
			utils.ExecTemplate(wrt, "C:/Users/Ruslan/Desktop/go-social-network/static/signup.html", "Password do not match")
			return
		}

		ctx := context.Background()
		id_user, err := DB.CreateUser(ctx, user)

		if err != nil {
			utils.ExecTemplate(wrt, "C:/Users/Ruslan/Desktop/go-social-network/static/signup.html", err)
			return
		}

		log.Printf("User success created: %s", id_user.ID)

		http.Redirect(wrt, req, "/login", http.StatusSeeOther)
	}
	utils.ExecTemplate(wrt, "C:/Users/Ruslan/Desktop/go-social-network/static/signup.html", nil)
}

func ResetPassword(wrt http.ResponseWriter, req *http.Request) {

}
