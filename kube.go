// Copyright 2017 Mohanarangan Muthukumar

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//    http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func sendToKube(obj interface{}, endpoint string) {

	//Gets various kubernetes objects, marshalls them and
	//sends them to the kubernetes api

	var b []byte
	reader := bytes.NewBuffer(b)
	encoder := json.NewEncoder(reader)
	encoder.SetEscapeHTML(false)
	encoder.Encode(obj)

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: transport}

	req, err := http.NewRequest("POST", kubehost+endpoint, reader)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+kubetoken)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)

}

func deleteApp(name, namespace string) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: transport}

	//TODO: Scale replicaset to 0

	//delete replicaset
	endpoint := fmt.Sprintf("/apis/extensions/v1beta1/namespaces/%s/replicasets/", namespace)

	req, err := http.NewRequest("DELETE", kubehost+endpoint+name, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "Bearer "+kubetoken)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)

	//delete service

	endpoint = fmt.Sprintf("/api/v1/namespaces/%s/services/", namespace)

	req, err = http.NewRequest("DELETE", kubehost+endpoint+name, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "Bearer "+kubetoken)

	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)
}
