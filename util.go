package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/crypto/bcrypt"
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

func authorize(username string, password string) (autherr error) {
	//Looks up the username and password in the database
	//check its validity
	var dbpassword string
	rows, err := DB.Query("Select password from login where name=$1", username)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		err = rows.Scan(&dbpassword)
		if err != nil {
			panic(err)
		}
		break
	}
	err = bcrypt.CompareHashAndPassword([]byte(dbpassword), []byte(password))
	if err == nil {
		return nil
	} else {
		autherr = fmt.Errorf("Cannot authorize %q", username)
		return autherr
	}
}

func getShortHash(f io.Reader) string {
	//Takes the file content
	//creates a sha256 hash
	//Returns only the first 6 digits.
	hash := sha256.New()
	io.Copy(hash, f)
	key := hex.EncodeToString(hash.Sum(nil))
	return key[:6]
}
