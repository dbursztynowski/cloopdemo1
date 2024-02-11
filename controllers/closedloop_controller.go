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
Version 2
*/

package controllers

import (
	"context"
	"sort"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	closedlooppoocv1 "closedloop/api/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClosedLoopReconciler reconciles a ClosedLoop object
type ClosedLoopReconciler struct {
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
// the ClosedLoop object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *ClosedLoopReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	verbosityLog := new (VerbosityLog)
	verbosityLog.SetMaxLevel(1)
	l := verbosityLog.FromContext(ctx)
	l.V(1).Info("Enter Reconcile ClosedLoop")

	//Retreiving ClosedLoop Object who triggered the Reconciler
	closedLoop := &closedlooppoocv1.ClosedLoop{}
	err := r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, closedLoop)

	// Update Status if it's not already done
	// Here we don't use Status for our logic, so the section is not needed but only to show how to use it
	updateStatus := false
	if closedLoop.Name != closedLoop.Status.Name {
		closedLoop.Status.Name = closedLoop.Name
		updateStatus = true
	}

	if closedLoop.Status.IncreaseRank == "" {
		l.V(2).Info("Set start")
		closedLoop.Status.IncreaseRank = "start"
		closedLoop.Status.IncreaseTime = time.Now().String()
		updateStatus = true
	}

	if updateStatus {
		if err := r.Status().Update(ctx, closedLoop); err != nil {
			l.Error(err, "Failed to update closedLoop status")
			return ctrl.Result{}, err
		}
		l.V(2).Info("Enter Reconcile View Spec & Status", "spec", closedLoop.Spec, "status", closedLoop.Status)
	}

	// If the closedLoop who triggered is not find it means that it's been deleted
	if err != nil {

		// The Closedloop object has been deleted or is not found, so we should delete the associated CR (if they exist)
		l.Info("Close loop Instance not found, Deletion Close Loop ressources")

		if err := r.deleteClosedLoopMonitoringType1(ctx, req.Name, req.Namespace, l); err != nil {
			l.Error(err, "No Monitoring found, no deletions have been made")

		}
		if err := r.deleteClosedLoopMonitoringType2(ctx, req.Name, req.Namespace, l); err != nil {
			l.Error(err, "No Monitoringv2 found, no deletions have been made")

		}

		if err := r.deleteClosedLoopDecisionType1(ctx, req.Name, req.Namespace, l); err != nil {
			l.Error(err, "No Decision type 1 found, no deletions have been made")

		}

		if err := r.deleteClosedLoopExecutionType1(ctx, req.Name, req.Namespace, l); err != nil {
			l.Error(err, "No Execution type 1 found, no deletions have been made")

		}

		//l.Info()
		return ctrl.Result{}, nil

	}
	updateSpec := false
	betterrank := closedLoop.Status.IncreaseRank
	l.V(1).Info("Closedloop receive message: " + "increase rank " + betterrank)
	//betterrankkey := ""
	betterrankindex := -1
	m := make(map[string]string)
	tk := make(map[string]string)
	tv := make(map[string]string)
	for key, value := range closedLoop.Spec.Monitoring.MonitoringPolicies.Data {
			m[key] = value
			//if value == betterrank {
			//              betterrankkey = key
			//}
	}
	for key, value := range closedLoop.Spec.Monitoring.MonitoringPolicies.TresholdKind {
			tk[key] = value
	}
	for key, value := range closedLoop.Spec.Monitoring.MonitoringPolicies.TresholdValue {
			tv[key] = value
	}
	keys := make([]string, 0, len(m))

	for k := range m {
			keys = append(keys, k)
	}
	sort.Strings(keys)

	for i := 0; i < len(keys); i++ {
			if m[keys[i]] == betterrank {
					betterrankindex = i
			}
	}
	if betterrankindex > 0 {
			pom := m[keys[betterrankindex-1]]
			m[keys[betterrankindex-1]] = betterrank
			m[keys[betterrankindex]] = pom

			ptk := tk[keys[betterrankindex-1]+"-thresholdkind"]
			tk[keys[betterrankindex-1]+"-thresholdkind"] = tk[keys[betterrankindex]+"-thresholdkind"]
			tk[keys[betterrankindex]+"-thresholdkind"] = ptk

			ptv := tv[keys[betterrankindex-1]+"-thresholdvalue"]
			tv[keys[betterrankindex-1]+"-thresholdvalue"] = tv[keys[betterrankindex]+"-thresholdvalue"]
			tv[keys[betterrankindex]+"-thresholdvalue"] = ptv

			closedLoop.Spec.Monitoring.MonitoringPolicies.Data = m
			closedLoop.Spec.Monitoring.MonitoringPolicies.TresholdKind = tk
			closedLoop.Spec.Monitoring.MonitoringPolicies.TresholdValue = tv
			closedLoop.Spec.Monitoring.MonitoringPolicies.Time = closedLoop.Status.IncreaseTime
			updateSpec = true
	}

	betterrankindex = -1
	p := make(map[string]string)
	for key, value := range closedLoop.Spec.Decision.DecisionPolicies.PrioritySpec.PriorityRank {
			p[key] = value
	}
	pkeys := make([]string, 0, len(p))

	for k := range p {
			pkeys = append(pkeys, k)
	}
	sort.Strings(pkeys)

	for i := 0; i < len(pkeys); i++ {
			if p[pkeys[i]] == betterrank {
					betterrankindex = i
			}
	}
	if betterrankindex > 0 {
			pom := p[pkeys[betterrankindex-1]]
			p[pkeys[betterrankindex-1]] = betterrank
			p[pkeys[betterrankindex]] = pom

			closedLoop.Spec.Decision.DecisionPolicies.PrioritySpec.PriorityRank = p
			closedLoop.Spec.Decision.DecisionPolicies.PrioritySpec.Time = closedLoop.Status.IncreaseTime
			updateSpec = true
	}




	if updateSpec {
			if err := r.Update(ctx, closedLoop); err != nil {
					l.Error(err, "Failed to update closedLoop spec")
					return ctrl.Result{}, err
			}
			l.V(2).Info("Enter Reconcile View Spec & Status", "spec", closedLoop.Spec, "status", closedLoop.Status)
	}


	// Creation of CR implementation layer ressources (Second layer) are created based on Kinds defined in the primary closedloop ressource
	// These functions are called each time, so if they've already been created, nothing will happen.
	switch closedLoop.Spec.Execution.Kind {
	case "Execution":
		if err := r.createExecutionType1(ctx, closedLoop, l); err != nil {
			l.Error(err, "Failed to createExecutionType1")
			return ctrl.Result{}, err
		}
		// If we had more Execution Types/version we would have to add some other cases like below
		/*
			case "Executionv2":
				if err := r.createExecutionType2(ctx, closedLoop, l); err != nil {
					l.Error(err, "Failed to createExecutionType2")
					return ctrl.Result{}, err
				}
		*/

	}

	switch closedLoop.Spec.Decision.DecisionKind.DecisionKindName {
	case "Decision":
		if err := r.createDecisionType1(ctx, closedLoop, l); err != nil {
			l.Error(err, "Failed to createDecisionType1")
			return ctrl.Result{}, err
		}
		// If we had more Decision Types/version we would have to add some other cases like below
		/*
			case "Decisionv2":
				if err := r.createDecisionType2(ctx, closedLoop, l); err != nil {
					l.Error(err, "Failed to createDecisionType2")
					return ctrl.Result{}, err
				}
		*/

	}

	switch closedLoop.Spec.Monitoring.MonitoringKind.MonitoringKindName {
	case "Monitoring":
		if err := r.createMonitoringType1(ctx, closedLoop, l); err != nil {
			l.Error(err, "Failed to createMonitoringType1")
			return ctrl.Result{}, err
		}
	case "Monitoringv2":
		if err := r.createMonitoringType2(ctx, closedLoop, l); err != nil {
			l.Error(err, "Failed to createMonitoringType2")
			return ctrl.Result{}, err
		}
		// If we had more Monitoring Types/version we would have to add some other cases like below
		/*
			case "Monitoringv3":
				if err := r.createMonitoringType3(ctx, closedLoop, l); err != nil {
					l.Error(err, "Failed to createMonitoringType3")
					return ctrl.Result{}, err
				}
		*/

	}

	return ctrl.Result{}, nil
}

