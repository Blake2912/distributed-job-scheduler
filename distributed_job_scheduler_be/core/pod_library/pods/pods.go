package pods

import (
	"context"

	"example.com/pod_library/client"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func List(ctx context.Context, c *client.K8sClient, namespace string) ([]corev1.Pod, error) {
	pods, err := c.Clientset.
		CoreV1().
		Pods(namespace).
		List(ctx, metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	return pods.Items, nil
}
