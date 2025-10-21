package middleware

import (
	"context"
	"net/http"
	"todo-app/backend/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type contextKey string

const userIDKey contextKey = "userID"

func AuthRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session_token")
		if err != nil || c.Value == "" {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		uid, err := utils.ParseSessionToken(c.Value)
		if err != nil {
			utils.ClearSessionCookie(w)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, uid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(r *http.Request) (primitive.ObjectID, bool) {
	v := r.Context().Value(userIDKey)
	if v == nil {
		// Fallback: try cookie (used on home route)
		c, err := r.Cookie("session_token")
		if err == nil && c.Value != "" {
			uid, err2 := utils.ParseSessionToken(c.Value)
			if err2 == nil {
				return uid, true
			}
		}
		return primitive.NilObjectID, false
	}
	uid, ok := v.(primitive.ObjectID)
	return uid, ok
}
