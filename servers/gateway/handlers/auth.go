package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bch0ng/master-debater/servers/gateway/models/users"
	"github.com/dgrijalva/jwt-go"
)

// Create the JWT key used to create the signature
var jwtKey = []byte("my_secret_key")

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
			generateJWT(w, insertedUser.Username)
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
	generateJWT(w, creds.Username)
	w.Write([]byte(jsonUser))
}

/*
func (context *HandlerContext) LogoutUserHandler(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("token")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	// Put JWT token into redis blacklist
}
*/
