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

type Scale struct {
	ApiVersion string    `json:"apiVersion,omitempty"`
	Kind       string    `json:"kind,omitempty"`
	Metadata   Metadata  `json:"metadata"`
	Spec       ScaleSpec `json:"spec,omitempty"`
}

type ScaleSpec struct {
	Replicas int64 `json:"replicas,omitempty"`
}

type Namespace struct {
	Kind       string        `json:"kind"`
	ApiVersion string        `json:"apiVersion"`
	Meta       NamespaceMeta `json:"metadata"`
}

type NamespaceMeta struct {
	Name   string            `json:"name"`
	Labels map[string]string `json:"labels"`
}

type Service struct {
	ApiVersion string          `json:"apiVersion"`
	Kind       string          `json:"kind"`
	Meta       ServiceMetadata `json:"metadata"`
	Spec       ServiceSpec     `json:"spec"`
}

type ServiceSpec struct {
	Selector    map[string]string `json:"selector"`
	ServiceType string            `json:"type"`
	Ports       []ServicePort     `json:"ports"`
}

type ServicePort struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	NodePort int    `json:"nodePort"`
}

type ServiceMetadata struct {
	Name string `json:"name"`
}
