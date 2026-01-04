package deployments

import (
	"context"
	"fmt"

	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/pod_library/client"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeploymentSpec struct {
	Name     string
	Image    string
	Replicas int32
	Port     int32
	Labels   map[string]string
}

// Get Deployments
func ListDeployments(
	ctx context.Context,
	c *client.K8sClient,
	namespace string,
) ([]appsv1.Deployment, error) {

	result, err := c.Clientset.
		AppsV1().
		Deployments(namespace).
		List(ctx, metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	return result.Items, nil
}

// Creates a single deployment
func CreateDeployment(
	ctx context.Context,
	c *client.K8sClient,
	namespace string,
	spec DeploymentSpec,
) (*appsv1.Deployment, error) {

	// Ensuring there is min of 1 pod running during deployment creation
	if spec.Replicas == 0 {
		spec.Replicas = 1
	}

	if spec.Labels == nil {
		spec.Labels = map[string]string{
			"app": spec.Name,
		}
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:   spec.Name,
			Labels: spec.Labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: spec.Labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: spec.Labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  spec.Name,
							Image: spec.Image,
							Ports: []corev1.ContainerPort{
								{ContainerPort: spec.Port},
							},
							Env: []corev1.EnvVar{
								{
									Name: "POD_NAME",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.name",
										},
									},
								},
								{
									Name: "POD_NAMESPACE",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.namespace",
										},
									},
								},
								{
									Name: "POD_UID",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.uid",
										},
									},
								},
								{
									Name: "POD_IP",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "status.podIP",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return c.Clientset.
		AppsV1().
		Deployments(namespace).
		Create(ctx, deployment, metav1.CreateOptions{})
}

// Scale pods using deployments - create a replica set
func ScaleDeployment(
	ctx context.Context,
	c *client.K8sClient,
	namespace string,
	name string,
	replicas int32,
) error {

	scale, err := c.Clientset.
		AppsV1().
		Deployments(namespace).
		GetScale(ctx, name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("get scale for deployment %s/%s: %w", namespace, name, err)
	}

	scale.Spec.Replicas = replicas

	_, err = c.Clientset.
		AppsV1().
		Deployments(namespace).
		UpdateScale(ctx, name, scale, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("update scale for deployment %s/%s: %w", namespace, name, err)
	}

	return nil
}

// Delete Deployments
func DeleteDeployment(
	ctx context.Context,
	c *client.K8sClient,
	namespace, name string,
) error {

	policy := metav1.DeletePropagationForeground

	return c.Clientset.
		AppsV1().
		Deployments(namespace).
		Delete(ctx, name, metav1.DeleteOptions{
			PropagationPolicy: &policy,
		})
}

// Update Deployments - during code changes
func UpdateDeploymentImage(
	ctx context.Context,
	c *client.K8sClient,
	namespace, name, newImage string,
) error {

	dep, err := c.Clientset.
		AppsV1().
		Deployments(namespace).
		Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	dep.Spec.Template.Spec.Containers[0].Image = newImage

	_, err = c.Clientset.
		AppsV1().
		Deployments(namespace).
		Update(ctx, dep, metav1.UpdateOptions{})

	return err
}
