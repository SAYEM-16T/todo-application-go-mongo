package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"todo-app/backend/middleware"
	"todo-app/backend/models"
	"todo-app/backend/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AppPage(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.GetUserID(r)

	ctx, cancel := context.WithTimeout(r.Context(), 6*time.Second)
	defer cancel()

	// fetch user
	var user models.User
	_ = models.UsersColl(DB).FindOne(ctx, bson.M{"_id": uid}).Decode(&user)

	// fetch todos (latest first)
	cur, err := models.TodosColl(DB).Find(ctx, bson.M{"user_id": uid}, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	defer cur.Close(ctx)

	var todos []models.Todo
	for cur.Next(ctx) {
		var t models.Todo
		if err := cur.Decode(&t); err == nil {
			todos = append(todos, t)
		}
	}
	utils.RenderAppPage(w, user.Email, todos)
}

func AddTodo(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.GetUserID(r)
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/app", http.StatusFound)
		return
	}
	title := strings.TrimSpace(r.FormValue("title"))
	if title == "" {
		http.Redirect(w, r, "/app", http.StatusFound)
		return
	}
	t := models.Todo{
		ID:        primitive.NewObjectID(),
		UserID:    uid,
		Title:     title,
		Done:      false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	_, _ = models.TodosColl(DB).InsertOne(ctx, t)
	http.Redirect(w, r, "/app", http.StatusFound)
}

func ToggleTodo(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.GetUserID(r)
	idStr := chi.URLParam(r, "id")
	oid, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Redirect(w, r, "/app", http.StatusFound)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var t models.Todo
	err = models.TodosColl(DB).FindOne(ctx, bson.M{"_id": oid, "user_id": uid}).Decode(&t)
	if err == nil {
		_, _ = models.TodosColl(DB).UpdateOne(ctx, bson.M{"_id": oid, "user_id": uid}, bson.M{
			"$set": bson.M{"done": !t.Done, "updated_at": time.Now()},
		})
	}
	http.Redirect(w, r, "/app", http.StatusFound)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.GetUserID(r)
	idStr := chi.URLParam(r, "id")
	oid, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Redirect(w, r, "/app", http.StatusFound)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	_, _ = models.TodosColl(DB).DeleteOne(ctx, bson.M{"_id": oid, "user_id": uid})
	http.Redirect(w, r, "/app", http.StatusFound)
}
