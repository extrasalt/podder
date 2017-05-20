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

func CreateReplicaSet(cmdstr string, objectName string, namespace string) {
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

	endpoint := fmt.Sprintf("/apis/extensions/v1beta1/namespaces/%s/replicasets", namespace)
	sendToKube(rs, endpoint)
}
