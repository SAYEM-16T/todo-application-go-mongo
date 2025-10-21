package utils

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Claims struct {
	UID string `json:"uid"`
	jwt.RegisteredClaims
}

func jwtKey() []byte {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		s = "dev-insecure-secret"
	}
	return []byte(s)
}

func NewSessionToken(userID primitive.ObjectID) (string, error) {
	claims := &Claims{
		UID: userID.Hex(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey())
}

func ParseSessionToken(tokenStr string) (primitive.ObjectID, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtKey(), nil
	})
	if err != nil {
		return primitive.NilObjectID, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		oid, err := primitive.ObjectIDFromHex(claims.UID)
		if err != nil {
			return primitive.NilObjectID, err
		}
		return oid, nil
	}
	return primitive.NilObjectID, errors.New("invalid token")
}

func SetSessionCookie(w http.ResponseWriter, token string) {
	c := &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	}
	http.SetCookie(w, c)
}

func ClearSessionCookie(w http.ResponseWriter) {
	c := &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	}
	http.SetCookie(w, c)
}
