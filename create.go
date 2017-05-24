package main

func CreateReplicaSet(cmdstr string, objectName string, namespace string) ReplicaSet {
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

	//	endpoint := fmt.Sprintf("/apis/extensions/v1beta1/namespaces/%s/replicasets", namespace)
	//	sendToKube(rs, endpoint)

	return rs
}

func CreateService(objectName string, namespace string) Service {
	//Creates a service by constructing a golang object
	//of the required type signature and marshalls it
	//and sends it to the kubernetes api endpoint with
	//the specified namespace

	//Plants the upload binary objectname in different
	//fields in the go object

	spec := ServiceSpec{
		Selector: map[string]string{
			"name": objectName,
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
		Meta: ServiceMetadata{
			Name: objectName,
		},
		Spec: spec,
	}

	//endpoint := fmt.Sprintf("/api/v1/namespaces/%s/services", namespace)
	//kube.sendToKube(serv, endpoint)

	return serv
}

func CreateNamespace(name string) Namespace {
	//Creates the namespace for the given value
	//and sends it to kubernetes api
	//by calling the function

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

	//	endpoint := "/api/v1/namespaces"
	//	kube.sendToKube(ns, endpoint)

	return ns
}

// *Reference*
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
