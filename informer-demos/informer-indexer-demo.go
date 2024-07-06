package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

func main() {
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		flag.StringVar(&kubeconfig, "kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		flag.StringVar(&kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("Error building kubeconfig: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error building kubernetes clientset: %v", err)
	}

	// 创建一个 SharedIndexInformer 来监视 Pod 资源
	sharedInformerFactory := informers.NewSharedInformerFactory(clientset, time.Second*30)
	podInformer := sharedInformerFactory.Core().V1().Pods().Informer()

	// 创建一个 Indexer
	indexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})

	// 定义 Watcher 事件处理函数，并将事件添加到 Indexer 中
	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*corev1.Pod)
			fmt.Printf("Pod added: %s/%s\n", pod.Namespace, pod.Name)
			indexer.Add(pod)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			newPod := newObj.(*corev1.Pod)
			// oldPod := oldObj.(*corev1.Pod)
			fmt.Printf("Pod updated: %s/%s\n", newPod.Namespace, newPod.Name)
			indexer.Update(newPod)
		},
		DeleteFunc: func(obj interface{}) {
			pod := obj.(*corev1.Pod)
			fmt.Printf("Pod deleted: %s/%s\n", pod.Namespace, pod.Name)
			indexer.Delete(pod)
		},
	})

	// 启动 Informer，开始监视 Pod 资源
	go podInformer.Run(make(chan struct{}))

	// 等待 Informer 启动和同步完成
	if !cache.WaitForCacheSync(make(chan struct{}), podInformer.HasSynced) {
		log.Fatalf("Error syncing cache")
	}

	log.Printf("Pod Informer synced and ready")

	// 示例：通过 Indexer 查询所有命名空间下的 Pod 列表
	allPods, err := indexer.ByIndex(cache.NamespaceIndex, "")
	if err != nil {
		log.Fatalf("Error listing pods: %v", err)
	}

	// 打印所有 Pod 的名称和命名空间
	for _, obj := range allPods {
		pod := obj.(*corev1.Pod)
		fmt.Printf("Pod: %s, Namespace: %s\n", pod.Name, pod.Namespace)
	}

	// 保持程序运行，等待事件处理
	// runtime.NewControllerStopChannel().Run()
}
