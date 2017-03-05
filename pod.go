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

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	return nil
}
