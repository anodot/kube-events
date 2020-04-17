package configuration

import (
	"fmt"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"testing"
)

func TestMarshl(t *testing.T) {
	configuration := Configuration{
		Deployment: Resource{
			Enabled:   false,
			Namespace: "default",
			FilterOptions: ObjectFiler{
				Namespace:   []string{"as", "bb"},
				Labels:      map[string]string{"a": "3"},
				Annotations: map[string]string{"a": "3"},
			},
		},
		ReplicationController: Resource{},
		ReplicaSet:            Resource{},
		DaemonSet:             Resource{},
		Services:              Resource{},
		Pod:                   Resource{},
		Job:                   Resource{},
		PersistentVolume:      Resource{},
		Namespace:             Resource{},
		Secret:                Resource{},
		ConfigMap:             Resource{},
		Ingress:               Resource{},
	}

	bytes, err := yaml.Marshal(configuration)
	if err != nil {
		t.Fatalf(err.Error())
	}

	fmt.Println(string(bytes))
}

func TestUnmarshall(t *testing.T) {

	os.Setenv("ANODOT_EVENT_SOURCE", "vova")
	os.Setenv("ANODOT_EVENT_CATEGORY", "vova-cat")

	configYaml := `deployment:
  enabled: true
  namespace: ""
  eventConfig:
    category: "test-category"
  exclude:
    namespace: []
    labels: {}
    annotations: {}
configmap:
  enabled: true
  namespace: ""
  exclude:
    namespace: []
    labels:
      OWNER: "TILLER"
    annotations:
      "control-plane.alpha.kubernetes.io/leader": ".*"`
	config, err := NewFromYaml([]byte(configYaml))
	if err != nil {
		t.Fatalf(err.Error())
	}

	if config.Deployment.Enabled != true {
		t.Fatalf("deployment should be enabled")
	}

	if config.Deployment.EventConfig.Source != "vova" {
		t.Fatal("value:", config.Deployment.EventConfig.Source)
	}

}

func TestExclude(t *testing.T) {
	configYaml := `deployment:
  enabled: true
  namespace: ""
  eventConfig:
    category: "test-category"
    source: "test-source"
  exclude:
    namespace: []
    labels: {}
    annotations: {}
configmap:
  enabled: true
  namespace: ""
  exclude:
    namespace: []
    labels:
      OWNER: "TILLER"
      MANAGED_BY: ".*"
    annotations:
      "control-plane.alpha.kubernetes.io/leader": ".*"`
	config, err := NewFromYaml([]byte(configYaml))
	if err != nil {
		t.Fatalf(err.Error())
	}

	excludedByLabel := &v1.ConfigMap{
		TypeMeta: meta_v1.TypeMeta{},
		ObjectMeta: meta_v1.ObjectMeta{
			Name:   "test",
			Labels: map[string]string{"OWNER": "TILLER"},
		},
		Data:       nil,
		BinaryData: nil,
	}

	excludeByLabel := config.ConfigMap.FilterOptions.IsExcluded(excludedByLabel)
	if !excludeByLabel {
		t.Fatalf("should be exluded by labels")
	}

	excludedByLabelRegex := &v1.ConfigMap{
		TypeMeta: meta_v1.TypeMeta{},
		ObjectMeta: meta_v1.ObjectMeta{
			Name:   "test",
			Labels: map[string]string{"MANAGED_BY": "bla-bla-bla"},
		},
		Data:       nil,
		BinaryData: nil,
	}

	excludeByLabelRegex := config.ConfigMap.FilterOptions.IsExcluded(excludedByLabelRegex)
	if !excludeByLabelRegex {
		t.Fatalf("should be exluded by labels")
	}

}
