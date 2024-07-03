package main

import (
	"context"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}

	client, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	deploymentResource := schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "deployments",
	}
	deploymentSpec := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name": "nginx-deploy-test",
			},
			"spec": map[string]interface{}{
				"replicas": 1,
				"selector": map[string]interface{}{
					"matchLabels": map[string]interface{}{
						"app": "nginx-deploy-test",
					},
				},
				"template": map[string]interface{}{
					"metadata": map[string]interface{}{
						"labels": map[string]interface{}{
							"app": "nginx-deploy-test",
						},
					},
					"spec": map[string]interface{}{
						"containers": []map[string]interface{}{
							{
								"name":  "nginx-deploy-test1",
								"image": "nginx:1.7.9",
								"ports": []map[string]interface{}{
									{
										"containerPort": 80,
										"protocol":      "TCP",
										"name":          "http",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	fmt.Println("creating deployment")
	result, err := client.Resource(deploymentResource).Namespace("default").Create(context.Background(), deploymentSpec, v1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Create deployment %q.\n", result.GetName())
}
