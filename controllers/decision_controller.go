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
	"strconv"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	closedlooppoocv1 "closedloop/api/v1"
)

// DecisionReconciler reconciles a Decision object
type DecisionReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// Here you give all the permission your controller need to be able to work

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
// the Decision object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *DecisionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	verbosityLog := new (VerbosityLog)
	verbosityLog.SetMaxLevel(1)
	l := verbosityLog.FromContext(ctx)
	l.V(2).Info("\n********************************")
	l.V(1).Info("Enter Reconcile Decision")

	//Retreiving Decision Object who triggered the Reconciler

	Decision := &closedlooppoocv1.Decision{}
	r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, Decision)

	// Update Status if it's not already done
	// Here we don't use Status for our logic, so the section is not needed but only to show how to use it
	if Decision.Name != Decision.Status.Affix {
		Decision.Status.Affix = Decision.Name

		if err := r.Status().Update(ctx, Decision); err != nil {
			l.Error(err, "Failed to update Decision status")
			return ctrl.Result{}, err
		}
		l.V(2).Info("Update Decision Spec & Status", "spec", Decision.Spec, "status", Decision.Status)
	}

	// START of data processing in the "message" field

	//	fmt.Println("Message received from Monitoring: " + Decision.Spec.Message + " | timestamp: " + Decision.Spec.Time)
	l.V(1).Info("Message received from Monitoring: " + Decision.Spec.Message + " | timestamp: " + Decision.Spec.Time)
	if Decision.Spec.Message == "No Message" || Decision.Spec.Message == "" {
		l.V(1).Info("No Data to process")
		return ctrl.Result{}, nil
	}

	// Convert the string into an array of values
	values := strings.Split(Decision.Spec.Message, ",")

	// Check for the presence of the required values

	searchValues := make([]string, len(Decision.Spec.DecisionPolicies.PrioritySpec.PriorityRank) + 1)
    searchValues[0] ="No Event"


	switch Decision.Spec.DecisionPolicies.DecisionType {
	case "Priority":
		for keyStr, value := range Decision.Spec.DecisionPolicies.PrioritySpec.PriorityRank {
			//ranks = append(ranks, key)
			key, _ := strconv.ParseInt(keyStr[5:],10,0)
		    searchValues[key] = value	

		}
		// If we had more DecisionType  we would have to add some other cases like below
		/*
			case "Something_Else":
				//Other Logic here
		*/

	}

	for _, searchValue := range searchValues {
		found := false
		for _, value := range values {
			if strings.Contains(value, searchValue) {
				found = true
				break
			}
		}
		if found {
			switch Decision.Spec.DecisionPolicies.DecisionType {
			case "Priority":
				for _, rankValue := range Decision.Spec.DecisionPolicies.PrioritySpec.PriorityRank {

					// here, the conversion of the decision onto a real action to be done by Execution should take place
					if searchValue == rankValue {
						action := "React to " + rankValue
						//						fmt.Println(action)
						l.Info(action)
						switch Decision.Spec.ExecutionKind {
						case "Execution":
							if err := r.ApplyExecution(ctx, Decision, l, action); err != nil {
								l.Error(err, "Failed to ApplyExecution")
								return ctrl.Result{}, err
							}
							if err := r.InformDeliberate(ctx, Decision, l, rankValue); err != nil {
								l.Error(err, "Failed to Inform Deliberate")
								return ctrl.Result{}, err
						}

						}
						return ctrl.Result{}, nil

					}
				}
				// If we had more DecisionType  we would have to add some other cases like below
				/*
					case "Something_Else":
						//Other Logic here
				*/

			}

		}
	}

	// END of data processing in the "message" field

	return ctrl.Result{}, nil
}

// Function to update Execution Action
func (r *DecisionReconciler) ApplyExecution(ctx context.Context, decision *closedlooppoocv1.Decision, l VerbosityLog, Action string) error {
	// Try to retrieve the CR that we want to update

	if decision.Spec.Time == decision.Spec.DecisionPolicies.PrioritySpec.Time {
		return nil
	} 
	Execution := &closedlooppoocv1.Execution{}
	r.Get(ctx, types.NamespacedName{Name: decision.Spec.Affix + "-execution", Namespace: decision.Namespace}, Execution)

	l.V(2).Info("Update Action on Execution")
	//Update it's field with the variable message
	Execution.Spec.Action = Action
	Execution.Spec.Time = decision.Spec.Time
	l.V(1).Info("Send message to Execution: " + Execution.Spec.Action + " " + Execution.Spec.Time)

	return r.Update(ctx, Execution)

}

func (r *DecisionReconciler) InformDeliberate(ctx context.Context, decision *closedlooppoocv1.Decision, l VerbosityLog, Message string) error {

	MonitoringDv2 := &closedlooppoocv1.MonitoringDv2{}
	r.Get(ctx, types.NamespacedName{Name: "closedloopd-v2-monitoringd", Namespace: decision.Namespace}, MonitoringDv2)
	
	if (MonitoringDv2.Spec.Time != decision.Spec.Time) {
		l.V(1).Info("Update Message on MonitoringDv2 " + MonitoringDv2.Name)
		l.V(2).Info("Data before " + fmt.Sprint(MonitoringDv2.Spec.Data))
		m := make(map[string]string)
		m["metric"] = Message
		MonitoringDv2.Spec.Data = m
		l.V(1).Info("Message send to ClosedLoop " + fmt.Sprint(MonitoringDv2.Spec.Data))
		MonitoringDv2.Spec.Time = decision.Spec.Time
		return r.Update(ctx, MonitoringDv2)
	}
	return nil

}


// SetupWithManager sets up the controller with the Manager.
func (r *DecisionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&closedlooppoocv1.Decision{}).
		Complete(r)
}
