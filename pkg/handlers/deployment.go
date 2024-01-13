package handlers

import (
	"fmt"

	"github.com/anodot/anodot-common/pkg/common"
	"github.com/anodot/anodot-common/pkg/events"
	"github.com/anodot/kube-events/pkg/configuration"
	"github.com/anodot/kube-events/pkg/utils"
	apps_v1 "k8s.io/api/apps/v1"
)

type DeploymentHandler struct {
	configuration.EventConfig
}

func (d *DeploymentHandler) SupportedEvent() string {
	return "deployment"
}

func (d *DeploymentHandler) EventData(event Event) ([]events.Event, error) {
	allEvents := make([]events.Event, 0)

	switch event.EventType {
	case "update":
		if event.New == nil || event.Old == nil {
			//TODO better error message
			return nil, fmt.Errorf("unable to retrieve deployment information")
		}

		newDep, ok := event.New.(*apps_v1.Deployment)
		if !ok {
			return nil, fmt.Errorf("%v is not deployment", event)
		}

		oldDeployment, ok := event.Old.(*apps_v1.Deployment)
		if !ok {
			return nil, fmt.Errorf("%v is not deployment", event)
		}

		deploymentName := newDep.Name

		if *newDep.Spec.Replicas != *oldDeployment.Spec.Replicas {
			res := events.Event{
				Title:       fmt.Sprintf("'%s' deployment replica number changed", deploymentName),
				Description: fmt.Sprintf("%s replicas changed from '%d' to '%d'", deploymentName, *oldDeployment.Spec.Replicas, *newDep.Spec.Replicas),
				Category:    d.Category,
				Source:      d.Source,
				Properties: []events.EventProperties{
					{Key: "deployment", Value: deploymentName},
					{Key: "namespace", Value: newDep.Namespace}},
				StartDate: common.AnodotTimestamp{Time: event.EventTime},
			}
			allEvents = append(allEvents, res)
		}

		//image changed
		for _, newC := range newDep.Spec.Template.Spec.Containers {
			for _, oldC := range oldDeployment.Spec.Template.Spec.Containers {

				if newC.Name == oldC.Name {

					if !newC.Resources.Limits.Cpu().Equal(*oldC.Resources.Limits.Cpu()) {

						before := oldC.Resources.Limits.Cpu().AsDec().String()
						after := newC.Resources.Limits.Cpu().AsDec().String()

						res := events.Event{
							Title:       fmt.Sprintf("'%s' deployment cpu limit changed", deploymentName),
							Description: fmt.Sprintf("%s  cpu limit changed from '%s' to '%s'", deploymentName, before, after),
							Category:    d.Category,
							Source:      d.Source,
							Properties: []events.EventProperties{
								{Key: "deployment", Value: deploymentName},
								{Key: "namespace", Value: newDep.Namespace}},
							StartDate: common.AnodotTimestamp{Time: event.EventTime},
						}

						allEvents = append(allEvents, res)
					}

					if !newC.Resources.Requests.Cpu().Equal(*oldC.Resources.Requests.Cpu()) {

						before := oldC.Resources.Requests.Cpu().AsDec().String()
						after := newC.Resources.Requests.Cpu().AsDec().String()

						res := events.Event{
							Title:       fmt.Sprintf("'%s' deployment cpu requests changed", deploymentName),
							Description: fmt.Sprintf("%s  cpu requests changed from '%s' to '%s'", deploymentName, before, after),
							Category:    d.Category,
							Source:      d.Source,
							Properties: []events.EventProperties{
								{Key: "deployment", Value: deploymentName},
								{Key: "namespace", Value: newDep.Namespace}},
							StartDate: common.AnodotTimestamp{Time: event.EventTime},
						}

						allEvents = append(allEvents, res)
					}

					if !newC.Resources.Limits.Memory().Equal(*oldC.Resources.Limits.Memory()) {
						before := oldC.Resources.Limits.Memory().AsDec().String()
						after := newC.Resources.Limits.Memory().AsDec().String()

						res := events.Event{
							Title:       fmt.Sprintf("'%s' deployment memory limits changed", deploymentName),
							Description: fmt.Sprintf("%s  memory limits changed from '%s' to '%s'", deploymentName, before, after),
							Category:    d.Category,
							Source:      d.Source,
							Properties: []events.EventProperties{
								{Key: "deployment", Value: deploymentName},
								{Key: "namespace", Value: newDep.Namespace}},
							StartDate: common.AnodotTimestamp{Time: event.EventTime},
						}

						allEvents = append(allEvents, res)
					}

					if !newC.Resources.Requests.Memory().Equal(*oldC.Resources.Requests.Memory()) {
						before := oldC.Resources.Requests.Memory().AsDec().String()
						after := newC.Resources.Requests.Memory().AsDec().String()

						res := events.Event{
							Title:       fmt.Sprintf("'%s' deployment memory requests changed", deploymentName),
							Description: fmt.Sprintf("%s  memory requests changed from '%s' to '%s'", deploymentName, before, after),
							Category:    d.Category,
							Source:      d.Source,
							Properties: []events.EventProperties{
								{Key: "deployment", Value: deploymentName},
								{Key: "namespace", Value: newDep.Namespace}},
							StartDate: common.AnodotTimestamp{Time: event.EventTime},
						}

						allEvents = append(allEvents, res)
					}

					if newC.Image != oldC.Image {
						res := events.Event{
							Title:       fmt.Sprintf("'%s' deployment container image changed", deploymentName),
							Description: utils.ImageChangedMessage(oldC.Image, newC.Image),
							Category:    d.Category,
							Source:      d.Source,
							Properties: []events.EventProperties{
								{Key: "deployment", Value: deploymentName},
								{Key: "namespace", Value: newDep.Namespace},
								{Key: "container", Value: newC.Name}},
							StartDate: common.AnodotTimestamp{Time: event.EventTime},
						}
						allEvents = append(allEvents, res)
					}
				}
				//TODO env check ?
				/*if !reflect.DeepEqual(c.Env, oldC.Env) {
				}*/
			}
		}

	case "delete":

	}

	return allEvents, nil
}
