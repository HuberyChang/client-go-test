package main

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
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

	watch, err := clientset.AppsV1().Deployments("default").Watch(context.TODO(), v1.ListOptions{LabelSelector: "app=nginx-deploy-test"})
	if err != nil {
		panic(err)
	}

	for {
		select {
		case event, ok := <-watch.ResultChan():
			if !ok {
				fmt.Println("channel close")
				break
			}
			fmt.Println("Event Type:", event.Type)
			dp, ok := event.Object.(*appsv1.Deployment)
			if !ok {
				fmt.Println("not deploy")
				continue
			}
			fmt.Println(dp)
		}
	}
}
