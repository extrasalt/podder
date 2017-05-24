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
	"strings"
)

type Kube struct {
	Host  string
	Token string
}

func CreateReplicaSet(cmdstr string, objectName string, namespace string) ReplicaSet {
	//Creates a replicaset by constructing a golang object
	//of the required type signature and marshalls it
	//and sends it to the kubernetes api endpoint with
	//the specified namespace

	rs := ReplicaSet{
		Kind:       "ReplicaSet",
		ApiVersion: "extensions/v1beta1",
		Meta: Metadata{
			Name:      objectName,
			Namespace: namespace,
			Labels: map[string]string{
				"name": objectName,
			},
		},
		Spec: ReplicaSpec{
			Replicas: 3,
			Selector: ReplicaSelector{
				MatchLabels: map[string]string{
					"name": objectName,
				},
			},
			Template: ReplicaTemplate{
				Meta: Metadata{
					Labels: map[string]string{
						"name": objectName,
					},
				},
				Spec: map[string][]Container{
					"containers": []Container{
						Container{
							Image:   "extrasalt/wgettu",
							Name:    objectName,
							Command: []string{"sh", "-c", cmdstr},
						},
					},
				},
			},
		},
	}

	//	endpoint := fmt.Sprintf("/apis/extensions/v1beta1/namespaces/%s/replicasets", namespace)
	//	sendToKube(rs, endpoint)

	return rs
}

func (kube *Kube) sendToKube(obj interface{}, endpoint string) {

	//Gets various kubernetes objects, marshalls them and
	//sends them to the kubernetes api

	var method string

	if strings.HasSuffix(endpoint, "scale") {
		method = "PUT"
	} else {
		method = "POST"
	}

	var b []byte
	reader := bytes.NewBuffer(b)
	encoder := json.NewEncoder(reader)
	encoder.SetEscapeHTML(false)
	encoder.Encode(obj)

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: transport}

	req, err := http.NewRequest(method, kube.Host+endpoint, reader)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+kube.Token)

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

	req, err := http.NewRequest("DELETE", kube.Host+endpoint+name, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "Bearer "+kube.Token)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)

	//delete service

	endpoint = fmt.Sprintf("/api/v1/namespaces/%s/services/", namespace)

	req, err = http.NewRequest("DELETE", kube.Host+endpoint+name, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "Bearer "+kube.Token)

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

	req, err := http.NewRequest("GET", kube.Host+endpoint, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "Bearer "+kube.Token)

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
	kube.sendToKube(scale, endpoint)

	return nil
}

func CreateService(objectName string, namespace string) Service {
	//Creates a service by constructing a golang object
	//of the required type signature and marshalls it
	//and sends it to the kubernetes api endpoint with
	//the specified namespace

	//Plants the upload binary objectname in different
	//fields in the go object

	spec := ServiceSpec{
		Selector: map[string]string{
			"name": objectName,
		},
		ServiceType: "NodePort",
		Ports: []ServicePort{
			ServicePort{
				Port:     8000,
				Protocol: "TCP",
			},
		},
	}

	serv := Service{
		ApiVersion: "v1",
		Kind:       "Service",
		Meta: ServiceMetadata{
			Name: objectName,
		},
		Spec: spec,
	}

	//endpoint := fmt.Sprintf("/api/v1/namespaces/%s/services", namespace)
	//kube.sendToKube(serv, endpoint)

	return serv
}

func CreateNamespace(name string) Namespace {
	//Creates the namespace for the given value
	//and sends it to kubernetes api
	//by calling the function

	ns := Namespace{
		Kind:       "Namespace",
		ApiVersion: "v1",
		Meta: NamespaceMeta{
			Name: name,
			Labels: map[string]string{
				"name": name,
			},
		},
	}

	//	endpoint := "/api/v1/namespaces"
	//	kube.sendToKube(ns, endpoint)

	return ns
}

// *Reference*
// {
//   "kind": "Namespace",
//   "apiVersion": "v1",
//   "metadata": {
//     "name": "development",
//     "labels": {
//       "name": "development"
//     }
//   }
// }
