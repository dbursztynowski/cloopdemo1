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
	"time"

	//"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	//"sigs.k8s.io/controller-runtime/pkg/log"

	closedlooppoocv1 "closedloop/api/v1"
)

// DecisionReconciler reconciles a Decision object
type DecisionDReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// Here you give all the permission your controller need to be able to work

//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=closedloopds,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=closedloopds/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=closedloopds/finalizers,verbs=update
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=monitoringds,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=monitoringds/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=monitoringds/finalizers,verbs=update
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=decisionds,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=decisionds/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=decisionds/finalizers,verbs=update
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=executionds,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=executionds/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=executionds/finalizers,verbs=update
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=monitoringdv2s,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=monitoringdv2s/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=monitoringdv2s/finalizers,verbs=update
//+kubebuilder:rbac:groups="apps",resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="apps",resources=deployments/status,verbs=get;watch;list
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=services/status,verbs=get;watch;list
//+kubebuilder:rbac:groups="networking.k8s.io",resources=ingresses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="networking.k8s.io",resources=ingresses/status,verbs=get;watch;list

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Decision object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *DecisionDReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	verbosityLog := new (VerbosityLog)
	verbosityLog.SetMaxLevel(1)
	l := verbosityLog.FromContext(ctx)
	
	l.V(1).Info("Enter Reconcile DecisionD")

	//Retreiving Decision Object who triggered the Reconciler

	DecisionD := &closedlooppoocv1.DecisionD{}
	r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, DecisionD)

	// Update Status if it's not already done
	// Here we don't use Status for our logic, so the section is not needed but only to show how to use it
	if DecisionD.Name != DecisionD.Status.Affix {
		DecisionD.Status.Affix = DecisionD.Name

		if err := r.Status().Update(ctx, DecisionD); err != nil {
			l.Error(err, "Failed to update DecisionD status")
			return ctrl.Result{}, err
		}
		l.V(2).Info("Enter Reconcile View Spec & Status", "spec", DecisionD.Spec, "status", DecisionD.Status)
	}

	// START of data processing in the "message" field

	//readData := make(map[string]string)
	readData := DecisionD.Spec.Data
	arrDataCpu := strings.Split(readData["cpu"], ";")
	arrDataMemory := strings.Split(readData["memory"], ";")
	newincreaserank := ""
	closedLoop := &closedlooppoocv1.ClosedLoop{}
	r.Get(ctx, types.NamespacedName{Name: "closedloop-v2", Namespace: DecisionD.Namespace}, closedLoop)

	increasetime := closedLoop.Status.IncreaseTime
	increaserank := closedLoop.Status.IncreaseRank
	if (arrDataMemory[0] != "" && arrDataCpu[0] != "" && (len(arrDataMemory)==10 || len(arrDataCpu)==10)) {
		if increasetime < time.Now().Add(time.Duration(-1)*time.Minute).String() {
			if len(arrDataMemory) < len(arrDataCpu) || arrDataMemory[0] < arrDataCpu[0] && increaserank != "memory" {
				newincreaserank = "memory"
				increasetime = arrDataCpu[len(arrDataCpu)-1]
				l.V(1).Info(fmt.Sprintf("%s %s %s", arrDataMemory[0], time.Now().Add(time.Duration(-1)*time.Minute).String(), increaserank))
			} else if len(arrDataMemory) > len(arrDataCpu) || arrDataMemory[0] > arrDataCpu[0] && increaserank != "cpu"  {
				newincreaserank = "cpu"
				increasetime = arrDataMemory[len(arrDataMemory)-1]
				l.V(1).Info(fmt.Sprintf("%s %s %s", arrDataCpu[0], time.Now().Add(time.Duration(-1)*time.Minute).String(), increaserank))
			}
		}
	}
	if newincreaserank != "" && closedLoop.Status.IncreaseRank != newincreaserank{
		l.V(2).Info("Previous value " + closedLoop.Status.IncreaseRank)
/*			closedLoop.Status.IncreaseRank = newincreaserank
			closedLoop.Status.IncreaseTime = increasetime
			l.V(2).Info("New value " + closedLoop.Status.IncreaseRank)
			l.V(1).Info("Send message to ClosedLoop: " + "IncreaseRank=" + newincreaserank + " IncreaseTime=" + increasetime)

			r.Status().Update(ctx, closedLoop)
*/
			if err := r.ApplyExecution(ctx, DecisionD, l, "Increase rank", newincreaserank); err != nil {
				l.V(2).Error(err, "Failed to ApplyExecutionD")
				return ctrl.Result{}, err
			}
	} else {
		if err := r.ApplyExecution(ctx, DecisionD, l, "No action", "No meric"); err != nil {
			l.V(2).Error(err, "Failed to ApplyExecutionD")
			return ctrl.Result{}, err
		}

	}

	return ctrl.Result{}, nil
}

// Function to update Execution Action
func (r *DecisionDReconciler) ApplyExecution(ctx context.Context, decision *closedlooppoocv1.DecisionD, 
		l VerbosityLog, Action string, Metric string) error {
	// Try to retrieve the CR that we want to update
	ExecutionD := &closedlooppoocv1.ExecutionD{}
	r.Get(ctx, types.NamespacedName{Name: decision.Spec.Affix + "-executiond", Namespace: decision.Namespace}, ExecutionD)

	l.V(2).Info("Update Action on ExecutionD")
	//Update it's field with the variable message
	ExecutionD.Spec.Action = Action
	ExecutionD.Spec.Metric = Metric	
	ExecutionD.Spec.Time = decision.Spec.Time
	l.V(1).Info("Send message to ExecutionD " +	"Action=" + Action + " Metric=" + Metric + " Time=" + decision.Spec.Time)

	return r.Update(ctx, ExecutionD)

}

// SetupWithManager sets up the controller with the Manager.
func (r *DecisionDReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&closedlooppoocv1.DecisionD{}).
		Complete(r)
}
