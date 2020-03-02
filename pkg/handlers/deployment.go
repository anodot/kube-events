package handlers

import (
	apps_v1beta1 "k8s.io/api/apps/v1beta1"
	"kube-events/pkg/controller"
	"reflect"
)

type DeploymentHandler struct {
}

// 1. replica number changed
// 2. Env variable changed
// 3. Image changed ?

func (d *DeploymentHandler) DoHandle(event controller.Event) {
	panic("implement me")
	newDep := event.New.(*apps_v1beta1.Deployment)
	oldDeployment := event.Old.(*apps_v1beta1.Deployment)

	if newDep.Spec.Replicas != oldDeployment.Spec.Replicas {
		//create event replicas changed
	}

	for _, c := range newDep.Spec.Template.Spec.Containers {
		for _, oldC := range oldDeployment.Spec.Template.Spec.Containers {
			if c.Name == oldC.Name {
				if c.Image != oldC.Image {
					//image changed event
				}
			}

			if !reflect.DeepEqual(c.Env, oldC.Env) {

			}
		}
	}

}
