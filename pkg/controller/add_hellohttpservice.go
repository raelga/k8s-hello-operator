package controller

import (
	"github.com/raelga/k8s-hello-operator/pkg/controller/hellohttpservice"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, hellohttpservice.Add)
}
