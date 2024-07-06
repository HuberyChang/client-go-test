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

	appsv1 "k8s.io/api/apps/v1"
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

	// 创建一个 SharedIndexInformer 来监视 Deployment 资源
	sharedInformerFactory := informers.NewSharedInformerFactory(clientset, time.Second*30)
	deploymentInformer := sharedInformerFactory.Apps().V1().Deployments().Informer()

	// 定义 Controller 控制逻辑
	// controller := NewDeploymentController(clientset, deploymentInformer)

	// 启动 Informer，开始监视 Deployment 资源
	go deploymentInformer.Run(make(chan struct{}))

	// 等待 Informer 启动和同步完成
	if !cache.WaitForCacheSync(make(chan struct{}), deploymentInformer.HasSynced) {
		log.Fatalf("Error syncing cache")
	}

	log.Printf("Deployment Informer synced and ready")

	// 保持程序运行，等待事件处理
	// runtime.NewControllerStopChannel().Run()
}

// DeploymentController 定义一个简单的 Deployment 控制器
type DeploymentController struct {
	clientset *kubernetes.Clientset
	informer  cache.SharedIndexInformer
}

// NewDeploymentController 创建一个新的 DeploymentController
func NewDeploymentController(clientset *kubernetes.Clientset, informer cache.SharedIndexInformer) *DeploymentController {
	controller := &DeploymentController{
		clientset: clientset,
		informer:  informer,
	}

	// 注册事件处理函数
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    controller.handleDeploymentAdd,
		UpdateFunc: controller.handleDeploymentUpdate,
		DeleteFunc: controller.handleDeploymentDelete,
	})

	return controller
}

// 处理 Deployment 资源的添加事件
func (c *DeploymentController) handleDeploymentAdd(obj interface{}) {
	deployment := obj.(*appsv1.Deployment)
	fmt.Printf("Deployment added: %s/%s\n", deployment.Namespace, deployment.Name)

	// 在这里编写自定义的添加处理逻辑，例如自动扩展 Pod 副本数
}

// 处理 Deployment 资源的更新事件
func (c *DeploymentController) handleDeploymentUpdate(oldObj, newObj interface{}) {
	newDeployment := newObj.(*appsv1.Deployment)
	// oldDeployment := oldObj.(*appsv1.Deployment)
	fmt.Printf("Deployment updated: %s/%s\n", newDeployment.Namespace, newDeployment.Name)

	// 在这里编写自定义的更新处理逻辑
}

// 处理 Deployment 资源的删除事件
func (c *DeploymentController) handleDeploymentDelete(obj interface{}) {
	deployment := obj.(*appsv1.Deployment)
	fmt.Printf("Deployment deleted: %s/%s\n", deployment.Namespace, deployment.Name)

	// 在这里编写自定义的删除处理逻辑
}
