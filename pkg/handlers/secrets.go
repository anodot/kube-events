package handlers

import (
	"fmt"

	"github.com/anodot/anodot-common/pkg/common"
	"github.com/anodot/anodot-common/pkg/events"
	"github.com/anodot/kube-events/pkg/configuration"
	api_v1 "k8s.io/api/core/v1"
)

type SecretHandler struct {
	configuration.EventConfig
}

func (d *SecretHandler) SupportedEvent() string {
	return "secret"
}

func (d *SecretHandler) EventData(event Event) ([]events.Event, error) {
	allEvents := make([]events.Event, 0)
	switch event.EventType {
	case "update":
		//newSecret := event.New.(*api_v1.Secret)
		oldSecret := event.Old.(*api_v1.Secret)
		res := events.Event{
			Title:       fmt.Sprintf("'%s' secret has been changed", oldSecret.Name),
			Description: fmt.Sprintf("%s secret has been changed", oldSecret.Name),
			Category:    d.Category,
			Source:      d.Source,
			Properties: []events.EventProperties{
				{Key: "secret", Value: oldSecret.Name},
				{Key: "namespace", Value: oldSecret.Namespace}},
			StartDate: common.AnodotTimestamp{Time: event.EventTime},
		}
		allEvents = append(allEvents, res)

	case "delete":

	}
	fmt.Println(allEvents)
	return allEvents, nil
}
