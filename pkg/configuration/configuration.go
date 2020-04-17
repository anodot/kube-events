package configuration

import (
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/api/meta"
	log "k8s.io/klog/v2"
	"os"
	"regexp"
	"strings"
)

type Configuration struct {
	Deployment            Resource
	StatefulSet           Resource
	ReplicationController Resource
	ReplicaSet            Resource
	DaemonSet             Resource
	Services              Resource
	Pod                   Resource
	Job                   Resource
	PersistentVolume      Resource
	Namespace             Resource
	Secret                Resource
	ConfigMap             Resource
	Ingress               Resource

	Properties map[string]string `yaml:"-,omitempty"`
}

func NewFromYaml(d []byte) (*Configuration, error) {
	configuration := &Configuration{}
	err := yaml.Unmarshal(d, configuration)
	if err != nil {
		return nil, err
	}

	configuration.Properties = map[string]string{}
	for _, s := range os.Environ() {
		split := strings.Split(s, "=")
		k := split[0]
		v := split[1]

		if strings.HasPrefix(k, "ANODOT_EVENTS_PROPS_") {
			configuration.Properties[strings.ToLower(strings.TrimPrefix(k, "ANODOT_EVENTS_PROPS_"))] = v
		}
	}

	return configuration, nil
}

type Resource struct {
	Enabled       bool
	Namespace     string
	EventConfig   EventConfig `yaml:"eventConfig,omitempty"`
	FilterOptions ObjectFiler `yaml:"exclude,omitempty"`
}

func (r *Resource) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type ref Resource
	err := unmarshal((*ref)(r))
	if err != nil {
		return err
	}
	eventCategory := os.Getenv("ANODOT_EVENT_CATEGORY")
	eventSource := os.Getenv("ANODOT_EVENT_SOURCE")

	if len(strings.TrimSpace(r.EventConfig.Source)) == 0 {
		r.EventConfig.Source = eventSource
	}
	if len(strings.TrimSpace(r.EventConfig.Category)) == 0 {
		r.EventConfig.Category = eventCategory
	}
	return nil
}

type EventConfig struct {
	Category string
	Source   string
}

type ObjectFiler struct {
	Name        string
	Namespace   []string
	Labels      map[string]string
	Annotations map[string]string
}

func (o *ObjectFiler) IsExcluded(obj interface{}) bool {
	object, err := meta.Accessor(obj)
	if err != nil {
		log.Error(err.Error() + ". IsExcluded is set to 'true' ")
		return true
	}

	matchNamespace := o.MatchNamespace(object.GetNamespace())
	matchAnnotations := o.MatchAnnotations(object.GetAnnotations())
	matchLabels := o.MatchLabels(object.GetLabels())
	matchName := o.MatchName(object.GetName())

	//TODO better logs
	if matchAnnotations || matchNamespace || matchLabels || matchName {
		log.V(5).Infof("'%s' excluded by name=%t, annotations=%t, by labels=%t, by namespace=%t", object.GetName(), matchName, matchAnnotations, matchLabels, matchNamespace)
		return true
	}

	return false
}

func (o ObjectFiler) MatchName(name string) bool {
	//name not set
	if len(strings.TrimSpace(o.Name)) == 0 {
		return false
	}

	return regexp.MustCompile(o.Name).MatchString(name)
}

func (o ObjectFiler) MatchLabels(m map[string]string) bool {
	for kk, vv := range m {
		if v, ok := o.Labels[kk]; ok && regexp.MustCompile(v).MatchString(vv) {
			return true
		}
	}
	return false
}

func (o ObjectFiler) MatchAnnotations(m map[string]string) bool {
	for kk, vv := range m {
		if v, ok := o.Annotations[kk]; ok && regexp.MustCompile(v).MatchString(vv) {
			return true
		}
	}
	return false
}

func (o ObjectFiler) MatchNamespace(namespace string) bool {
	for _, v := range o.Namespace {
		if regexp.MustCompile(v).MatchString(namespace) {
			return true
		}
	}
	return false
}
