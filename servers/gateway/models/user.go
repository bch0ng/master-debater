package users

import (
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

//bcryptCost is the default bcrypt cost to use when hashing passwords
var bcryptCost = 13

//User represents a user account in the database
type User struct {
	ID       int64  `json:"id"`
	PassHash []byte `json:"-"`
	Username string `json:"username"`
}

//Credentials represents user sign-in credentials
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//NewUser represents a new user signing up for an account
type NewUser struct {
	Password     string `json:"password"`
	PasswordConf string `json:"passwordConf"`
	Username     string `json:"username"`
}

//Validate validates the new user and returns an error if
//any of the validation rules fail, or nil if its valid
func (nu *NewUser) Validate() error {
	if len(nu.Password) < 6 {
		return fmt.Errorf("password cannot be less than 6 characters")
	}
	if nu.Password != nu.PasswordConf {
		return fmt.Errorf("password and password confirmation do not match")
	}
	if nu.Username == "" {
		return fmt.Errorf("Username cannot be empty")
	}
	if strings.Contains(" ", nu.Username) {
		return fmt.Errorf("Username cannot have spaces")
	}
	return nil
}

//ToUser converts the NewUser to a User, setting the
//PhotoURL and PassHash fields appropriately
func (nu *NewUser) ToUser() (*User, error) {
	err := nu.Validate()
	if err != nil {
		return nil, err
	}
	user := &User{
		ID:       0,
		Username: nu.Username,
	}
	user.SetPassword(nu.Password)
	return user, nil
}

//SetPassword hashes the password and stores it in the PassHash field
func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return err
	}
	u.PassHash = hash
	return nil
}

//Authenticate compares the plaintext password against the stored hash
//and returns an error if they don't match, or nil if they do
func (u *User) Authenticate(password string) error {
	err := bcrypt.CompareHashAndPassword(u.PassHash, []byte(password))
	if err != nil {
		return fmt.Errorf("given password does not match the hashed password")
	}
	return nil
}
