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
	"net/http"
)

type Container struct {
	Image   string           `json:"image"`
	Name    string           `json:"name"`
	Command []string         `json:"command"`
	Ports   []map[string]int `json:"ports"`
}

type Pod struct {
	Kind       string                 `json:"kind"`
	ApiVersion string                 `json:"apiVersion"`
	Metadata   map[string]string      `json:"metadata"`
	Spec       map[string][]Container `json:"spec"`
}

func CreatePod(cmdstr string) error {

	ports := []map[string]int{
		map[string]int{
			"hostPort":      8000,
			"containerPort": 8000,
		},
	}
	container := Container{"extrasalt/wgettu", "binary", []string{"sh", "-c", cmdstr}, ports}
	metadata := map[string]string{
		"name":      "goo",
		"namespace": "default",
	}
	pod := Pod{"Pod", "v1", metadata,
		map[string][]Container{"containers": []Container{container}}}

	var b []byte
	reader := bytes.NewBuffer(b)
	encoder := json.NewEncoder(reader)
	encoder.SetEscapeHTML(false)
	encoder.Encode(pod)

	req, err := http.NewRequest("POST", kubehost+"/api/v1/namespaces/default/pods", reader)
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

	return nil
}
