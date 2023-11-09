/*
Copyright 2023.

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

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	webappv1alpha1 "arpan.dev/kubecon-us-2019/api/v1alpha1"
)

// GuestBookReconciler reconciles a GuestBook object
type GuestBookReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=webapp.arpan.dev,resources=guestbooks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=webapp.arpan.dev,resources=guestbooks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=webapp.arpan.dev,resources=guestbooks/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the GuestBook object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.0/pkg/reconcile
func (r *GuestBookReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("guestbook", req.NamespacedName)
	log.Info("Starting Reconcile", req.NamespacedName)

	guestbook := webappv1alpha1.GuestBook{}
	if err := r.Get(ctx, req.NamespacedName, &guestbook); err != nil {
		if errors.IsNotFound(err) {
			log.Info("Unable to Find The Object", req.NamespacedName)
			return ctrl.Result{}, nil
		} else if err != nil {
			return ctrl.Result{}, err
		}
	}
	var redis webappv1alpha1.Redis

	err := r.Get(ctx, types.NamespacedName{Name: guestbook.Spec.RedisName, Namespace: req.Namespace}, &redis)

	if err != nil && errors.IsNotFound(err) {
		log.Info("Unable to find The Redis Object with name", guestbook.Spec.RedisName)
		return ctrl.Result{}, nil
	} else if err != nil {
		log.Info("Unable to get Data Check Api and Permissions")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GuestBookReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webappv1alpha1.GuestBook{}).
		Complete(r)
}
