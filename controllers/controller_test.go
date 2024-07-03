package main

import (
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type Controller struct {
	indexer  cache.Indexer
	queue    workqueue.TypedRateLimitingInterface[string]
	informer cache.Controller
}

func NewController(indexer cache.Indexer, queue workqueue.TypedRateLimitingInterface[string], informer cache.Controller) *Controller {
	return &Controller{
		indexer:  indexer,
		queue:    queue,
		informer: informer,
	}
}

func (c *Controller) Run(workers int, stopCh chan struct{}) {} {
	defer runtime.HandleCrash()
	defer c.queue.ShutDown()

	klog.Info("Starting Pod controller")
	go c.informer.Run(stopCh)

	if !cache.waitForCacheSync(stopCh, c.informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	for i := 0; i < workers; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}
	<-stopCh
	klog.Info("Stopping Pod controller")
}

func (c *Controller) runWorker() {
	for c.processNextItem() {
	}
}

func (c *Controller) processNextItem() bool {
	c.queue
}
