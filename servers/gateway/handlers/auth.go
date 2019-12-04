package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bch0ng/master-debater/servers/gateway/models/users"
	"github.com/dgrijalva/jwt-go"
)

// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func (context *HandlerContext) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if r.Header.Get("Content-Type") == "application/json" {
			decoder := json.NewDecoder(r.Body)
			var newUser users.NewUser
			err := decoder.Decode(&newUser)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			user, err := newUser.ToUser()
			if err != nil {
				fmt.Printf("Error Converting User:%v", err)
				w.WriteHeader(http.StatusBadRequest)
			}
			insertedUser, err := context.Users.Insert(user)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			jsonUser, err := json.Marshal(insertedUser)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusCreated)
			token, exp, err := generateJWT(insertedUser.Username)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			http.SetCookie(w, &http.Cookie{
				Name:     "token",
				Value:    token,
				Expires:  exp,
				HttpOnly: true,
			})
			w.Write([]byte(jsonUser))
		}
	}
}

func (context *HandlerContext) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var creds users.Credentials
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := context.Users.GetByUserName(creds.Username)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = user.Authenticate(creds.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	jsonUser, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	token, exp, err := generateJWT(creds.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  exp,
		Path:     "/",
		HttpOnly: true,
	})
	w.Write([]byte(jsonUser))
}

func (context *HandlerContext) LogoutUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    "",
			Expires:  time.Unix(0, 0),
			Path:     "/",
			HttpOnly: true,
		})
		w.Write([]byte("Successfully logged out."))
	}
}
