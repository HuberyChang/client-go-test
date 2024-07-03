package main

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err.Error())
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Get all groups
	groupList, err := discoveryClient.ServerGroups()
	if err != nil {
		panic(err.Error())
	}
	spew.Dump(groupList.Groups)

	resources, err := discoveryClient.ServerResourcesForGroupVersion("v1")
	if err != nil {
		panic(err.Error())
	}

	for _, r := range resources.APIResources {
		fmt.Println(r.Name)
	}
}
