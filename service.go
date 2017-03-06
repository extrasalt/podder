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
	"encoding/json"
	"io"
	"net/http"
	"os"
)

type Service struct {
	ApiVersion string            `json:"apiVersion"`
	Kind       string            `json:"kind"`
	Meta       map[string]string `json:"metadata"`
	Spec       ServiceSpec       `json:"spec"`
}

type ServiceSpec struct {
	Selector    map[string]string `json:"selector"`
	ServiceType string            `json:"type"`
	Ports       []ServicePort     `json:"ports"`
}

type ServicePort struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
}

func CreateService() {
	spec := ServiceSpec{
		Selector: map[string]string{
			"name": "binary",
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
		Meta: map[string]string{
			"name": "gooo",
		},
		Spec: spec,
	}

	var b []byte
	reader := bytes.NewBuffer(b)
	encoder := json.NewEncoder(reader)
	encoder.SetEscapeHTML(false)
	encoder.Encode(serv)

	req, err := http.NewRequest("POST", kubehost+"/api/v1/namespaces/default/services", reader)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Authorization", "Bearer " + kubetoken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)

}
