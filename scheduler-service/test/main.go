package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/pod_library/client"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/pod_library/deployments"
)

func main() {
	ctx := context.Background()

	// Create Kubernetes client
	k8sClient, err := client.New()
	if err != nil {
		panic(fmt.Errorf("failed to create k8s client: %w", err))
	}

	namespace := "default"

	// Define deployment spec
	spec := deployments.DeploymentSpec{
		Name:     "demo-worker-deployment",
		Image:    "worker:1.0",
		Replicas: 2,
		Port:     80,
	}

	// Create deployment
	dep, err := deployments.CreateDeployment(ctx, k8sClient, namespace, spec)
	if err != nil {
		panic(fmt.Errorf("failed to create deployment: %w", err))
	}

	fmt.Printf("Deployment created: %s\n", dep.Name)

	// Wait a bit for pods to come up
	time.Sleep(5 * time.Second)

	// List deployments
	deploymentsList, err := deployments.ListDeployments(ctx, k8sClient, namespace)
	if err != nil {
		panic(fmt.Errorf("failed to list deployments: %w", err))
	}

	fmt.Println("\nDeployments in namespace:")
	for _, d := range deploymentsList {
		fmt.Printf(
			"- %s (replicas: %d)\n",
			d.Name,
			*d.Spec.Replicas,
		)
	}

	// Delete deplyments

}
