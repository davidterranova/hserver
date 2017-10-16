// Created by davidterranova on 16/10/2017.

package middlewares

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

type TokenGenerator struct {
	salt          string
	signingMethod jwt.SigningMethod
}

func NewTokenGenerator(salt string, signingMethod jwt.SigningMethod) *TokenGenerator {
	return &TokenGenerator{
		salt:          salt,
		signingMethod: signingMethod,
	}
}

func (tg *TokenGenerator) MountRoutes(r *mux.Router) {
	logrus.WithField("route", "[GET] /token").Info("add route")
	r.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		token := jwt.New(tg.signingMethod)

		claims := token.Claims.(jwt.MapClaims)
		claims["admin"] = true
		claims["user"] = "admin@tellmeplus.com"
		claims["exp"] = time.Now().Add(time.Minute * 60).Unix()

		tokenString, _ := token.SignedString([]byte(tg.salt))

		w.Write([]byte(tokenString))
	}).Methods("GET")
}

func NewJWTMiddleware(jwtSalt string, jwtSigningMethod string) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var signingMethod = jwt.GetSigningMethod(jwtSigningMethod)

	return jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSalt), nil
		},
		SigningMethod: signingMethod,
	}).HandlerWithNext
}
