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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"bytes"
	"encoding/json"
	"net/http"

	"github.com/go-logr/logr"
	greenopsv1 "github.com/vicentinileonardo/operator-workload-placement-policy/api/v1"
)

const (
	REGION_FILTERING_SERVER_BASE_URL   = "region-filtering-server"
	REGION_FILTERING_SERVER_K8S_SUFFIX = ".default.svc.cluster.local"
	REGION_FILTERING_SERVER_PORT       = "8080"
	REGION_FILTERING_SERVER_ENDPOINT   = "/regions/eligible"

	LOCAL_REGION_FILTERING_SERVER_BASE_URL = "localhost"
	LOCAL_REGION_FILTERING_SERVER_PORT     = "8080"
	LOCAL_REGION_FILTERING_SERVER_ENDPOINT = "/regions/eligible"
)

type requestPayload struct {
	CloudProviderOriginRegion string `json:"cloudProviderOriginRegion"`
	MaxLatency                int    `json:"maxLatency"`
	CloudProvider             string `json:"cloudProvider"`
}

type responsePayload struct {
	CloudProvider   string              `json:"cloudProvider"`
	EligibleRegions []greenopsv1.Region `json:"eligibleRegions"`
}

// WorkloadPlacementPolicyReconciler reconciles a WorkloadPlacementPolicy object
type WorkloadPlacementPolicyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=greenops.greenops.test,resources=workloadplacementpolicies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=greenops.greenops.test,resources=workloadplacementpolicies/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=greenops.greenops.test,resources=workloadplacementpolicies/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the WorkloadPlacementPolicy object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *WorkloadPlacementPolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	l.Info("Starting WorkloadPlacementPolicy reconciliation")

	wpp := &greenopsv1.WorkloadPlacementPolicy{}
	if err := r.Get(ctx, req.NamespacedName, wpp); err != nil {
		l.Error(err, "unable to fetch WorkloadPlacementPolicy")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	//fetch eligible regions with a http post request
	eligibleRegions := fetchEligibleRegions(wpp.Spec.OriginRegion, wpp.Spec.MaxLatency, wpp.Spec.CloudProvider, l)

	l.Info("Eligible regions", "regions", eligibleRegions)

	// fill status with eligible regions
	wpp.Status.EligibleRegions = eligibleRegions

	if err := r.Status().Update(ctx, wpp); err != nil {
		l.Error(err, "unable to update WorkloadPlacementPolicy status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// fetchEligibleRegions fetches eligible regions from an external service
func fetchEligibleRegions(originRegion greenopsv1.Region, maxLatency int, cloudProvider string, l logr.Logger) []greenopsv1.Region {

	payload := requestPayload{
		CloudProviderOriginRegion: originRegion.CloudProviderRegion,
		MaxLatency:                maxLatency,
		CloudProvider:             cloudProvider,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		// handle error
		return []greenopsv1.Region{}
	}

	l.Info("[PAYLOAD]", "payload", string(jsonPayload))

	//url := fmt.Sprintf("http://%s%s:%s%s", REGION_FILTERING_SERVER_BASE_URL, REGION_FILTERING_SERVER_K8S_SUFFIX, REGION_FILTERING_SERVER_PORT, REGION_FILTERING_SERVER_ENDPOINT)

	//local url for testing
	url := fmt.Sprintf("http://%s:%s%s", LOCAL_REGION_FILTERING_SERVER_BASE_URL, LOCAL_REGION_FILTERING_SERVER_PORT, LOCAL_REGION_FILTERING_SERVER_ENDPOINT)

	l.Info("[URL]", "url", url)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		// handle error
		l.Info("[ERROR]", "error1", err)
		return []greenopsv1.Region{}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// handle error
		l.Info("[ERROR]", "error2", resp.StatusCode)
		return []greenopsv1.Region{}
	}

	l.Info("[STATUS]", "status", resp.StatusCode)

	var response responsePayload
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		// handle error
		l.Info("[ERROR]", "error3", err)
		return []greenopsv1.Region{}
	}

	l.Info("[RESPONSE]", "response", response)

	return response.EligibleRegions

}

// SetupWithManager sets up the controller with the Manager.
func (r *WorkloadPlacementPolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&greenopsv1.WorkloadPlacementPolicy{}).
		Named("workloadplacementpolicy").
		Complete(r)
}
