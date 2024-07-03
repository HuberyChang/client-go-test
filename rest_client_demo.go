package main

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic("failed to build config")
	}

	config.GroupVersion = &v1.SchemeGroupVersion
	config.NegotiatedSerializer = scheme.Codecs
	config.APIPath = "/api"

	restClient, err := rest.RESTClientFor(config)
	if err != nil {
		panic(err)
	}

	pod := &v1.Pod{}
	err = restClient.Get().
		Namespace("default").
		Resource("pods").
		Name("nginx-deploy-797488cdc6-hhjlt").
		Do(context.TODO()).
		Into(pod)
	if err != nil {
		panic(err)
	}
	fmt.Printf("GVK:%s\n", pod.GroupVersionKind())
}
