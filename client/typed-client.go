package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"k8s.io/client-go/util/retry"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "demo-deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "demo",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "demo",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "web",
							Image: "nginx:1.7.9",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: 80,
									Protocol:      apiv1.ProtocolTCP,
								},
							},
						},
					},
				},
			},
		},
	}

	// Create Deployment
	fmt.Printf("Creating deployment...\n")
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

	// Update Deployments
	prompt()
	fmt.Println("Updating deployment...")
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		deployment, getErr := deploymentsClient.Get(context.TODO(), "demo-deployment", metav1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("failed to get latest version of Deployment: %v", getErr))
		}

		deployment.Spec.Replicas = int32Ptr(1)
		deployment.Spec.Template.Spec.Containers[0].Image = "nginx:1.9.1"
		deployment.Spec.Template.Spec.Containers[0].ImagePullPolicy = apiv1.PullIfNotPresent
		_, updateErr := deploymentsClient.Update(context.TODO(), deployment, metav1.UpdateOptions{})
		return updateErr
	})

	if retryErr != nil {
		panic(fmt.Errorf("update deployment failed: %v", err))
	}
	fmt.Println("Updated deployment.")

	// List Deployments
	prompt()
	fmt.Printf("Listing deployments in namespace %q:\n", apiv1.NamespaceDefault)

	list, err := deploymentsClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, d := range list.Items {
		fmt.Printf(" * %s (%d replicas)\n", d.Name, *d.Spec.Replicas)
	}

	// Delete Deployment
	prompt()
	fmt.Printf("Deleting deployment...\n")
	deletePolicy := metav1.DeletePropagationForeground
	if err := deploymentsClient.Delete(context.TODO(), "demo-deployment", metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err.Error())
	}
	fmt.Println("Deleted deployment.")
}

func prompt() {
	fmt.Printf("-> Press Return key to continue.")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		break
	}

	if err := scanner.Err(); err != nil {
		panic(err.Error())
	}

	fmt.Println()
}

func int32Ptr(i int32) *int32 {
	return &i
}
