package controller

import (
	"fmt"
	"github.com/anodot/kube-events/pkg/configuration"
	"github.com/anodot/kube-events/pkg/handlers"
	"github.com/anodot/kube-events/pkg/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	k8sruntime "k8s.io/apimachinery/pkg/util/runtime"
	log "k8s.io/klog/v2"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	apps_v1beta1 "k8s.io/api/apps/v1beta1"
	batch_v1 "k8s.io/api/batch/v1"
	api_v1 "k8s.io/api/core/v1"
	ext_v1beta1 "k8s.io/api/extensions/v1beta1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

var (
	k8sErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "anodot_kube_event_kubernetes_error_count",
		Help: "Total number of error while processing events",
	})
)

var serverStartTime time.Time

type EventHandler interface {
	DoHandle()
}

type Controller struct {
	clientset    kubernetes.Interface
	queue        workqueue.RateLimitingInterface
	informer     cache.SharedIndexInformer
	eventHandler *handlers.AnodotEventhandler
}

// Start prepares watchers and run their controllers, then waits for process termination signals
func Start(conf configuration.Configuration, eventHandler *handlers.AnodotEventhandler) {
	var kubeClient kubernetes.Interface
	_, err := rest.InClusterConfig()
	if err != nil {
		kubeClient = utils.GetClientOutOfCluster()
	} else {
		kubeClient = utils.GetClient()
	}

	go func() {
		serveMux := http.NewServeMux()
		serveMux.Handle("/metrics", promhttp.Handler())
		serveMux.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
			_, err := fmt.Fprintf(writer, `{"time": "%d"}`, time.Now().Unix())
			if err != nil {
				log.Error(err.Error())
			}
		})
		err := http.ListenAndServe(":8080", serveMux)
		if err != nil {
			log.Error("failed to initialize metrics endpoint: ", err.Error())
		}
	}()

	if conf.Pod.Enabled {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.CoreV1().Pods(conf.Pod.Namespace).List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.CoreV1().Pods(conf.Pod.Namespace).Watch(options)
				},
			},
			&api_v1.Pod{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandler, informer, "pod", conf.Pod.FilterOptions)
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.DaemonSet.Enabled {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.ExtensionsV1beta1().DaemonSets(conf.DaemonSet.Namespace).List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.ExtensionsV1beta1().DaemonSets(conf.DaemonSet.Namespace).Watch(options)
				},
			},
			&ext_v1beta1.DaemonSet{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandler, informer, "daemonset", conf.DaemonSet.FilterOptions)
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.ReplicaSet.Enabled {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.ExtensionsV1beta1().ReplicaSets(conf.ReplicaSet.Namespace).List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.ExtensionsV1beta1().ReplicaSets(conf.ReplicaSet.Namespace).Watch(options)
				},
			},
			&ext_v1beta1.ReplicaSet{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandler, informer, "replicaset", conf.ReplicaSet.FilterOptions)
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.Services.Enabled {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.CoreV1().Services(conf.Services.Namespace).List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.CoreV1().Services(conf.Services.Namespace).Watch(options)
				},
			},
			&api_v1.Service{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandler, informer, "service", conf.Services.FilterOptions)
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.Deployment.Enabled {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.AppsV1beta1().Deployments(conf.Deployment.Namespace).List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.AppsV1beta1().Deployments(conf.Deployment.Namespace).Watch(options)
				},
			},
			&apps_v1beta1.Deployment{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandler, informer, "deployment", conf.Deployment.FilterOptions)
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.Namespace.Enabled {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.CoreV1().Namespaces().List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.CoreV1().Namespaces().Watch(options)
				},
			},
			&api_v1.Namespace{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandler, informer, "namespace", conf.Namespace.FilterOptions)
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.ReplicationController.Enabled {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.CoreV1().ReplicationControllers(conf.ReplicationController.Namespace).List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.CoreV1().ReplicationControllers(conf.ReplicationController.Namespace).Watch(options)
				},
			},
			&api_v1.ReplicationController{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandler, informer, "replication controller", conf.ReplicationController.FilterOptions)
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.Job.Enabled {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.BatchV1().Jobs(conf.Job.Namespace).List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.BatchV1().Jobs(conf.Job.Namespace).Watch(options)
				},
			},
			&batch_v1.Job{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandler, informer, "job", conf.Job.FilterOptions)
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.PersistentVolume.Enabled {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.CoreV1().PersistentVolumes().List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.CoreV1().PersistentVolumes().Watch(options)
				},
			},
			&api_v1.PersistentVolume{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandler, informer, "persistent volume", conf.PersistentVolume.FilterOptions)
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.Secret.Enabled {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.CoreV1().Secrets(conf.Secret.Namespace).List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.CoreV1().Secrets(conf.Secret.Namespace).Watch(options)
				},
			},
			&api_v1.Secret{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandler, informer, "secret", conf.Secret.FilterOptions)
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.ConfigMap.Enabled {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.CoreV1().ConfigMaps(conf.ConfigMap.Namespace).List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.CoreV1().ConfigMaps(conf.ConfigMap.Namespace).Watch(options)
				},
			},
			&api_v1.ConfigMap{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandler, informer, "configmap", conf.ConfigMap.FilterOptions)
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.Ingress.Enabled {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.ExtensionsV1beta1().Ingresses(conf.Ingress.Namespace).List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.ExtensionsV1beta1().Ingresses(conf.Ingress.Namespace).Watch(options)
				},
			},
			&ext_v1beta1.Ingress{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandler, informer, "ingress", conf.Ingress.FilterOptions)
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.StatefulSet.Enabled {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.AppsV1beta1().StatefulSets(conf.StatefulSet.Namespace).List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.AppsV1beta1().StatefulSets(conf.StatefulSet.Namespace).Watch(options)
				},
			},
			&apps_v1beta1.StatefulSet{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandler, informer, "ingress", conf.Ingress.FilterOptions)
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM)
	signal.Notify(sigterm, syscall.SIGINT)
	<-sigterm
}

