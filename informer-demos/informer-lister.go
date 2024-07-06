package main

import (
	"log"
	"time"

	"k8s.io/apimachinery/pkg/labels"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// 创建一个 SharedIndexInformer 来监视 Pod 资源
	sharedInformerFactory := informers.NewSharedInformerFactory(clientset, time.Second*30)
	podInformer := sharedInformerFactory.Core().V1().Pods().Informer()

	// 启动 Informer，开始监视 Pod 资源
	stopCh := make(chan struct{})
	go func() {
		defer close(stopCh)
		podInformer.Run(stopCh)
	}()

	// 等待 Informer 启动和同步完成
	if !cache.WaitForCacheSync(make(chan struct{}), podInformer.HasSynced) {
		panic("failed to sync")
	}

	log.Printf("Pod Informer synced and ready")

	// 获取 Lister
	podLister := sharedInformerFactory.Core().V1().Pods().Lister()

	// 示例：通过 Lister 获取 default 命名空间下的 Pod 列表
	// pods, err := podLister.Pods("default").List(labels.Everything())
	// 示例：通过 Lister 获取所有命名空间下的 Pod 列表
	pods, err := podLister.List(labels.Everything())
	if err != nil {
		panic(err)
	}

	// 打印所有 Pod 的名称和命名空间
	for _, pod := range pods {
		log.Printf("Pod: %s", pod.Name)
	}
}
