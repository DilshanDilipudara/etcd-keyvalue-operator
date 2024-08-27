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
	"reflect"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	etcdv1 "github.com/DilshanDilipudara/etcd-keyvalue-operator/api/v1"
	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EtcdConfigReconciler reconciles a EtcdConfig object
type EtcdConfigReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=etcd.dilshan.com,resources=etcdconfigs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=etcd.dilshan.com,resources=etcdconfigs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=etcd.dilshan.com,resources=etcdconfigs/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the EtcdConfig object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.3/pkg/reconcile

func (r *EtcdConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	//_ = log.FromContext(ctx)

	// create log instance
	log := r.Log.WithValues("etcdconfig", req.NamespacedName)
	log.Info("namespace name: ", "ns", req.Namespace)
	// Check if the req.Name ends with "-last-synced"
	if strings.HasSuffix(req.Name, "-last-synced") {
		log.Info("Skipping reconciliation for -last-synced resource", "resource", req.Name)
		return ctrl.Result{}, nil
	}

	// Fetch the EtcdConfig instance
	var etcdConfig etcdv1.EtcdConfig

	log.Info(req.Namespace)
	if err := r.Get(ctx, req.NamespacedName, &etcdConfig); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("EtcdConfig resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get EtcdConfig")
		return ctrl.Result{}, err
	}

	// Check if the configuration has changed
	var lastSyncedConfig etcdv1.EtcdConfig
	if err := r.Get(ctx, client.ObjectKey{Namespace: req.Namespace, Name: req.Name + "-last-synced"}, &lastSyncedConfig); err != nil && !apierrors.IsNotFound(err) {
		log.Error(err, "Failed to get last synced EtcdConfig")
		return ctrl.Result{}, err
	}

	// Only update the cluster if the configuration has changed
	if !reflect.DeepEqual(etcdConfig.Spec.Items, lastSyncedConfig.Spec.Items) {
		log.Info("Detected configuration change, updating Etcd cluster")
		// log.Info("etcdConfig: ","key: ",etcdConfig.Spec.Items)
		// log.Info("lastSync: ","key: ",lastSyncedConfig.Spec.Items)

		// Authenticate etcd cluster
		if err := getEtcdConfigFromSecret(ctx, r.Client); err != nil {
			log.Error(err, "Failed to authenticate with etcd")
			return ctrl.Result{}, err
		}

		// update the etcd config
		for _, configItem := range etcdConfig.Spec.Items {
			if err := updateEtcdCluster(configItem.Key, configItem.Value); err != nil {
				log.Error(err, "Failed to update Etcd cluster", "key", configItem.Key)
				return ctrl.Result{}, err
			}
		}

		// Check if lastSyncedConfig exists
		if lastSyncedConfig.Name != "" {
			if err := r.Delete(ctx, &lastSyncedConfig); err != nil {
				log.Error(err, "Failed to Delete last synced EtcdConfig")
				return ctrl.Result{}, err
			}
		}
		// Creating a new lastSyncedConfig
		lastSyncedConfig = etcdConfig
		lastSyncedConfig.Name = req.Name + "-last-synced"

		// Clear fields that should not be copied over to a new resource
		lastSyncedConfig.ResourceVersion = ""
		lastSyncedConfig.UID = ""
		lastSyncedConfig.CreationTimestamp = metav1.Time{}

		// Create the lastSyncedConfig
		if err := r.Create(ctx, &lastSyncedConfig); err != nil {
			log.Error(err, "Failed to create last synced EtcdConfig")
			return ctrl.Result{}, err
		}
		log.Info("Successfully created lastSyncedConfig", "resource", lastSyncedConfig.Name)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *EtcdConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.Log = ctrl.Log.WithName("ETCD-Controllers")
	return ctrl.NewControllerManagedBy(mgr).
		For(&etcdv1.EtcdConfig{}).
		Complete(r)
}
