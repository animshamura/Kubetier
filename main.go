package main

import (
	"context"
	"fmt"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func int32Ptr(i int32) *int32 { return &i }

func main() {
	kubeconfig := filepath.Join(
		homeDir(), ".kube", "config",
	)

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	namespace := "default"
	ctx := context.Background()

	// Define your apps
	components := []struct {
		Name       string
		Image      string
		Port       int32
		TargetPort int32
	}{
		{"frontend", "your-dockerhub-username/frontend:latest", 80, 80},
		{"backend", "your-dockerhub-username/backend:latest", 8080, 8080},
		{"ml-model", "your-dockerhub-username/ml-model:latest", 8000, 8000},
	}

	for _, comp := range components {
		// --- Create Deployment ---
		deploy := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name: comp.Name,
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: int32Ptr(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"app": comp.Name},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{"app": comp.Name},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  comp.Name,
								Image: comp.Image,
								Ports: []corev1.ContainerPort{
									{ContainerPort: comp.Port},
								},
							},
						},
					},
				},
			},
		}

		_, err := clientset.AppsV1().Deployments(namespace).Create(ctx, deploy, metav1.CreateOptions{})
		if err != nil {
			fmt.Printf("Deployment %s: %v\n", comp.Name, err)
		} else {
			fmt.Printf("Deployment created: %s\n", comp.Name)
		}

		// --- Create Service ---
		service := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name: comp.Name + "-service",
			},
			Spec: corev1.ServiceSpec{
				Selector: map[string]string{"app": comp.Name},
				Ports: []corev1.ServicePort{
					{
						Port:       comp.Port,
						TargetPort: intstr.FromInt(int(comp.TargetPort)),
					},
				},
			},
		}

		_, err = clientset.CoreV1().Services(namespace).Create(ctx, service, metav1.CreateOptions{})
		if err != nil {
			fmt.Printf("Service %s: %v\n", comp.Name, err)
		} else {
			fmt.Printf("Service created: %s-service\n", comp.Name)
		}
	}
}

func homeDir() string {
	if h := filepath.Join("/home", "your-username"); h != "" {
		return h
	}
	return "/root"
}