// Function use to create the CR Monitoring type 1 (kind : Monitoring)
func (r *ClosedLoopReconciler) createMonitoringType1(ctx context.Context, CL *closedlooppoocv1.ClosedLoop, l VerbosityLog) error {

	// Try to retrieve the CR to see if it's already within the Cluster
	clcm := &closedlooppoocv1.Monitoring{}
	err := r.Get(ctx, types.NamespacedName{Name: CL.Name + "-monitoring", Namespace: CL.Namespace}, clcm)

	if err == nil {
		l.V(2).Info("monitoring type 1 Found - No Creation")
		return nil
	}

	l.V(2).Info("Creating monitoring")

	//Creating the Monitoring Object with the right Spec
	clcm = &closedlooppoocv1.Monitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: CL.Namespace,
			Name:      CL.Name + "-monitoring",
		},
		Spec: closedlooppoocv1.MonitoringSpec{
			Affix:        CL.Name,
			DecisionKind: CL.Spec.Decision.DecisionKind.DecisionKindName,
			Source: closedlooppoocv1.Source{
				Addresse: CL.Spec.Monitoring.MonitoringKind.Source.Addresse,
				Port:     CL.Spec.Monitoring.MonitoringKind.Source.Port,
				Interval: CL.Spec.Monitoring.MonitoringKind.Source.Interval,
			},
			MonitoringPolicies: closedlooppoocv1.MonitoringPolicies{
				Data:          CL.Spec.Monitoring.MonitoringPolicies.Data,
				TresholdKind:  CL.Spec.Monitoring.MonitoringPolicies.TresholdKind,
				TresholdValue: CL.Spec.Monitoring.MonitoringPolicies.TresholdValue,
			},
		},
	}

	return r.Create(ctx, clcm)

}

