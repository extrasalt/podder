package main

type Namespace struct {
	Kind       string        `json:"kind"`
	ApiVersion string        `json:"apiVersion"`
	Meta       NamespaceMeta `json:"metadata"`
}

type NamespaceMeta struct {
	Name   string            `json:"name"`
	Labels map[string]string `json:"labels"`
}

func CreateNamespace(name string) {
	ns := Namespace{
		Kind:       "Namespace",
		ApiVersion: "v1",
		Meta: NamespaceMeta{
			Name: name,
			Labels: map[string]string{
				"name": name,
			},
		},
	}

	endpoint := "/api/v1/namespaces"
	sendToKube(ns, endpoint)
}

// {
//   "kind": "Namespace",
//   "apiVersion": "v1",
//   "metadata": {
//     "name": "development",
//     "labels": {
//       "name": "development"
//     }
//   }
// }
