package handlers

import (
    "fmt"
    "github.com/anodot/anodot-common/pkg/common"
    "github.com/anodot/anodot-common/pkg/events"
    "github.com/anodot/kube-events/pkg/configuration"
    v1 "k8s.io/api/core/v1"
	// klog "k8s.io/klog/v2"
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

    // Check node 'NotReady' condition
    for _, condition := range node.Status.Conditions {
        if condition.Type == v1.NodeReady && condition.Status != v1.ConditionTrue {
            anodotEvent := events.Event{
                Title:       fmt.Sprintf("Node '%s' is not ready", node.Name),
                Description: fmt.Sprintf("Node %s status is now: NotReady", node.Name),
                Category:    n.Category,
                Source:      n.Source,
                Properties: []events.EventProperties{
                    {Key: "node", Value: node.Name},
                    {Key: "nodegroup", Value: node.Labels["eks.amazonaws.com/nodegroup"]},
                    {Key: "instance-type", Value: node.Labels["beta.kubernetes.io/instance-type"]},
                },
                StartDate: common.AnodotTimestamp{Time: event.EventTime},
            }
            allEvents = append(allEvents, anodotEvent)
            break
        }
    }

	// DiskPressure
	// MemoryPressure
	// PIDPressure

    // other status	

    return allEvents, nil
}


// default              69s         Normal    NodeAllocatableEnforced   node/ip-10-152-22-146.ec2.internal                                    Updated Node Allocatable limit across pods
// default              69s         Normal    Starting                  node/ip-10-152-22-146.ec2.internal                                    Starting kubelet.
// default              69s         Warning   InvalidDiskCapacity       node/ip-10-152-22-146.ec2.internal                                    invalid capacity 0 on image filesystem
// default              69s         Normal    NodeHasSufficientMemory   node/ip-10-152-22-146.ec2.internal                                    Node ip-10-152-22-146.ec2.internal status is now: NodeHasSufficientMemory
// default              69s         Normal    NodeHasNoDiskPressure     node/ip-10-152-22-146.ec2.internal                                    Node ip-10-152-22-146.ec2.internal status is now: NodeHasNoDiskPressure
// default              69s         Normal    NodeHasSufficientPID      node/ip-10-152-22-146.ec2.internal                                    Node ip-10-152-22-146.ec2.internal status is now: NodeHasSufficientPID