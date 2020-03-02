package main

import (
	"kube-events/pkg/configuration"
	"kube-events/pkg/controller"
)

func main() {
	controller.Start(configuration.Configuration{
		ExcludeNamespace: "",
		Resource:         configuration.Resource{Pod: false, Deployment: true},
		Namespace:        "",
	}, nil)
}
