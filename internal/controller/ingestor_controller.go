/*
Copyright 2025.

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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	utilv1alpha "github.com/end-2/api/v1alpha"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IngestorReconciler reconciles a Ingestor object
type IngestorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=util.github.com,resources=ingestors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=util.github.com,resources=ingestors/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=util.github.com,resources=ingestors/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Ingestor object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.20.4/pkg/reconcile
func (r *IngestorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := logf.FromContext(ctx)

	logger.Info("start reconciling ingestor")
	defer func() {
		logger.Info("done reconciling ingestor")
	}()

	ingestor := &utilv1alpha.Ingestor{}
	if err := r.Get(ctx, req.NamespacedName, ingestor); err != nil {
		logger.Error(err, "failed to get ingestor")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	logger.Info("ingestor", "name", ingestor.Name, "namespace", ingestor.Namespace)

	configMapList := &corev1.ConfigMapList{}
	if err := r.List(ctx, configMapList, client.InNamespace(ingestor.Namespace)); err != nil {
		logger.Error(err, "failed to list configmaps")
		return ctrl.Result{}, err
	}

	logger.Info("configmaps", "count", len(configMapList.Items))

	for _, configMap := range configMapList.Items {
		logger.Info("configmap", "name", configMap.Name, "namespace", configMap.Namespace)
		for key, value := range configMap.Data {
			logger.Info("configmap", "key", key, "value", value)
		}
	}

	secretList := &corev1.SecretList{}
	if err := r.List(ctx, secretList, client.InNamespace(ingestor.Namespace)); err != nil {
		logger.Error(err, "failed to list secrets")
		return ctrl.Result{}, err
	}

	logger.Info("secrets", "count", len(secretList.Items))

	for _, secret := range secretList.Items {
		logger.Info("secret", "name", secret.Name, "namespace", secret.Namespace)
		for key, value := range secret.Data {
			logger.Info("secret", "key", key, "value", value)
		}
	}

	ingestor.Status.LastIngested = metav1.Now()
	if err := r.Status().Update(ctx, ingestor); err != nil {
		logger.Error(err, "failed to update ingestor status")
		return ctrl.Result{}, err
	}

	if ingestor.Spec.Interval.Duration > 0 {
		return ctrl.Result{RequeueAfter: ingestor.Spec.Interval.Duration}, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *IngestorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&utilv1alpha.Ingestor{}).
		Named("ingestor").
		Watches(&corev1.ConfigMap{}, handler.EnqueueRequestForOwner(mgr.GetScheme(), mgr.GetRESTMapper(), &utilv1alpha.Ingestor{})).
		Watches(&corev1.Secret{}, handler.EnqueueRequestForOwner(mgr.GetScheme(), mgr.GetRESTMapper(), &utilv1alpha.Ingestor{})).
		// example of Owns
		// Owns(&corev1.ConfigMap{}).
		// Owns(&corev1.Secret{}).
		Complete(r)
}
