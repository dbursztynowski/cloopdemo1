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

package controllers

import (
	"context"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	closedlooppoocv1 "closedloop/api/v1"
)

// Monitoringv2Reconciler reconciles a Monitoringv2 object
type MonitoringDv2Reconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=closedloops,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=closedloops/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=closedloops/finalizers,verbs=update
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=monitorings,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=monitorings/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=monitorings/finalizers,verbs=update
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=decisions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=decisions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=decisions/finalizers,verbs=update
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=executions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=executions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=executions/finalizers,verbs=update
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=monitoringv2s,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=monitoringv2s/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=monitoringv2s/finalizers,verbs=update
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=executionv2s,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=executionv2s/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=executionv2s/finalizers,verbs=update
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=decisionv2s,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=decisionv2s/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=decisionv2s/finalizers,verbs=update
//+kubebuilder:rbac:groups="apps",resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="apps",resources=deployments/status,verbs=get;watch;list
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=services/status,verbs=get;watch;list
//+kubebuilder:rbac:groups="networking.k8s.io",resources=ingresses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="networking.k8s.io",resources=ingresses/status,verbs=get;watch;list

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Monitoringv2 object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *MonitoringDv2Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	verbosityLog := new (VerbosityLog)
	verbosityLog.SetMaxLevel(1)
	l := verbosityLog.FromContext(ctx)
	l.V(1).Info("D>>>>>>>>>>>>>>>>>>>>>>>>>")
	l.V(1).Info("Enter Reconcile MonitoringD")

	//Retreiving ClosedLoop Object who triggered the Reconciler

	MonitoringD := &closedlooppoocv1.MonitoringDv2{}
	r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, MonitoringD)
	// Update Status if it's not already done
	// Here we don't use Status for our logic, so the section is not needed but only to show how to use it
	if MonitoringD.Name != MonitoringD.Status.Affix {
		MonitoringD.Status.Affix = MonitoringD.Name
		//closedLoop.Status.ContextMgr.Spec.DeploymentConfig.DeployReplicas = 4
		if err := r.Status().Update(ctx, MonitoringD); err != nil {
			l.Error(err, "Failed to update Monitoring status")
			return ctrl.Result{}, err
		}
		l.V(2).Info("Enter Reconcile View Spec & Status", "spec", MonitoringD.Spec, "status", MonitoringD.Status)
	}

	/* ---------------------------------- START Monitoring Part ---------------------------------- */
	DecisionD := &closedlooppoocv1.DecisionD{}
	r.Get(ctx, types.NamespacedName{Name: MonitoringD.Spec.Affix + "-decisiond", Namespace: MonitoringD.Namespace}, DecisionD)

	//ewma_data := make(map[string]string)
	ewma_data := DecisionD.Spec.Data
	message := ""
	fmt.Printf("%+v", ewma_data)
	new_ewma_data := MonitoringD.Spec.Data
	//new_ewma_data = ewma_data
	// Treatment of data in the "Data" field

	/* ---------------------------------- END Monitoring Part ---------------------------------- */

	l.V(2).Info("message Send to DecisionD " + message)

	/* -------------------------------- Apply modification on Decision --------------------- */

	if message != "" && string(message[0]) == "," {
		message = strings.TrimPrefix(message, string(message[0]))
	} else {
		message = "No Event"
	}

	switch MonitoringD.Spec.DecisionKind {
	case "DecisionD":
		if err := r.ApplyDecision(ctx, MonitoringD, l, message, new_ewma_data); err != nil {
			l.Error(err, "Failed to Apply DecisionD")
			return ctrl.Result{}, err
		}

	}

	return ctrl.Result{}, nil
}

func (r *MonitoringDv2Reconciler) ApplyDecision(ctx context.Context, monitoringd *closedlooppoocv1.MonitoringDv2, l VerbosityLog, Message string, data map[string]string) error {

	DecisionD := &closedlooppoocv1.DecisionD{}
	r.Get(ctx, types.NamespacedName{Name: monitoringd.Spec.Affix + "-decisiond", Namespace: monitoringd.Namespace}, DecisionD)

	l.V(2).Info("Update Message on DecisionD " + DecisionD.Name)
	//DecisionD.Spec.Message = Message
	DecisionD.Spec.Time = monitoringd.Spec.Time

	//readData := make(map[string]string)
	newData := make(map[string]string)
	readData := DecisionD.Spec.Data
	mtime := monitoringd.Spec.Time
	indeks := data["metric"]
	l.V(2).Info(indeks)
	l.V(2).Info(fmt.Sprintf("readData %#v", readData))

	//arrData := make([]string)
	arrData := strings.Split(readData[indeks],";")
	l.V(2).Info(fmt.Sprintf("arrData %#v", arrData))
	if len(arrData)!=0 && mtime != arrData[len(arrData)-1] {
		if len(arrData) < 10 {
			arrData = append(arrData, mtime)
		} else {
			arrData = arrData[1:]
			arrData = append(arrData, mtime)
		}
	}
	newvalue := strings.Trim(strings.TrimSpace(strings.Join(arrData, ";")),";")
	l.V(2).Info(fmt.Sprintf("arrData %#v", arrData))
	newData["cpu"] = readData["cpu"] // strings.Join(arrData, ";")
	newData["memory"] = readData["memory"]
	if (indeks == "cpu" || indeks == "memory") {
		newData[indeks] = newvalue
	}
	DecisionD.Spec.Data = newData
	l.V(1).Info(fmt.Sprintf("send message to DecisionD %#v", newData))

	//return nil
	return r.Update(ctx, DecisionD)

}

// SetupWithManager sets up the controller with the Manager.
func (r *MonitoringDv2Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&closedlooppoocv1.MonitoringDv2{}).
		Complete(r)
}
