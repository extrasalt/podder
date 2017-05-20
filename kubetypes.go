package main

type ServiceList struct {
	Items []Service `json:"items"`
}

type ReplicaSetList struct {
	Items []ReplicaSet `json:"items"`
}

type ReturnedResult struct {
	Username string
	Items    []ServiceResponse
}

type ServiceResponse struct {
	Name     string
	Port     int
	Replicas int
}

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
