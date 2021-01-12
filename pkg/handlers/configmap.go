package handlers

import (
	"fmt"
	"github.com/anodot/anodot-common/pkg/common"
	"github.com/anodot/anodot-common/pkg/events"
	"github.com/anodot/kube-events/pkg/configuration"
	api_v1 "k8s.io/api/core/v1"
)

type ConfigmapHandler struct {
	configuration.EventConfig
}

func (c *ConfigmapHandler) EventData(event Event) ([]events.Event, error) {
	if event.New == nil {
		return nil, fmt.Errorf("unable to retrieve configmap information")
	}

	newCM := event.New.(*api_v1.ConfigMap)
	cmName := newCM.Name

	allEvents := make([]events.Event, 0)
	switch event.EventType {
	case "update":
		res := events.Event{
			Title: fmt.Sprintf("'%s' configmap updated", cmName),
			//Description: fmt.Sprintf("%s replicas changed from '%d' to '%d'", deploymentName, *oldDeployment.Spec.Replicas, *newDep.Spec.Replicas),
			Category: c.Category,
			Source:   c.Source,
			Properties: []events.EventProperties{
				{Key: "configmap", Value: cmName},
				{Key: "namespace", Value: newCM.Namespace}},
			StartDate: common.AnodotTimestamp{Time: event.EventTime},
		}
		allEvents = append(allEvents, res)
	case "delete":

	}

	return allEvents, nil
}

func (c *ConfigmapHandler) SupportedEvent() string {
	return "configmap"
}
