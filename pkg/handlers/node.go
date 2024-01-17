package handlers

import (
    "fmt"
	"time"
    "github.com/anodot/anodot-common/pkg/common"
    "github.com/anodot/anodot-common/pkg/events"
    "github.com/anodot/kube-events/pkg/configuration"
    v1 "k8s.io/api/core/v1"
)

type NodeHandler struct {
    configuration.EventConfig
}

func (n *NodeHandler) SupportedEvent() string {
    return "node"
}
func (n *NodeHandler) EventData(event Event) ([]events.Event, error) {
    allEvents := make([]events.Event, 0)

    node, ok := event.New.(*v1.Node)
    if !ok {
        return nil, fmt.Errorf("event is not a Node")
    }

	for _, condition := range node.Status.Conditions {
        if condition.Type == v1.NodeMemoryPressure && condition.Status == v1.ConditionTrue {
            anodotEvent := createAnodotEvent("MemoryPressure", node, event.EventTime, n)
            allEvents = append(allEvents, anodotEvent)
        } else if condition.Type == v1.NodeDiskPressure && condition.Status == v1.ConditionTrue {
            anodotEvent := createAnodotEvent("DiskPressure", node, event.EventTime, n)
            allEvents = append(allEvents, anodotEvent)
        }
    }
    return allEvents, nil
}

func createAnodotEvent(conditionType string, node *v1.Node, eventTime time.Time, n *NodeHandler) events.Event {
    return events.Event{
        Title:       fmt.Sprintf("Node '%s' has %s", node.Name, conditionType),
        Description: fmt.Sprintf("Node %s status is now: %s", node.Name, conditionType),
        Category:    n.Category,
        Source:      n.Source,
        Properties: []events.EventProperties{
            {Key: "node", Value: node.Name},
            {Key: "nodegroup", Value: node.Labels["eks.amazonaws.com/nodegroup"]},
            {Key: "instance-type", Value: node.Labels["beta.kubernetes.io/instance-type"]},
        },
        StartDate: common.AnodotTimestamp{Time: eventTime},
    }
}