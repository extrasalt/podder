// Copyright 2017 Mohanarangan Muthukumar

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License..

package main

import (
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/minio/minio-go"
)

var store *Store
var err error
var DB *sql.DB

var (
	dat, _ = ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
)
var kube = &Kube{
	Host:  "https://" + os.Getenv("KUBERNETES_SERVICE_HOST") + ":" + os.Getenv("KUBERNETES_PORT_443_TCP_PORT"),
	Token: string(dat),
}

func main() {

	var err error
	DB, err = sql.Open("postgres", "postgres://user:password@10.0.0.86/my_db?sslmode=disable")
	//Hardcoded IP. Fix DNS and Replace this.

	if err != nil {
		fmt.Println(err)
	}
	_, err = DB.Exec("CREATE TABLE IF NOT EXISTS login(name varchar, password varchar)")
	if err != nil {
		fmt.Println(err)
	}

	//HARDCODED values. Replace after fixing DNS
	endpoint := "10.0.0.121:9000"
	accessKeyID := "A3CS41BWB9J37FAZGTPT"
	secretAccessKey := "mPtRh7OvMxkDZYpJ63eWHGdemlSbk7pQ6kFl0kmP"
	useSSL := false

	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Fatalln(err)
	}

	store = &Store{
		Client: minioClient,
	}

	r := mux.NewRouter()
	r.HandleFunc("/scale/{name}/{count}", authenticate(ScaleAppHandler))
	r.HandleFunc("/delete/{name}", authenticate(DeleteAppHandler))
	r.HandleFunc("/getbinary", authenticate(UserBinaryHandler))
	r.HandleFunc("/", authenticate(ListServicesHandler))
	r.HandleFunc("/login", ShowLoginPageHandler).Methods("GET")
	r.HandleFunc("/login", LoginHandler).Methods("POST")
	r.HandleFunc("/signup", SignUpHandler).Methods("POST")

	http.ListenAndServe(":8000", r)
}

func UserBinaryHandler(w http.ResponseWriter, r *http.Request) {
	//Gets the [golang] binary from the user
	//Uploads it to a bucket and get a public url
	//Creates command string from the object name and url
	//Creates replicasets and kubernetes service

	r.ParseMultipartForm(32 << 20)
	binary, header, err := r.FormFile("upload")
	if err != nil {
		panic(err)
	}
	url, objectName, err := store.Upload(header.Filename, binary)
	cookie, _ := r.Cookie("rcs")
	ns := cookie.Value

	cmdstr := createCommandString(url.String(), objectName)
	CreateReplicaSet(cmdstr, objectName, ns)
	CreateService(objectName, ns)

	//TODO: Check errors and return back to the start page if there's a problem
	http.Redirect(w, r, "/", 302)
}

func ListServicesHandler(w http.ResponseWriter, r *http.Request) {
	//Shows a list of Services owned by the current user.
	//Gets the current user information from the cookie
	//and GETs from kubernetes api the services pertaining
	//to the user's namespace.
	cookie, _ := r.Cookie("rcs")
	ns := cookie.Value

	//Get Services
	endpoint := fmt.Sprintf("/api/v1/namespaces/%s/services", ns)
	req, err := http.NewRequest("GET", kube.Host+endpoint, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "Bearer "+kube.Token)
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transport}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	servicelist := ServiceList{}
	err = json.Unmarshal(data, &servicelist)

	if err != nil {
		panic(err)
	}

	//Create a separate response list object
	//to make templating easier.
	var responselist []ServiceResponse
	for _, service := range servicelist.Items {
		var resp ServiceResponse
		resp.Name = service.Meta.Name
		resp.Port = service.Spec.Ports[0].NodePort
		responselist = append(responselist, resp)
	}
	result := ReturnedResult{
		Username: ns,
		Items:    responselist,
	}
	tmpl, err := template.ParseFiles("templates/services.html")
	if err != nil {
		panic(err)
	}
	tmpl.Execute(w, result)
}

func DeleteAppHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	cookie, _ := r.Cookie("rcs")
	deleteApp(name, cookie.Value)
	http.Redirect(w, r, "/", 302)
}

func ScaleAppHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	countString := vars["count"]
	count, err := strconv.Atoi(countString)
	if err != nil {
		panic(err)
	}

	cookie, _ := r.Cookie("rcs")
	namespace := cookie.Value
	scaleApp(namespace, name, count)
	http.Redirect(w, r, "/", 302)
}

func ShowLoginPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/login.html")
	tmpl.Execute(w, nil)
}
