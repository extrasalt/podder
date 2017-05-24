package main

import (
	"testing"
)

func TestCreateNamespace(t *testing.T) {
	ns := CreateNamespace("hala")

	if ns.Meta.Name != "hala" {
		t.Fatal("The Meta doesn't match the given name")
	}

	if ns.Meta.Labels["name"] != "hala" {
		t.Fatal("The Namespace label doesn't match")
	}
}

func TestCreateReplicaSet(t *testing.T) {
	cmdstring := "cmdstr"
	objectName := "objname"
	namespace := "namespace"

	rs := CreateReplicaSet(cmdstring, objectName, namespace)

	if rs.Meta.Name != objectName {
		t.Fatal("Meta name doesn't match object name")
	}

	if rs.Meta.Namespace != namespace {
		t.Fatal("Meta namespace doesn't match with given namespace")
	}

	if rs.Meta.Labels["name"] != objectName {
		t.Fatal("Meta Label doesn't match object name")
	}

	if rs.Spec.Selector.MatchLabels["name"] != objectName {
		t.Fatal("Object name doesn't match in the spec label")
	}

	if rs.Spec.Template.Meta.Labels["name"] != objectName {
		t.Fatal("Template meta label doesn't match object name")
	}

	if rs.Spec.Template.Spec["containers"][0].Name != objectName {
		t.Fatal("container name doesn't match object name")
	}

}

func TestCreateService(t *testing.T) {
	objectName := "objectname"
	namespace := "namespace"

	serv := CreateService(objectName, namespace)

	if serv.Meta.Name != objectName {
		t.Fatal("Meta.Name doesn't match object name")
	}

	if serv.Spec.Selector["name"] != objectName {
		t.Fatal("Selector name doesn't match object name")
	}

}