func newResourceController(client kubernetes.Interface, eventHandler *handlers.AnodotEventhandler, informer cache.SharedIndexInformer, resourceType string, filter configuration.ObjectFiler) *Controller {
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	var newEvent handlers.Event
	var err error
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			newEvent.Key, err = cache.MetaNamespaceKeyFunc(obj)
			newEvent.EventType = "create"
			newEvent.ResourceType = resourceType
			newEvent.EventTime = time.Now()
			log.V(4).Infof("Processing add to %v: %s", resourceType, newEvent.Key)
			if err == nil {
				if !filter.IsExcluded(obj) {
					queue.Add(newEvent)
				}
			}
		},
		UpdateFunc: func(old, new interface{}) {
			newEvent.Key, err = cache.MetaNamespaceKeyFunc(old)
			newEvent.EventType = "update"
			newEvent.ResourceType = resourceType
			newEvent.EventTime = time.Now()

			newEvent.Old = old
			newEvent.New = new
			log.V(4).Infof("Processing update to %v: %s", resourceType, newEvent.Key)

			if err == nil {
				if !filter.IsExcluded(new) {
					queue.Add(newEvent)
				}
			}
		},
		DeleteFunc: func(obj interface{}) {
			newEvent.Key, err = cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			newEvent.EventType = "delete"
			newEvent.ResourceType = resourceType
			newEvent.EventTime = time.Now()
			log.V(4).Infof("Processing delete to %v: %s", resourceType, newEvent.Key)
			if err == nil {
				if !filter.IsExcluded(obj) {
					queue.Add(newEvent)
				}
			}
		},
	})

	return &Controller{
		clientset:    client,
		informer:     informer,
		queue:        queue,
		eventHandler: eventHandler,
	}
}

// Run starts the kubewatch controller
func (c *Controller) Run(stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	log.V(3).Info("Starting kubewatch controller")
	serverStartTime = time.Now().Local()

	k8sruntime.ErrorHandlers = append(k8sruntime.ErrorHandlers, func(err error) {
		k8sErrors.Inc()
	})

	go c.informer.Run(stopCh)

	if !cache.WaitForCacheSync(stopCh, c.HasSynced) {
		utilruntime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	log.V(3).Info("Kubewatch controller synced and ready")

	wait.Until(c.runWorker, time.Second, stopCh)
}

// HasSynced is required for the cache.Controller interface.
func (c *Controller) HasSynced() bool {
	return c.informer.HasSynced()
}

// LastSyncResourceVersion is required for the cache.Controller interface.
func (c *Controller) LastSyncResourceVersion() string {
	return c.informer.LastSyncResourceVersion()
}

func (c *Controller) runWorker() {
	for c.processNextItem() {
		// continue looping
	}
}

func (c *Controller) processNextItem() bool {
	newEvent, quit := c.queue.Get()

	if quit {
		return false
	}
	defer c.queue.Done(newEvent)
	err := c.processItem(newEvent.(handlers.Event))
	if err == nil {
		// No error, reset the ratelimit counters
		c.queue.Forget(newEvent)
		//TODO max retries
	} else if c.queue.NumRequeues(newEvent) < 5 {
		log.Errorf("Error processing %s (will retry): %v", newEvent.(handlers.Event).Key, err)
		c.queue.AddRateLimited(newEvent)
	} else {
		// err != nil and too many retries
		log.Errorf("Error processing %s (giving up): %v", newEvent.(handlers.Event).Key, err)
		c.queue.Forget(newEvent)
		utilruntime.HandleError(err)
	}

	return true
}

func (c *Controller) processItem(newEvent handlers.Event) error {
	obj, _, err := c.informer.GetIndexer().GetByKey(newEvent.Key)
	if err != nil {
		return fmt.Errorf("error fetching object with key %s from store: %v", newEvent.Key, err)
	}

	objectMeta := utils.GetObjectMetaData(obj)

	switch newEvent.EventType {
	case "create":
		// compare CreationTimestamp and serverStartTime and alert only on latest events
		// Could be Replaced by using Delta or DeltaFIFO
		if objectMeta.CreationTimestamp.Sub(serverStartTime).Seconds() > 0 {
			c.eventHandler.Handle(newEvent)
			return nil
		}
	case "update":
		c.eventHandler.Handle(newEvent)
		return nil
	case "delete":
		c.eventHandler.Handle(newEvent)
		return nil
	}
	return nil
}
