package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/getbinary", UserBinaryHandler)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	http.ListenAndServe(":3000", r)
}

func UserBinaryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	binary, _, err := r.FormFile("upload")

	if err != nil {
		panic(err)
	}

	io.Copy(w, binary)
}
