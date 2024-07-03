package main

import (
	"context"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/scale"

	"k8s.io/client-go/tools/clientcmd"

	autoscalingapi "k8s.io/api/autoscaling/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type resourceMapper struct{}

func (r *resourceMapper) ResourceFor(resource schema.GroupVersionResource) (schema.GroupVersionResource, error) {
	fmt.Printf("ResourceFor was called  with resource %s\n", resource.String())
	if resource.Group == "apps" && resource.Resource == "deployments" {
		return schema.GroupVersionResource{
			Group:    resource.Group,
			Version:  "v1",
			Resource: resource.Resource,
		}, nil
	}
	return schema.GroupVersionResource{}, fmt.Errorf("", clientcmd.RecommendedHomeFile)
}

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		panic(err)
	}

	scaleClient, err := scale.NewForConfig(config, &resourceMapper{}, dynamic.LegacyAPIPathResolverFunc, scale.NewDiscoveryScaleKindResolver(discoveryClient))
	if err != nil {
		panic(err)
	}

	scaleSpec := &autoscalingapi.Scale{
		Spec: autoscalingapi.ScaleSpec{
			Replicas: 1,
		},
		Status: autoscalingapi.ScaleStatus{
			Replicas: 0,
			Selector: "app=nginx-deploy-test",
		},
	}

	scaleSpec.Namespace = "default"
	scaleSpec.ObjectMeta.Name = "nginx-deploy-test"

	_, err = scaleClient.Scales("default").Update(context.Background(), schema.GroupResource{
		Group:    "apps",
		Resource: "deployments",
	}, scaleSpec, v1.UpdateOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s was scaled", scaleSpec.Name)
}
