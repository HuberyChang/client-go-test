package main

import (
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type Controller struct {
	indexer  cache.Indexer
	informer cache.Controller
	queue    workqueue.RateLimitingInterface
}

func NewController(indexer cache.Indexer, informer cache.Controller, queue workqueue.RateLimitingInterface) *Controller {
	return &Controller{
		indexer:  indexer,
		informer: informer,
		queue:    queue,
	}
}

func (c *Controller) Run(workers int, stopCh chan struct{}) {
	go c.informer.Run(stopCh)
	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		panic("failed to wait for cache to sync")
	}
	for i := 0; i < workers; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}
	<-stopCh
	fmt.Println("stop")
	//for c.processNextItem() {
	//}
}

func (c *Controller) processNextItem() bool {
	key, quit := c.queue.Get()

	if quit {
		return false
	}

	defer c.queue.Done(key)

	keystr, ok := key.(string)
	if !ok {
		fmt.Println("Error: key is not of type string")
		return false
	}

	obj, exists, err := c.indexer.GetByKey(keystr)
	if err != nil {
		panic(err)
	}
	if !exists {
		fmt.Println("not exists")
	} else {
		fmt.Println("exists", obj.(*corev1.Pod).GetName())
	}

	return true
}

func (c *Controller) runWorker() {
	for c.processNextItem() {
	}
}

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}

	kubernetesClient, err := kubernetes.NewForConfig(config)

	podListWatcher := cache.NewListWatchFromClient(kubernetesClient.CoreV1().RESTClient(), "pods", v1.NamespaceDefault, fields.Everything())

	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	indexer, informer := cache.NewIndexerInformer(podListWatcher, &corev1.Pod{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				fmt.Println("add")
				queue.Add(key)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(newObj)
			if err == nil {
				fmt.Println("update")
				queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				fmt.Println("delete")
				queue.Add(key)
			}
		},
	}, cache.Indexers{})

	stop := make(chan struct{})
	controller := NewController(indexer, informer, queue)
	go controller.Run(1, stop)

	select {}
}
