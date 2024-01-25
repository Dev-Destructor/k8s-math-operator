/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	"io"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	mathv1 "math-controller/api/v1"
)

// ArithmeticReconciler reconciles a Arithmetic object
type ArithmeticReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=math.stream.com,resources=arithmetics,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=math.stream.com,resources=arithmetics/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=math.stream.com,resources=arithmetics/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Arithmetic object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *ArithmeticReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	var problem mathv1.Arithmetic
	if err := r.Get(ctx, req.NamespacedName, &problem); err != nil {
		log.Log.Error(err, "unable to fetch Problem")
		return ctrl.Result{}, err
	}
	log.Log.Info(fmt.Sprintf("Problem: %v", problem.Spec))

	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("job-%s", req.NamespacedName.Name),
			Namespace: req.NamespacedName.Namespace,
		},
		Spec: corev1.PodSpec{
			RestartPolicy: "Never",
			Containers: []corev1.Container{
				{
					Name:  "math",
					Image: "python:latest",
					Args:  []string{"python", "-c", fmt.Sprintf("print(%s)", problem.Spec.Expression)},
				},
			},
		},
	}

	if err := r.Create(ctx, &pod, &client.CreateOptions{}); err != nil {
		log.Log.Error(err, "unable to create pod")
		return ctrl.Result{}, err
	}
	log.Log.Info("Created the pod")
	time.Sleep(5 * time.Second)

	answer, err := readPodLogs(pod)
	if err != nil {
		log.Log.Error(err, "unable to read pod logs")
		return ctrl.Result{}, err
	}

	log.Log.Info(fmt.Sprintf("Answer: %s", answer))
	problem.Status.Answer = answer

	if err := r.Status().Update(ctx, &problem); err != nil {
		log.Log.Error(err, "unable to update Problem status")
		return ctrl.Result{}, err
	}

	if err := r.Delete(ctx, &pod, &client.DeleteOptions{}); err != nil {
		log.Log.Error(err, "unable to delete pod")
		return ctrl.Result{}, err
	}
	log.Log.Info("Deleted the pod")

	return ctrl.Result{}, nil
}

func readPodLogs(pod corev1.Pod) (string, error) {
	config := ctrl.GetConfigOrDie()

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", err
	}

	req := clientSet.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &corev1.PodLogOptions{})

	podLogs, err := req.Stream(context.Background())
	if err != nil {
		return "", err
	}

	defer podLogs.Close()

	answer, err := io.ReadAll(podLogs)
	if err != nil {
		return "", err
	}

	return string(answer), nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ArithmeticReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mathv1.Arithmetic{}).
		Complete(r)
}
