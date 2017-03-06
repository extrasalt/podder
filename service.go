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

	endpoint := "/api/v1/namespaces/default/services"
	sendToKube(serv, endpoint)

}
