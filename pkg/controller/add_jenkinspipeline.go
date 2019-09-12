package controller

import (
	"github.com/jakubbujny/jenkins-pipeline-operator/pkg/controller/jenkinspipeline"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, jenkinspipeline.Add)
}
