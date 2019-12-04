package handlers

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Create the JWT key used to create the signature
var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type JWTMiddleware struct {
	handler http.Handler
}

// ServeHTTP serves HTTP with CORS enabled.
func (jwt *JWTMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	token, exp, err := validateJWT(c.Value)
	if err == nil {
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    token,
			Expires:  exp,
			Path:     "/",
			HttpOnly: true,
		})
	}

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	jwt.handler.ServeHTTP(w, r)
}

// NewJWTMiddleware initializes a new JWTMiddleware struct with the given
// HTTP handler.
func NewJWTMiddleware(handler http.Handler) *JWTMiddleware {
	return &JWTMiddleware{handler}
}

// JWT Middleware
func validateJWT(token string) (string, time.Time, error) {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", time.Now(), err
		}
		return "", time.Now(), err
	}
	if !tkn.Valid {
		return "", time.Now(), errors.New("invalid token")
	}
	token, exp, err := Refresh(claims)
	if err != nil {
		return "", time.Now(), err
	}
	return token, exp, nil
}

func generateJWT(username string) (string, time.Time, error) {
	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		return "", time.Now(), err
	}

	return tokenString, expirationTime, nil
}

func Refresh(claims *Claims) (string, time.Time, error) {
	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		return "", time.Now(), errors.New("token still has > 30 seconds before expiring")
	}

	// Now, create a new token for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", time.Now(), err
	}

	return tokenString, expirationTime, nil
}
