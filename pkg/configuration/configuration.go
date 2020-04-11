package configuration

import (
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/api/meta"
	log "k8s.io/klog/v2"
	"regexp"
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
}

func NewFromYaml(d []byte) (*Configuration, error) {
	configuration := &Configuration{}
	err := yaml.Unmarshal(d, configuration)
	if err != nil {
		return nil, err
	}
	return configuration, nil
}

type Resource struct {
	Enabled       bool
	Namespace     string
	FilterOptions ObjectFiler `yaml:"exclude,omitempty"`
}

//TOOO implement this
type EventConfig struct {
	Category string
	Source   string
}

type ObjectFiler struct {
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

	//TODO better logs
	if matchAnnotations || matchNamespace || matchLabels {
		log.V(5).Infof("'%s' excluded by annotations=%t, by labels=%t, by namespace=%t", object.GetName(), matchAnnotations, matchLabels, matchNamespace)
		return true
	}

	return false
}

func (o *ObjectFiler) MatchLabels(m map[string]string) bool {
	for kk, vv := range m {
		if v, ok := o.Labels[kk]; ok && regexp.MustCompile(v).MatchString(vv) {
			return true
		}
	}
	return false
}

func (o *ObjectFiler) MatchAnnotations(m map[string]string) bool {
	for kk, vv := range m {
		if v, ok := o.Annotations[kk]; ok && regexp.MustCompile(v).MatchString(vv) {
			return true
		}
	}
	return false
}

func (o *ObjectFiler) MatchNamespace(namespace string) bool {
	for _, v := range o.Namespace {
		if v == namespace {
			return true
		}
	}
	return false
}
