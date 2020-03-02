package handlers

import "kube-events/pkg/controller"

type KubernetesEventsHandler interface {
	DoHandle(event controller.Event)
}
