package handlers

import (
	"context"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"todo-app/backend/middleware"
	"todo-app/backend/models"
	"todo-app/backend/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var DB *mongo.Database
var FrontendDir string

func Init(db *mongo.Database, frontendDir string) {
	DB = db
	FrontendDir = frontendDir
}

func Home(w http.ResponseWriter, r *http.Request) {
	if _, ok := middleware.GetUserID(r); ok {
		http.Redirect(w, r, "/app", http.StatusFound)
		return
	}
	http.ServeFile(w, r, filepath.Join(FrontendDir, "index.html"))
}

func ServeLoginPage(w http.ResponseWriter, r *http.Request) {
	if _, ok := middleware.GetUserID(r); ok {
		http.Redirect(w, r, "/app", http.StatusFound)
		return
	}
	http.ServeFile(w, r, filepath.Join(FrontendDir, "login.html"))
}

func ServeRegisterPage(w http.ResponseWriter, r *http.Request) {
	if _, ok := middleware.GetUserID(r); ok {
		http.Redirect(w, r, "/app", http.StatusFound)
		return
	}
	http.ServeFile(w, r, filepath.Join(FrontendDir, "register.html"))
}

func RegisterPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad form", http.StatusBadRequest)
		return
	}
	email := strings.ToLower(strings.TrimSpace(r.FormValue("email")))
	pw := r.FormValue("password")
	if len(email) < 5 || len(pw) < 6 {
		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}
	hash, err := utils.HashPassword(pw)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	u := models.User{
		ID:           primitive.NewObjectID(),
		Email:        email,
		PasswordHash: hash,
		CreatedAt:    time.Now(),
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	_, err = models.UsersColl(DB).InsertOne(ctx, u)
	if err != nil {
		http.Redirect(w, r, "/register?exists=1", http.StatusFound)
		return
	}
	token, err := utils.NewSessionToken(u.ID)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	utils.SetSessionCookie(w, token)
	http.Redirect(w, r, "/app", http.StatusFound)
}

func LoginPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad form", http.StatusBadRequest)
		return
	}
	email := strings.ToLower(strings.TrimSpace(r.FormValue("email")))
	pw := r.FormValue("password")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var u models.User
	err := models.UsersColl(DB).FindOne(ctx, bson.M{"email": email}).Decode(&u)
	if err != nil || !utils.CheckPassword(u.PasswordHash, pw) {
		http.Redirect(w, r, "/login?fail=1", http.StatusFound)
		return
	}
	token, err := utils.NewSessionToken(u.ID)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	utils.SetSessionCookie(w, token)
	http.Redirect(w, r, "/app", http.StatusFound)
}

func LogoutPost(w http.ResponseWriter, r *http.Request) {
	utils.ClearSessionCookie(w)
	http.Redirect(w, r, "/", http.StatusFound)
}