// Function use to create the CR Monitoring type 2 (kind : Monitoringv2)
func (r *ClosedLoopReconciler) createMonitoringType2(ctx context.Context, CL *closedlooppoocv1.ClosedLoop, l VerbosityLog) error {

	// Try to retrieve the CR to see if it's already within the Cluster
	clcm := &closedlooppoocv1.Monitoringv2{}
	err := r.Get(ctx, types.NamespacedName{Name: CL.Name + "-monitoring", Namespace: CL.Namespace}, clcm)

	m := make(map[string]string)
	m["Value"] = "No Value"
	if err == nil {
		l.V(2).Info("monitoring type 2 Found - No Creation")
		clcm.Spec.MonitoringPolicies = CL.Spec.Monitoring.MonitoringPolicies
		if err := r.Update(ctx, clcm); err != nil {
				l.Error(err, "Failed to update policy Monitoring")
		}
		return nil
	}

	l.V(2).Info("Creating monitoring v2")
	l.V(2).Info("Value Decision kind : ", CL.Spec.Decision.DecisionKind.DecisionKindName)
	//Creating the Monitoringv2 Object with the right Spec
	clcm = &closedlooppoocv1.Monitoringv2{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: CL.Namespace,
			Name:      CL.Name + "-monitoring",
		},
		Spec: closedlooppoocv1.Monitoringv2Spec{
			Affix:        CL.Name,
			Time:         time.Now().String(),
			DecisionKind: CL.Spec.Decision.DecisionKind.DecisionKindName,
			RequestedPod: CL.Spec.Monitoring.MonitoringKind.RequestedPod,
			MonitoringPolicies: closedlooppoocv1.MonitoringPolicies{
				Data:          CL.Spec.Monitoring.MonitoringPolicies.Data,
				Time:	   	   CL.Spec.Monitoring.MonitoringPolicies.Time,
				TresholdKind:  CL.Spec.Monitoring.MonitoringPolicies.TresholdKind,
				TresholdValue: CL.Spec.Monitoring.MonitoringPolicies.TresholdValue,
			},
			Data: m,
		},
	}

	//Ask Kubernetes to create it
	return r.Create(ctx, clcm)

}

// Function use to create the CR Decision type 1 (kind : Decision)
func (r *ClosedLoopReconciler) createDecisionType1(ctx context.Context, CL *closedlooppoocv1.ClosedLoop, l VerbosityLog) error {

	// Try to retrieve the CR to see if it's already within the Cluster
	clcm := &closedlooppoocv1.Decision{}
	err := r.Get(ctx, types.NamespacedName{Name: CL.Name + "-decision", Namespace: CL.Namespace}, clcm)

	if err == nil {
		l.V(2).Info("Decision type 1Found - No Creation")
		clcm.Spec.DecisionPolicies.PrioritySpec = CL.Spec.Decision.DecisionPolicies.PrioritySpec
		if err := r.Update(ctx, clcm); err != nil {
				l.Error(err, "Failed to update policy Decision")
		}
		return nil
	}

	l.V(2).Info("Creating decision")
	//Creating the Decision Object with the right Spec
	clcm = &closedlooppoocv1.Decision{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: CL.Namespace,
			Name:      CL.Name + "-decision",
		},
		Spec: closedlooppoocv1.DecisionSpec{

			Affix:         CL.Name,
			Time:          time.Now().String(),
			Message:       "No Message",
			ExecutionKind: CL.Spec.Execution.Kind,

			DecisionPolicies: closedlooppoocv1.DecisionPolicies{
				DecisionType: CL.Spec.Decision.DecisionPolicies.DecisionType,
				PrioritySpec: closedlooppoocv1.PrioritySpec{
					PriorityType: CL.Spec.Decision.DecisionPolicies.PrioritySpec.PriorityType,
					PriorityRank: CL.Spec.Decision.DecisionPolicies.PrioritySpec.PriorityRank,
					Time:	   CL.Spec.Decision.DecisionPolicies.PrioritySpec.Time,
				},
			},
		},
	}

	//Ask Kubernetes to create it
	return r.Create(ctx, clcm)

}

