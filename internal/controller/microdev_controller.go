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
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	kapps "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	testv1 "606.hdu.io/mycontroller/api/v1"
)

// MicroDevReconciler reconciles a MicroDev object
type MicroDevReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Clock
}

type realClock struct{}

func (c realClock) Now() time.Time { return time.Now() }

// clock knows how to get the current time.
// It can be used to fake out timing for testing.
type Clock interface {
	Now() time.Time
}

//+kubebuilder:rbac:groups=test.606.hdu.io,resources=microdevs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=test.606.hdu.io,resources=microdevs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=test.606.hdu.io,resources=microdevs/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MicroDev object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *MicroDevReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// TODO(user): your logic here
	//获取MicroDev资源
	var microDev testv1.MicroDev
	if err := r.Get(ctx, req.NamespacedName, &microDev); err != nil {
		log.Error(err, "unable to fetch MicroDev")
		return ctrl.Result{}, nil
	}

	//获取该MicroDev的管理的deployments
	//（当前设定一个deployment匹配一个MicroDev，通过Selector匹配获取）
	var deployments kapps.DeploymentList
	deploymentSelector, _ := metav1.LabelSelectorAsSelector(microDev.Spec.Selector)
	if err := r.List(ctx, &deployments, client.InNamespace(req.Namespace), client.MatchingLabelsSelector{Selector: deploymentSelector}); err != nil {
		log.Error(err, "unable to list child deployments")
		return ctrl.Result{}, err
	}

	//查询计划副本数量
	getReplicas := func(deploymentName, queryUrl string) (int32, error) {
		if queryUrl == "" {
			return 0, nil
		}
		hasSlash := strings.HasSuffix(queryUrl, "/")
		var url string
		if hasSlash {
			url = queryUrl + deploymentName
		} else {
			url = queryUrl + "/" + deploymentName
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Error(err, "Request canceled")
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error(err, "read Response Body error")
		}

		type Response struct {
			Replicas int32 `json:"replicas" `
		}

		var res Response
		json.Unmarshal(body, &res)
		return res.Replicas, err
	}

	var scheduledResult ctrl.Result
	//使用MicroDev中spec里填写的QueryUrl查询当前deployment是否需要更改副本数
	for _, deployment := range deployments.Items {
		name := deployment.ObjectMeta.Name
		queryUrl := microDev.Spec.QueryUrl
		currentReplicas := deployment.Spec.Replicas
		expectedReplicas, err := getReplicas(name, queryUrl)
		if err != nil {
			return ctrl.Result{}, err
		}

		//若计划副本数与当前副本数相同，则跳过更新，按MicroDev中spec里填写的时间间隔计划重新入队
		if *currentReplicas == expectedReplicas {
			queryInterval := *microDev.Spec.QueryInterval
			scheduledResult = ctrl.Result{RequeueAfter: time.Duration(queryInterval) * time.Second}
			continue
		}

		//若查询失败，则在30s后重新入队
		if expectedReplicas == 0 {
			queryInterval := 30
			scheduledResult = ctrl.Result{RequeueAfter: time.Duration(queryInterval) * time.Second}
			continue
		}

		//更新microDev的状态以及调整deployment的副本数
		microDev.Status.Replicas = &expectedReplicas
		if err := r.Status().Update(ctx, &microDev); err != nil {
			log.Error(err, "unable to update microDev status")
			return ctrl.Result{}, err
		}

		deployment.Spec.Replicas = &expectedReplicas
		if err := r.Update(ctx, &deployment); client.IgnoreNotFound(err) != nil {
			log.Error(err, "unable to update old deployment", "deployment", deployment)
			return ctrl.Result{}, err
		} else {
			log.V(0).Info("updated old deployment", "deployment", deployment)
		}
	}

	return scheduledResult, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MicroDevReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// set up a real clock, since we're not in a test
	if r.Clock == nil {
		r.Clock = realClock{}
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&testv1.MicroDev{}).
		Complete(r)
}
