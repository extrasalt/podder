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
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/minio/minio-go"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	kubehost = "http://" + os.Getenv("KUBERNETES_SERVICE_HOST") + ":" + os.Getenv("KUBERNETES_PORT_443_TCP_PORT")
	// dat, _   = ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")

	// kubetoken = string(dat)
)

func main() {

	var err error
	DB, err = sql.Open("postgres", "password=password  user=user dbname=my_db sslmode=disable")
	if err != nil {
		fmt.Println(err)
	}
	_, err = DB.Exec("CREATE TABLE IF NOT EXISTS login(name varchar, password varchar)")
	if err != nil {
		fmt.Println(err)
	}

	//Test bed values. Replace with real minio address and keys
	endpoint := "play.minio.io:9000"
	accessKeyID := "Q3AM3UQ867SPQQA43P2F"
	secretAccessKey := "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG"
	useSSL := true

	minioClient, err = minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Fatalln(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/getbinary", UserBinaryHandler)
	r.HandleFunc("/services", ListServicesHandler)
	r.HandleFunc("/login", LoginHandler)
	r.HandleFunc("/signup", SignUpHandler).Methods("POST")
	r.HandleFunc("/whoami", WhoAmiHandler)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	http.ListenAndServe(":8000", r)
}

func UserBinaryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	binary, header, err := r.FormFile("upload")

	if err != nil {
		panic(err)
	}

	url, objectName, err := uploadFile(header.Filename, binary)

	cmdstr := createCommandString(url.String(), objectName)
	CreateReplicaSet(cmdstr, objectName)
	CreateService(objectName)

	//TODO: Check errors and return back to the start page if there's a problem

	http.Redirect(w, r, "/services", 302)
}

func createCommandString(url, filename string) string {

	//"wget -O /bin/#{filename} '#url' && chmod +x /bin/{#filename} && {#filename}"

	return fmt.Sprintf("wget -O /bin/%[2]s '%[1]s' && chmod +x /bin/%[2]s && %[2]s", url, filename)

}

func ListServicesHandler(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest("GET", kubehost+"/api/v1/namespaces/default/services", nil)
	if err != nil {
		panic(err)
	}
	// req.Header.Set("Authorization", "Bearer " + kubetoken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)

	// fmt.Printf("%v", string(data))

	if err != nil {
		panic(err)
	}

	servicelist := ServiceList{}
	err = json.Unmarshal(data, &servicelist)

	if err != nil {
		panic(err)
	}

	var responselist []ServiceResponse
	for _, service := range servicelist.Items {
		var resp ServiceResponse
		resp.Name = service.Meta.Name
		resp.Port = service.Spec.Ports[0].NodePort
		responselist = append(responselist, resp)
	}

	fmt.Printf("%+v", responselist)

	tmpl, err := template.ParseFiles("templates/services.html")

	if err != nil {
		panic(err)
	}

	tmpl.Execute(w, responselist)
}

func WhoAmiHandler(w http.ResponseWriter, r *http.Request) {

	cookie, _ := r.Cookie("rcs")

	w.Write([]byte(cookie.Value))

}