// Function used to create the CR Execution type 1 (kind : Execution)
func (r *ClosedLoopReconciler) createExecutionType1(ctx context.Context, CL *closedlooppoocv1.ClosedLoop, l VerbosityLog) error {

	// Try to retrieve the CR to see if it's already within the Cluster
	clcm := &closedlooppoocv1.Execution{}
	err := r.Get(ctx, types.NamespacedName{Name: CL.Name + "-execution", Namespace: CL.Namespace}, clcm)

	if err == nil {
		l.V(2).Info("Execution type 1 Found - No Creation")
		return nil
	}

	l.V(2).Info("Creating Execution")
	//Creating the Execution Object with the right Spec
	clcm = &closedlooppoocv1.Execution{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: CL.Namespace,
			Name:      CL.Name + "-execution",
		},
		Spec: closedlooppoocv1.ExecutionSpec{

			Affix:           CL.Name,
			Time:            time.Now().String(),
			Action:          CL.Spec.Execution.Spec.Action,
			ExecutionTypeId: 1,
		},
	}

	//Ask Kubernetes to create it
	return r.Create(ctx, clcm)

}

func (r *ClosedLoopReconciler) deleteClosedLoopMonitoringType1(ctx context.Context, name string, namespace string, l VerbosityLog) error {
	// Try to retrieve the CR to see if it's already within the Cluster
	name = name + "-monitoring"
	clcm := &closedlooppoocv1.Monitoring{}
	err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, clcm)

	if err != nil {
		if errors.IsNotFound(err) {
			l.V(2).Info("Monitoring type 1 already deleted or no present")
			return nil
		}
		return err
	}
	//Delete Object if present within the CLuster
	err = r.Delete(ctx, clcm)
	if err != nil {
		l.Error(err, "Failed to delete monitoring")
		return err
	}

	l.V(2).Info("Deleted monitoring")

	return nil

}

func (r *ClosedLoopReconciler) deleteClosedLoopMonitoringType2(ctx context.Context, name string, namespace string, l VerbosityLog) error {
	// Try to retrieve the CR to see if it's already within the Cluster
	name = name + "-monitoring"
	clcm := &closedlooppoocv1.Monitoringv2{}
	err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, clcm)

	if err != nil {
		if errors.IsNotFound(err) {
			l.V(2).Info("Monitoring type 2 already deleted or no present")
			return nil
		}
		return err
	}
	//Delete Object if present within the CLuster
	err = r.Delete(ctx, clcm)
	if err != nil {
		l.Error(err, "Failed to delete monitoring")
		return err
	}

	l.Info("Deleted monitoring")

	return nil

}

func (r *ClosedLoopReconciler) deleteClosedLoopDecisionType1(ctx context.Context, name string, namespace string, l VerbosityLog) error {
	// Try to retrieve the CR to see if it's already within the Cluster
	name = name + "-decision"
	clcm := &closedlooppoocv1.Decision{}
	err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, clcm)

	if err != nil {
		if errors.IsNotFound(err) {
			l.V(2).Info("Decision type 1 already deleted or no present")
			return nil
		}
		return err
	}
	//Delete Object if present within the CLuster
	err = r.Delete(ctx, clcm)
	if err != nil {
		l.Error(err, "Failed to delete decision")
		return err
	}

	l.V(2).Info("Deleted decision")

	return nil

}

func (r *ClosedLoopReconciler) deleteClosedLoopExecutionType1(ctx context.Context, name string, namespace string, l VerbosityLog) error {
	// Try to retrieve the CR to see if it's already within the Cluster
	name = name + "-execution"
	clcm := &closedlooppoocv1.Execution{}
	err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, clcm)

	if err != nil {
		if errors.IsNotFound(err) {
			l.V(2).Info("Execution type 1 already deleted or no present")
			return nil
		}
		return err
	}
	//Delete Object if present within the CLuster
	err = r.Delete(ctx, clcm)
	if err != nil {
		l.Error(err, "Failed to delete execution")
		return err
	}
	l.V(2).Info("Deleted execution")

	return nil

}

// SetupWithManager sets up the controller with the Manager.
func (r *ClosedLoopReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&closedlooppoocv1.ClosedLoop{}).
		Complete(r)
}
