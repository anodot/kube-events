package main

import (
	"flag"
	"fmt"
	"github.com/anodot/kube-events/pkg/configuration"
	"github.com/anodot/kube-events/pkg/controller"
	"github.com/anodot/kube-events/pkg/handlers"
	"github.com/anodot/kube-events/pkg/version"
	"io/ioutil"
	log "k8s.io/klog/v2"
	"os"
	"runtime"
	"strings"
)

func main() {
	log.InitFlags(nil)
	err := flag.Set("v", defaultIfBlank(os.Getenv("ANODOT_LOG_LEVEL"), "2"))
	if err != nil {
		log.Fatal(err)
	}

	log.Info(fmt.Sprintf("Anodot kube-events version: '%s'. GitSHA: '%s'", version.VERSION, version.REVISION))
	log.V(4).Infof("Go Version: %s", runtime.Version())
	log.V(4).Infof("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)

	anodotURL := os.Getenv("ANODOT_URL")
	anodotApiToken := os.Getenv("ANODOT_API_TOKEN")

	log.V(4).Infof("Anodot Server URL: '%s'", anodotURL)
	log.V(4).Infof("Anodot api token: '%s'", anodotApiToken)
	configLocation := defaultIfBlank(os.Getenv("ANODOT_EVENT_CONFIG_LOCATION"), "/mnt/config.yaml")
	yamlFile, err := ioutil.ReadFile(configLocation)
	if err != nil {
		log.Fatal("failed to open configuration file: ", err.Error())
	}

	config, err := configuration.NewFromYaml(yamlFile)
	if err != nil {
		log.Fatalf(err.Error())
	}

	log.V(5).Infof("configuration: %+v\n", *config)

	anodotEventHandler, err := handlers.NewAnodotEventHandler(anodotURL, anodotApiToken, *config)
	if err != nil {
		log.Fatal(err)
	}

	controller.Start(*config, anodotEventHandler)
}

func defaultIfBlank(actual string, fallback string) string {
	if len(strings.TrimSpace(actual)) == 0 {
		return fallback
	}
	return actual
}
