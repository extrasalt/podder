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
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Scale struct {
	ApiVersion string    `json:"apiVersion,omitempty"`
	Kind       string    `json:"kind,omitempty"`
	Metadata   Metadata  `json:"metadata"`
	Spec       ScaleSpec `json:"spec,omitempty"`
}

type ScaleSpec struct {
	Replicas int64 `json:"replicas,omitempty"`
}

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

	//Remove all pods
	scaleApp(namespace, name, 0)

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

func getScale(namespace, name string) (*Scale, error) {

	var scale Scale
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transport}
	endpoint := fmt.Sprintf("/apis/extensions/v1beta1/namespaces/%s/replicasets/%s/scale", namespace, name)

	req, err := http.NewRequest("GET", kubehost+endpoint, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "Bearer "+kubetoken)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("Get scale error non 200 reponse: " + resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&scale)
	if err != nil {
		return nil, err
	}
	fmt.Printf("SCALE: %v", scale)

	return &scale, nil
}

func scaleApp(namespace, name string, count int) error {
	scale, err := getScale(namespace, name)
	if err != nil {
		return err
	}
	scale.Spec.Replicas = int64(count)

	endpoint := fmt.Sprintf("/apis/extensions/v1beta1/namespaces/%s/replicasets/%s/scale", namespace, name)
	sendToKube(scale, endpoint)

	return nil
}
