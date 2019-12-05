package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"bytes"
	"math/rand"
	"log"
	"os"
	"github.com/INFO441-19au-org/assignments-bch0ng/servers/gateway/sessions"
	"github.com/go-redis/redis"
	"github.com/INFO441-19au-org/assignments-bch0ng/servers/db"
	"encoding/json"
	"github.com/INFO441-19au-org/assignments-bch0ng/servers/gateway/models/users"
)

func GetSessionContext() *SessionContext{

	sessionKey, sessionKeyExists := os.LookupEnv("SESSIONKEY")
	if !sessionKeyExists {
		//log.Fatalf("Environment variable SESSIONKEY not defined.")
		//os.Exit(1)
		sessionKey="asdfslls"
	}
	redisAddr, redisAddrExists := os.LookupEnv("REDISADDR")
	if !redisAddrExists {
		redisAddr=":6379"
	}
	dsn, dsnExists := os.LookupEnv("DSN")
	if !dsnExists {
		dsn=fmt.Sprintf("root:%s@tcp(127.0.0.1:3306)/psql_db", "your-password")
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	redisSession := sessions.NewRedisStore(redisClient, 3600)

	userStore, postgressErr := db.ConnectToPostgres(dsn)
	if postgressErr != nil {
		log.Fatalf("Postgress not working")
		os.Exit(1)
	}

	handlerContext := &SessionContext{
		Key:     sessionKey,
		Session: redisSession,
		User:    userStore,
	}
	return handlerContext
}

func TestUsersHandler(t *testing.T) {
	cases := []struct {
		request			string
		contentType 	string
		newUser 		*users.NewUser
		expectedStatus	int
	}{
		{
			"POST",
			"application/json",
			&users.NewUser{
				Email:        fmt.Sprintf("bchong@uw.edu%v",rand.Intn(10000)),
				Password:     "password",
				PasswordConf: "password",
				UserName:     fmt.Sprintf("BC@uw.edu%v",rand.Intn(10000)),
				FirstName:    "Brandon",
				LastName:     "Chong",
			},
			200,
		},
	}

	for _, c := range cases {
		jsonUser, _ := json.Marshal(c.newUser)
		req, err := http.NewRequest(c.request, "/", bytes.NewBuffer(jsonUser))
		if err != nil {
			t.Fatalf(err.Error())
		}
		req.Header.Add("Content-Type", c.contentType)
		handlerContext:=GetSessionContext();
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlerContext.UsersHandler)

		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v, wanted %v",
				status, http.StatusOK)
		}/*
		expected := `{"alive":true}`
		if rr.Body.String() != expected {
			t.Errorf("Handler returned wrong body: got %v, wanted %v",
				rr.Body.String(), expected)
		}*/
	}
}


func TestInvalidSessionHandler(t *testing.T) {
	cases := []struct {
		request			string
		contentType 	string
		newUser 		*users.NewUser
		expectedStatus	int
	}{
		{
			"POST",
			"application/text",
			&users.NewUser{
				Email:        "bchong@uw.edu",
				Password:     "password",
				PasswordConf: "password",
				UserName:     "bchong",
				FirstName:    "Brandon",
				LastName:     "Chong",
			},
			415,
		},
		{
			"GET",
			"application/text",
			&users.NewUser{
				Email:        fmt.Sprintf("bcsfsong@uw.edu%v",rand.Intn(1000)),
				Password:     "password",
				PasswordConf: "password",
				UserName:     fmt.Sprintf("sss@uw.edu%v",rand.Intn(1000)),
				FirstName:    "Brandon",
				LastName:     "Chong",
			},
			405,
		},
	}

	for _, c := range cases {
		jsonUser, _ := json.Marshal(c.newUser)
		req, err := http.NewRequest(c.request, "/", bytes.NewBuffer(jsonUser))
		if err != nil {
			t.Fatalf(err.Error())
		}
		req.Header.Add("Content-Type", c.contentType)
		handlerContext:=GetSessionContext()
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlerContext.SessionsHandler)
		stringStatus := fmt.Sprintf("%s%d", "Expected Key:", c.expectedStatus)

		handler.ServeHTTP(rr, req)
		t.Log(stringStatus)
		if status := rr.Code; status != c.expectedStatus {
			t.Errorf("Handler returned wrong status code: got %v, wanted %v",
				status, http.StatusOK)
		}
		expected := ``
		if rr.Body.String() != expected {
			t.Errorf("Handler returned wrong body: got %v, wanted %v",
				rr.Body.String(), expected)
		}
	}
}
