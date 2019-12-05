package handlers

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/INFO441-19au-org/assignments-bch0ng/servers/gateway/models/users"
	"github.com/INFO441-19au-org/assignments-bch0ng/servers/gateway/sessions"
)

// UsersHandler handles all requests for general user actions
func (context *SessionContext) UsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		contentType := r.Header.Get("Content-Type")
		fmt.Println("User request is a post request, content type is:",contentType)
		if contentType == "application/json" {
			decoder := json.NewDecoder(r.Body)
			var t users.NewUser
			err := decoder.Decode(&t)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				user,err:=t.ToUser()
				if err !=nil{
					fmt.Printf("Error Converting User:%v",err)
					w.WriteHeader(http.StatusBadRequest)
				}
				newState := &SessionState{
					StartTime: time.Now(),
					User:      user,
				}
				//fmt.Printf("Sanity check, user:%v",(*user).Email)
				_user, idErr := context.User.Insert(user)
				if idErr != nil {
					w.Write([]byte("Insert Error"))
				} else {
					_user.PassHash = nil
					_user.Email = ""
				}
				newSession, sessionErr := sessions.BeginSession(context.Key, context.Session, newState, w)
				fmt.Println("New Session:",newSession)
				if sessionErr != nil {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Session Error"))
				} else {
					encodedUser, encodingErr := json.Marshal(_user)
					if encodingErr != nil {
						w.WriteHeader(http.StatusUnauthorized)
						w.Write([]byte("Invalid Auth"))
					} else {
						fmt.Println("Writing encoded user:",encodedUser)
						w.WriteHeader(http.StatusCreated)
						w.Header().Add("Content-Type", "application/json")
						w.Write(encodedUser)
					}
				}
			}
		} else {
			w.WriteHeader(http.StatusUnsupportedMediaType)
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// SpecificUserHandler handles all requests for specific user actions
func (context *SessionContext) SpecificUserHandler(w http.ResponseWriter, r *http.Request) {
	var creds users.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	user, _ := context.User.GetByEmail(creds.Email)
	authErr := user.Authenticate(creds.Password)
	if authErr != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if r.Method == http.MethodGet {
		userProfile, err := context.User.GetByID(user.ID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("User Not Found!"))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Header().Add("Content-Type", "application/json")
			encodedUser, _ := json.Marshal(userProfile)
			w.Write(encodedUser)
		}
	} else if r.Method == http.MethodPatch {
		uriSegments := strings.Split(r.URL.Path, "/")
		lastSegments := uriSegments[len(uriSegments)-1]
		if lastSegments != "me" || lastSegments != strconv.FormatInt(user.ID, 10) {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("No Such User!"))
			return
		}
		contentType := r.Header.Get("Content-type")
		if contentType != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte("The request body must be in JSON!"))
		} else {
			var update users.Updates
			err := json.NewDecoder(r.Body).Decode(&update)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				user.ApplyUpdates(&update)
				updatedUser, err := context.User.Update(user.ID, &update)
				if err != nil {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte("User could not be updated."))
				}
				w.WriteHeader(http.StatusOK)
				w.Header().Add("Content-Type", "application/json")
				encodedUser, _ := json.Marshal(updatedUser)
				w.Write(encodedUser)
			}
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// SessionsHandler handles sessions by allowing clients to being a new session
// using their user credentials.
func (context *SessionContext) SessionsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		debugStr := "Sessions request is a post request, Content Type is:" + r.Header.Get("Content-Type")
		fmt.Printf("%v\n", debugStr)
		contentType := r.Header.Get("Content-Type")
		if contentType == "application/json" {
			fmt.Println("Sessions request is a application/json")
			decoder := json.NewDecoder(r.Body)
			var t users.Credentials
			err := decoder.Decode(&t)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				//find user,
				//if not found or authentication failed,
				//return http.StatusUnauthorized with message 'invalid credentials'
				//
				//if authentication is successful make a new session.
				//respond with http.StatusCreated (201), set Content-Type header to application/json
				//a copy of the user's profile in the response body, encoded as a JSON object.
				user, userError := context.User.GetByEmail(t.Email)
				ip, port, err := net.SplitHostPort(r.RemoteAddr)
				if err != nil {
					fmt.Fprintf(w, "userip: %q is not IP:port%v", r.RemoteAddr, port)
				}

				userIP := net.ParseIP(ip)
				if userIP == nil {
					fmt.Fprintf(w, "userip: %q is not IP:port", r.RemoteAddr)
				}
				context.User.AddToLog(t.Email, userIP.String())
				if userError != nil {
					time.Sleep(1 * time.Second)
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("invalid credentials"))
				} else {
					//If you find a matching profile, hash the supplied password using
					// the same algorithm and parameters as you used when generating the password
					// hash during sign-up.
					//Compare the two hashes: if they match, the password was correct
					//and you should start a new authenticated session; if not, the password
					// was invalid, and you should respond with an Unauthorized (401) status code,
					// and the same vague message used when you can't find the user profile.
					authErr := user.Authenticate(t.Password)
					if authErr != nil {
						time.Sleep(1 * time.Second)
						w.WriteHeader(http.StatusUnauthorized)
						w.Write([]byte("Invalid Auth"))
					} else {
						newState := &SessionState{
							StartTime: time.Now(),
							User:      user,
						}
						newSession, sessionErr := sessions.BeginSession(context.Key, context.Session, newState, w)
						fmt.Println(newSession)
						if sessionErr != nil {
							w.WriteHeader(http.StatusUnauthorized)
							w.Write([]byte("Session Error"))
						} else {
							encodedUser, encodingErr := json.Marshal(user)
							if encodingErr != nil {
								w.WriteHeader(http.StatusUnauthorized)
								w.Write([]byte("Invalid Auth"))
							} else {
								w.WriteHeader(http.StatusCreated)
								w.Header().Add("Content-Type", "application/json")
								w.Write(encodedUser)
							}
						}
					}
				}
			}
		} else {
			fmt.Println("Sessions request not a application/json")
			w.WriteHeader(http.StatusUnsupportedMediaType)
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// SpecificSessionHandler handles closing a specific authenticated sessions.
func (context *SessionContext) SpecificSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		uriSegments := strings.Split(r.URL.Path, "/")
		if uriSegments[len(uriSegments)-1] != "mine" {
			w.WriteHeader(http.StatusForbidden)
		} else {
			w.Write([]byte("Signed out!"))
			return
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
