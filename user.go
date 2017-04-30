package main

import (
	"golang.org/x/crypto/bcrypt"
	"fmt"
	"crypto/md5"

	"github.com/briand787b/validation"
)


type User struct {
	ID		string
	Email 		string
	HashedPassword 	string
	Username	string
}

const (
	hashCost = 10
	passwordLength = 6
	userIDLength = 16
)

func NewUser(username, email, password string) (User, error) {
	user := User{
		Email:		email,
		Username:	username,
	}
	if username == "" {
		return user, validation.ErrNoUsername
	}

	if email == "" {
		return user, validation.ErrNoEmail
	}

	if password == "" {
		return user, validation.ErrNoPassword
	}

	if len(password) < passwordLength {
		return user, validation.ErrPasswordTooShort
	}

	// Check if the username exists
	existingUser, err := globalUserStore.FindByUsername(username)
	if err != nil {
		return user, err
	}
	if existingUser != nil {
		return user, validation.ErrUsernameExists
	}

	// Check if the email exists
	existingUser, err = globalUserStore.FindByEmail(email)
	if err != nil {
		return user, err
	}
	if existingUser != nil {
		return user, validation.ErrEmailExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)

	user.HashedPassword = string(hashedPassword)
	user.ID = GenerateID("usr", userIDLength)

	return user, err
}

func FindUser(username, password string) (*User, error) {
	out := &User{
		Username: username,
	}

	existingUser, err := globalUserStore.FindByUsername(username)
	if err != nil {
		return out, err
	}
	if existingUser == nil {
		return out, validation.ErrCredentialsIncorrect
	}

	if bcrypt.CompareHashAndPassword(
		[]byte(existingUser.HashedPassword),
		[]byte(password),
	) != nil {
		return out, validation.ErrCredentialsIncorrect
	}

	return existingUser, nil
}

func UpdateUser(user *User, email, currentPassword, newPassword string) (User, error) {
	out := *user
	out.Email = email

	// Check if the email exists
	existingUser, err := globalUserStore.FindByEmail(email)
	if err != nil {
		return out, err
	}
	if existingUser != nil && existingUser.ID != user.ID {
		return out, validation.ErrEmailExists
	}

	// At this point, we can update the email address
	user.Email = email

	// No current password? Don't try to update the password
	if currentPassword == "" {
		return out, nil
	}

	if bcrypt.CompareHashAndPassword(
		[]byte(user.HashedPassword),
		[]byte(currentPassword),
	) != nil {
		return out, validation.ErrPasswordIncorrect
	}

	if newPassword == "" {
		return out, validation.ErrNoPassword
	}

	if len(newPassword) < passwordLength {
		return out, validation.ErrPasswordTooShort
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), hashCost)
	user.HashedPassword = string(hashedPassword)
	return out, err
}

func (user *User) AvatarURL() string {
	return fmt.Sprintf(
		"www.gravatar.com/avatar/%x",
		md5.Sum([]byte(user.Email)),
	)
}

func (user *User) ImagesRoute() string {
	return "/user/" + user.ID
}