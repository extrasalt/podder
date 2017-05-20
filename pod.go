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

func CreatePod(cmdstr string, namespace string) error {

	//Creates a pod by constructing a golang object
	//of the required type signature and marshalls it
	//and sends it to the kubernetes api endpoint with
	//the specified namespace

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

	endpoint := fmt.Sprintf("/api/v1/namespaces/%s/pods", namespace)
	sendToKube(pod, endpoint)
	return nil
}
