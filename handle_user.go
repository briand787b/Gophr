package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"fmt"
	"github.com/briand787b/validation"
)

func HandleUserNew(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	RenderTemplate(w, r, "users/new", nil)
}

func HandleUserCreate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user, err := NewUser(
		r.FormValue("username"),
		r.FormValue("email"),
		r.FormValue("password"),
	)
	if err != nil {
		if validation.IsValidationError(err) {
			RenderTemplate(w, r, "users/new", map[string]interface{}{
				"Error": err.Error(),
				"User": user,
			})
			return
		}
		panic(err)
	}

	err = globalUserStore.Save(user)
	if err != nil {
		panic(err)
	}

	// Create a new session
	session := NewSession(w)
	session.UserID = user.ID
	fmt.Println("this is the session  ", *session)
	err = globalSessionStore.Save(session)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/?flash=User+created", http.StatusFound)
}

func HandleUserEdit(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user := RequestUser(r)
	RenderTemplate(w, r, "users/edit", map[string]interface{}{
		"User": user,
	})
}

func HandleUserUpdate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	currentUser := RequestUser(r)
	email := r.FormValue("email")
	currentPassword := r.FormValue("currentPassword")
	newPassword := r.FormValue("newPassword")

	user, err := UpdateUser(currentUser, email, currentPassword, newPassword)
	if err != nil {
		if validation.IsValidationError(err) {
			RenderTemplate(w, r, "users/edit", map[string]interface{}{
				"Error": err.Error(),
				"User": user,
			})
			return
		}
		panic(err)
	}

	err = globalUserStore.Save(*currentUser)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/account?flash=User+updated", http.StatusFound)
}

func HandleUserShow(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, err := globalUserStore.Find(params.ByName("userID"))
	if err != nil {
		panic(err)
	}

	// 404
	if user == nil {
		http.NotFound(w, r)
		return
	}

	images, err := globalImageStore.FindAllByUser(user, 0)
	if err != nil {
		panic(err)
	}

	RenderTemplate(w, r, "users/show", map[string]interface{}{
		"User": user,
		"Images": images,
	})
}