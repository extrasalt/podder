// Copyright 2017 Mohanarangan Muthukumar

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

type ServiceList struct {
	Items []Service `json:"items"`
}

type ServiceResponse struct {
	Name string
	Port int
}

var minioClient *minio.Client
var err error
var DB *sql.DB

var (
	kubehost = "https://" + os.Getenv("KUBERNETES_SERVICE_HOST") + ":" + os.Getenv("KUBERNETES_PORT_443_TCP_PORT")
	dat, _   = ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")

	kubetoken = string(dat)
)

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
	endpoint := "10.0.0.189:9000"
	accessKeyID := "CK07433U5QNCM9AT6XB5"
	secretAccessKey := "GoPgcppO0D2K3f0ndpj6ILzgAimltBty/Aemwf0B"
	useSSL := false

	minioClient, err = minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Fatalln(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/scale/{name}/{count}", authenticate(ScaleAppHandler))
	r.HandleFunc("/delete/{name}", authenticate(DeleteAppHandler))
	r.HandleFunc("/getbinary", authenticate(UserBinaryHandler))
	r.HandleFunc("/services", authenticate(ListServicesHandler))
	r.HandleFunc("/login", LoginHandler)
	r.HandleFunc("/signup", SignUpHandler).Methods("POST")
	r.HandleFunc("/whoami", authenticate(WhoAmiHandler))
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

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

	url, objectName, err := uploadFile(header.Filename, binary)
	cookie, _ := r.Cookie("rcs")
	ns := cookie.Value

	cmdstr := createCommandString(url.String(), objectName)
	CreateReplicaSet(cmdstr, objectName, ns)
	CreateService(objectName, ns)

	//TODO: Check errors and return back to the start page if there's a problem

	http.Redirect(w, r, "/services", 302)
}

func createCommandString(url, filename string) string {

	//Creates a command string in shell format that resolves down to the following format
	//"wget -O /bin/#{filename} '#url' && chmod +x /bin/{#filename} && {#filename}"

	return fmt.Sprintf("wget -O /bin/%[2]s '%[1]s' && chmod +x /bin/%[2]s && %[2]s", url, filename)
}

func ListServicesHandler(w http.ResponseWriter, r *http.Request) {
	//Shows a list of Services owned by the current user.
	//Gets the current user information from the cookie
	//and GETs from kubernetes api the services pertaining
	//to the user's namespace.

	cookie, _ := r.Cookie("rcs")
	ns := cookie.Value

	endpoint := fmt.Sprintf("/api/v1/namespaces/%s/services", ns)
	req, err := http.NewRequest("GET", kubehost+endpoint, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "Bearer "+kubetoken)

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

	tmpl, err := template.ParseFiles("templates/services.html")
	if err != nil {
		panic(err)
	}
	tmpl.Execute(w, responselist)
}

func WhoAmiHandler(w http.ResponseWriter, r *http.Request) {
	//Reads username from cookie and prints it to Response
	cookie, _ := r.Cookie("rcs")
	w.Write([]byte(cookie.Value))

}

func authenticate(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("rcs")
		if err != nil {
			http.Redirect(w, r, "/login.html", 302)
		} else {
			next(w, r)
		}

	}

}

func DeleteAppHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	name := vars["name"]

	cookie, _ := r.Cookie("rcs")
	deleteApp(name, cookie.Value)

	http.Redirect(w, r, "/services", 302)

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

	http.Redirect(w, r, "/services", 302)
}
