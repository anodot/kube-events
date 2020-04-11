package handlers

import (
	"fmt"
	"github.com/anodot/anodot-common/pkg/common"
	"github.com/anodot/anodot-common/pkg/events"
	v1 "k8s.io/api/batch/v1"
	log "k8s.io/klog/v2"
	"strings"
)

type JobHandler struct {
	UserEventConfiguration
}

func (d *JobHandler) SupportedEvent() string {
	return "job"
}

func (d *JobHandler) EventData(event Event) ([]events.Event, error) {
	allEvents := make([]events.Event, 0)
	newJob := event.New.(*v1.Job)
	if strings.HasPrefix(newJob.Name, "monitor-offline") {
		log.V(5).Infof("skipping job %q", newJob.Name)
		return allEvents, nil
	}

	switch event.EventType {
	case "create":
		if event.New == nil {
			return nil, fmt.Errorf("unable to retrieve job information")
		}

		newJob := event.New.(*v1.Job)
		jobName := newJob.Name

		allEvents = append(allEvents, events.Event{
			Title:       fmt.Sprintf("'%s' job created", jobName),
			Description: "",
			Category:    d.Category,
			Source:      d.Source,
			Properties: []events.EventProperties{
				{Key: "job", Value: jobName},
				{Key: "namespace", Value: newJob.Namespace}},
			StartDate: common.AnodotTimestamp{Time: newJob.Status.StartTime.Time},
			EndDate:   nil,
		})

	case "update":
		if event.New == nil || event.Old == nil {
			return nil, fmt.Errorf("unable to retrieve job information")
		}

		newJob := event.New.(*v1.Job)
		jobName := newJob.Name

		if newJob.Status.CompletionTime != nil {
			allEvents = append(allEvents, events.Event{
				Title:       fmt.Sprintf("'%s' job completed", jobName),
				Description: "",
				Category:    d.Category,
				Source:      d.Source,
				Properties: []events.EventProperties{
					{Key: "job", Value: jobName},
					{Key: "namespace", Value: newJob.Namespace}},
				StartDate: common.AnodotTimestamp{Time: newJob.Status.CompletionTime.Time},
				EndDate:   nil,
			})
		}

	case "delete":

	}

	return allEvents, nil
}
