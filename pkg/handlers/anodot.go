package handlers

import (
	"strings"
	"time"

	"github.com/anodot/anodot-common/pkg/api"
	"github.com/anodot/anodot-common/pkg/client"
	"github.com/anodot/anodot-common/pkg/events"
	"github.com/anodot/kube-events/pkg/configuration"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "k8s.io/klog/v2"
)

var labels = []string{"anodot_url"}

var (
	anodotServerResponseTime = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name:       "anodot_server_response_time_seconds",
		Help:       "Anodot server response time in seconds",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}, labels)

	totalEvents = promauto.NewCounter(prometheus.CounterOpts{
		Name: "anodot_kube_event_total_produced",
		Help: "Total number of events produced",
	})

	eventErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "anodot_kube_event_errors_count",
		Help: "Total number of error while processing events",
	})
)

type KubernetesEventsHandler interface {
	//multiple Anodot Events could be generated from single Kubernetes event
	EventData(event Event) ([]events.Event, error)
	SupportedEvent() string
}

type AnodotEventhandler struct {
	anodotApi     *api.Api
	handlers      map[string]KubernetesEventsHandler
	configuration configuration.Configuration
}

func NewAnodotEventHandler(anodotURL string, apiToken string, config configuration.Configuration) (*AnodotEventhandler, error) {
	anodotClient, err := client.NewAnodotClient(anodotURL, apiToken, nil)
	if err != nil {
		return nil, err
	}
	apiClient, err := api.NewApiClient(anodotClient)
	if err != nil {
		return nil, err
	}

	deploymentHandler := DeploymentHandler{EventConfig: config.Deployment.EventConfig}
	configmapHandler := ConfigmapHandler{EventConfig: config.ConfigMap.EventConfig}
	daemonsetHandler := DaemonsetHandler{EventConfig: config.DaemonSet.EventConfig}
	statefulsetHandler := StatefulSetHandler{EventConfig: config.StatefulSet.EventConfig}
	secretsHandler := SecretHandler{EventConfig: config.Secret.EventConfig}
	//jobHandler := JobHandler{eventConfig}

	return &AnodotEventhandler{
		anodotApi:     apiClient,
		configuration: config,
		handlers: map[string]KubernetesEventsHandler{
			strings.ToLower(deploymentHandler.SupportedEvent()):  &deploymentHandler,
			strings.ToLower(configmapHandler.SupportedEvent()):   &configmapHandler,
			strings.ToLower(daemonsetHandler.SupportedEvent()):   &deploymentHandler,
			strings.ToLower(statefulsetHandler.SupportedEvent()): &statefulsetHandler,
			strings.ToLower(secretsHandler.SupportedEvent()):     &secretsHandler,
		}}, nil
}

func (a *AnodotEventhandler) Handle(event Event) {
	log.V(5).Infof("Processing event: %s", event)
	if v, ok := a.handlers[strings.ToLower(event.ResourceType)]; ok {
		eventData, err := v.EventData(event)
		if err != nil {
			log.Error("failed to get event data: ", err.Error())
			eventErrors.Inc()
			return
		}

		for _, ev := range eventData {
			go func(e events.Event) {
				ts := time.Now()

				for k, v := range a.configuration.Properties {
					e.Properties = append(e.Properties, events.EventProperties{
						Key:   k,
						Value: v,
					})
				}

				totalEvents.Inc()
				_, err = a.anodotApi.Events.Create(e)
				if err != nil {
					log.Error("failed to send event: ", err.Error())
					eventErrors.Inc()
					return
				}
				anodotServerResponseTime.WithLabelValues(a.anodotApi.AnodotURL().Host).Observe(time.Since(ts).Seconds())
			}(ev)
		}
	} else {
		log.Warningf("unsupported event %s", event)
	}
}
