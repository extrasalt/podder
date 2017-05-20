package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	//Gets the username and password from the form data
	//Checks the validity of it. And sets cookie if valid.
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	username := r.Form["name"][0]
	password := r.Form["password"][0]

	err = authorize(username, password)

	if err == nil {
		cookie := &http.Cookie{Name: "rcs", Value: username, MaxAge: 3600, Secure: false, HttpOnly: true, Raw: username}
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/", 302)
	} else {
		w.Write([]byte("Wrong password"))
	}
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {

	//Obtains username and password from the form
	//Generates bcrypt password
	//Inserts the username and hashed password into db
	//Creates a k8s namespace for the new user
	//Redirects the user back to the login page so that
	//he can login.

	err := r.ParseForm()

	if err != nil {
		panic(err)
	}

	username := r.Form["name"][0]
	password := r.Form["password"][0]
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		fmt.Println(err)
	}
	_, err = DB.Exec("insert into login values($1, $2)", username, hashedPassword)
	CreateNamespace(username)
	http.Redirect(w, r, "/login", 302)
}
