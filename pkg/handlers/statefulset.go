package handlers

import (
	"fmt"
	"github.com/anodot/anodot-common/pkg/common"
	"github.com/anodot/anodot-common/pkg/events"
	"github.com/anodot/kube-events/pkg/configuration"
	"github.com/anodot/kube-events/pkg/utils"
	v1 "k8s.io/api/apps/v1"
)

type StatefulSetHandler struct {
	configuration.EventConfig
}

func (s *StatefulSetHandler) SupportedEvent() string {
	return "statefulset"
}

func (s *StatefulSetHandler) EventData(event Event) ([]events.Event, error) {
	allEvents := make([]events.Event, 0)

	switch event.EventType {
	case "update":
		if event.New == nil || event.Old == nil {
			//TODO better error message
			return nil, fmt.Errorf("unable to retrieve deployment information")
		}

		newSts := event.New.(*v1.StatefulSet)
		oldSts := event.Old.(*v1.StatefulSet)

		stsName := newSts.Name

		if *newSts.Spec.Replicas != *oldSts.Spec.Replicas {
			res := events.Event{
				Title:       fmt.Sprintf("'%s' statefulset replica number changed", stsName),
				Description: fmt.Sprintf("%s replicas changed from '%d' to '%d'", stsName, *oldSts.Spec.Replicas, *newSts.Spec.Replicas),
				Category:    s.Category,
				Source:      s.Source,
				Properties: []events.EventProperties{
					{Key: "statefulset", Value: stsName},
					{Key: "namespace", Value: newSts.Namespace}},
				StartDate: common.AnodotTimestamp{Time: event.EventTime},
			}
			allEvents = append(allEvents, res)
		}

		//image changed
		for _, newC := range newSts.Spec.Template.Spec.Containers {
			for _, oldC := range oldSts.Spec.Template.Spec.Containers {
				if newC.Name == oldC.Name {
					if newC.Image != oldC.Image {
						res := events.Event{
							Title:       fmt.Sprintf("'%s' statefulset container image changed", stsName),
							Description: utils.ImageChangedMessage(oldC.Image, newC.Image),
							Category:    s.Category,
							Source:      s.Source,
							Properties: []events.EventProperties{
								{Key: "statefulset", Value: stsName},
								{Key: "namespace", Value: newSts.Namespace},
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
