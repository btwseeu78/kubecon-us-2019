package controller

import (
	webappv1alpha1 "arpan.dev/kubecon-us-2019/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *GuestBookReconciler) desiredDeployment(book webappv1alpha1.GuestBook, redis webappv1alpha1.Redis) (appsv1.Deployment, error) {
	depl := appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: appsv1.SchemeGroupVersion.String(),
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      book.Name,
			Namespace: book.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: book.Spec.Frontend.Replicas, //wont be nil
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"guestbook": book.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "frontend",
							Image: "gcr.io/google-samples/gb-frontend:v4",
							Env: []corev1.EnvVar{
								{
									Name:  "GET_HOSTS_FROM",
									Value: "ENV",
								},
								{
									Name:  "REDIS_MASTER_SERVICE_HOST",
									Value: redis.Status.LeaderService,
								},
								{
									Name:  "REDIS_SLAVE_SERVICE_HOST",
									Value: redis.Status.FollowerService,
								},
							},
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 80,
									Name:          "http",
									Protocol:      "TCP",
								},
							},
							Resources: *book.Spec.Frontend.Resources.DeepCopy(),
						},
					},
				},
			},
		},
	}
	if err := ctrl.SetControllerReference(&book, &depl, r.Scheme); err != nil {
		return depl, err
	}
	return depl, nil
}
