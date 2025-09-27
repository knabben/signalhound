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
	"fmt"
	"reflect"
	"time"

	"github.com/go-logr/logr"
	testgridv1alpha1 "github.com/knabben/signalhound/api/v1alpha1"
	"github.com/knabben/signalhound/internal/testgrid"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// DashboardReconciler reconciles a Dashboard object
type DashboardReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	log    logr.Logger
}

// +kubebuilder:rbac:groups=testgrid.holdmybeer.io,resources=dashboards,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=testgrid.holdmybeer.io,resources=dashboards/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=testgrid.holdmybeer.io,resources=dashboards/finalizers,verbs=update

// Reconcile loops against the dashboard reconciler and set the final object status.
func (r *DashboardReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.log = logf.FromContext(ctx).WithValues("resource", req.NamespacedName)

	var dashboard testgridv1alpha1.Dashboard
	if err := r.Get(ctx, req.NamespacedName, &dashboard); err != nil {
		r.log.Error(err, "unable to fetch dashboard")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	grid := testgrid.NewTestGrid(testgrid.URL)
	summary, err := grid.FetchSummary(dashboard.Spec.DashboardTab, testgridv1alpha1.ERROR_STATUSES)
	if err != nil {
		r.log.Error(err, "error fetching summary from endpoint.")
		return ctrl.Result{}, err
	}

	if r.shouldRefresh(dashboard.Status, summary) {
		// set the dashboard summary on status if an update happened
		dashboard.Status.DashboardSummary = summary
		dashboard.Status.LastUpdate = metav1.Now()
		r.log.Info("updating dashboard object status.")
		if err := r.Status().Update(ctx, &dashboard); err != nil {
			r.log.Error(err, "unable to update dashboard status")
			return ctrl.Result{}, err
		}

		err = r.createOrUpdateConfigmaps(ctx, req, grid, dashboard)
		if err != nil {
			r.log.Error(err, "unable to create update a configmap")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// createOrUpdateConfigmaps
func (r *DashboardReconciler) createOrUpdateConfigmaps(
	ctx context.Context,
	req ctrl.Request,
	grid *testgrid.TestGrid,
	dashboard testgridv1alpha1.Dashboard,
) (err error) {
	for _, summary := range dashboard.Status.DashboardSummary {
		var configmap = &v1.ConfigMap{}

		table, err := grid.FetchTable(summary, dashboard.Spec.MinFlakes, dashboard.Spec.MinFailures)
		if err != nil {
			r.log.Error(err, "error fetching table", "tab", summary.TabName)
			continue
		}

		configmapName := fmt.Sprintf("%s-%s", summary.DashboardName, summary.TabName)
		configMapKey := client.ObjectKey{Namespace: req.Namespace, Name: configmapName}

		if err := r.Get(ctx, configMapKey, configmap); err != nil {
			if apierrors.IsNotFound(err) {
				configmap = &v1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      configmapName,
						Namespace: req.Namespace,
					},
					Data: map[string]string{
						"description": table.Description,
					},
				}
				// create the configmap
				if err := r.Create(ctx, configmap); err != nil {
					return err
				}
				r.log.Info("created configmap", "configmap", configmapName)
			}
		}
	}
	return nil
}

// shouldRefresh determines if it's time to refresh the dashboard data
func (r *DashboardReconciler) shouldRefresh(dashboardStatus testgridv1alpha1.DashboardStatus, summary []testgridv1alpha1.DashboardSummary) bool {
	if reflect.DeepEqual(dashboardStatus.DashboardSummary, summary) {
		return false
	}

	if dashboardStatus.LastUpdate.IsZero() {
		return true
	}

	refreshInterval := time.Duration(1) * time.Minute // should at least wait for 1 minute
	return time.Since(dashboardStatus.LastUpdate.Time) >= refreshInterval
}

// SetupWithManager sets up the controller with the Manager.
func (r *DashboardReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&testgridv1alpha1.Dashboard{}).
		Named("dashboard").
		Complete(r)
}
