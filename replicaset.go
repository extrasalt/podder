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
	"fmt"
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

func CreateReplicaSet(cmdstr string, objectName string) {
	rs := ReplicaSet{
		Kind:       "ReplicaSet",
		ApiVersion: "extensions/v1beta1",
		Meta: Metadata{
			Name:      objectName,
			Namespace: "default",
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

	endpoint := fmt.Sprintf("/apis/extensions/v1beta1/namespaces/%s/replicasets", "default")

	sendToKube(rs, endpoint)
}
