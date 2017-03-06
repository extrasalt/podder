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

type Metadata struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Labels    map[string]string `json:"labels"`
}

type ReplicaSet struct {
	Kind       string      `json:"kind"`
	ApiVersion string      `json:"apiVersion"`
	Meta       Metadata    `json:"metadata"`
	Spec       ReplicaSpec `json:"spec"`
}

type ReplicaSpec struct {
	Replicas int             `json:"replicas"`
	Selector ReplicaSelector `json:"selector"`
	Template ReplicaTemplate `json:"template"`
}

type ReplicaSelector struct {
	MatchLabels map[string]string `json:"matchLabels"`
}

type ReplicaTemplate struct {
	Meta Metadata               `json:"metadata"`
	Spec map[string][]Container `json:"spec"`
}

func CreateReplicaSet(cmdstr string) {
	rs := ReplicaSet{
		Kind:       "ReplicaSet",
		ApiVersion: "extensions/v1beta1",
		Meta: Metadata{
			Name:      "goo",
			Namespace: "default",
			Labels: map[string]string{
				"name": "binary",
			},
		},
		Spec: ReplicaSpec{
			Replicas: 3,
			Selector: ReplicaSelector{
				MatchLabels: map[string]string{
					"name": "binary",
				},
			},
			Template: ReplicaTemplate{
				Meta: Metadata{
					Labels: map[string]string{
						"name": "binary",
					},
				},
				Spec: map[string][]Container{
					"containers": []Container{
						Container{
							Image:   "extrasalt/wgettu",
							Name:    "binary",
							Command: []string{"sh", "-c", cmdstr},
						},
					},
				},
			},
		},
	}

	var b []byte
	reader := bytes.NewBuffer(b)
	encoder := json.NewEncoder(reader)
	encoder.SetEscapeHTML(false)
	encoder.Encode(rs)

	req, err := http.NewRequest("POST", kubehost+"/apis/extensions/v1beta1/namespaces/default/replicasets", reader)
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
}
