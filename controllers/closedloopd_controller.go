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
	"time"

	//"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	//"sigs.k8s.io/controller-runtime/pkg/log"

	closedlooppoocv1 "closedloop/api/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClosedLoopReconciler reconciles a ClosedLoop object
type ClosedLoopDReconciler struct {
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
//+kubebuilder:rbac:groups="apps",resources=deploymentds,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="apps",resources=deploymentds/status,verbs=get;watch;list
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
func (r *ClosedLoopDReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	verbosityLog := new (VerbosityLog)
	verbosityLog.SetMaxLevel(1)
	l := verbosityLog.FromContext(ctx)
	
	l.V(2).Info("Enter ******************* Reconcile", "req", req)

	//Retreiving ClosedLoop Object who triggered the Reconciler
	closedLoop := &closedlooppoocv1.ClosedLoopD{}
	err := r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, closedLoop)

	// Update Status if it's not already done
	// Here we don't use Status for our logic, so the section is not needed but only to show how to use it
	if closedLoop.Name != closedLoop.Status.Name {
		closedLoop.Status.Name = closedLoop.Name

		if err := r.Status().Update(ctx, closedLoop); err != nil {
			l.Error(err, "Failed to update closedLoop D status")
			return ctrl.Result{}, err
		}
		l.Info("Enter Reconcile View Spec & Status", "spec", closedLoop.Spec, "status", closedLoop.Status)
	}

	// If the closedLoop who triggered is not find it means that it's been deleted
	if err != nil {

		// The Closedloop object has been deleted or is not found, so we should delete the associated CR (if they exist)
		l.Info("Close loop D Instance not found, Deletion Close Loop D ressources")

		if err := r.deleteClosedLoopDMonitoringType1(ctx, req.Name, req.Namespace, l); err != nil {
			l.Error(err, "No MonitoringD found, no deletions have been made")

		}
		if err := r.deleteClosedLoopDMonitoringType2(ctx, req.Name, req.Namespace, l); err != nil {
			l.Error(err, "No MonitoringDv2 found, no deletions have been made")

		}

		if err := r.deleteClosedLoopDDecisionType1(ctx, req.Name, req.Namespace, l); err != nil {
			l.Error(err, "No DecisionD type 1 found, no deletions have been made")

		}

		if err := r.deleteClosedLoopDExecutionType1(ctx, req.Name, req.Namespace, l); err != nil {
			l.Error(err, "No ExecutionD type 1 found, no deletions have been made")

		}

		//l.Info()
		return ctrl.Result{}, nil

	}

	// Creation of CR implementation layer ressources (Second layer) are created based on Kinds defined in the primary closedloop ressource
	// These functions are called each time, so if they've already been created, nothing will happen.
	switch closedLoop.Spec.Execution.Kind {
	case "ExecutionD":
		if err := r.createExecutionType1(ctx, closedLoop, l); err != nil {
			l.Error(err, "Failed to createExecutiondType1")
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
	case "DecisionD":
		if err := r.createDecisionType1(ctx, closedLoop, l); err != nil {
			l.Error(err, "Failed to createDecisiondType1")
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
	case "MonitoringD":
		if err := r.createMonitoringType1(ctx, closedLoop, l); err != nil {
			l.Error(err, "Failed to createMonitoringdType1")
			return ctrl.Result{}, err
		}
	case "MonitoringDv2":
		if err := r.createMonitoringType2(ctx, closedLoop, l); err != nil {
			l.Error(err, "Failed to createMonitoringdType2")
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
func (r *ClosedLoopDReconciler) createMonitoringType1(ctx context.Context, CL *closedlooppoocv1.ClosedLoopD, l VerbosityLog) error {

	// Try to retrieve the CR to see if it's already within the Cluster
	clcm := &closedlooppoocv1.Monitoring{}
	err := r.Get(ctx, types.NamespacedName{Name: CL.Name + "-monitoringd", Namespace: CL.Namespace}, clcm)

	if err == nil {
		l.Info("monitoringd type 1 Found - No Creation")
		return nil
	}

	l.V(2).Info("Creating monitoringd")

	//Creating the Monitoring Object with the right Spec
	clcm = &closedlooppoocv1.Monitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: CL.Namespace,
			Name:      CL.Name + "-monitoringd",
		},
		Spec: closedlooppoocv1.MonitoringSpec{
			Affix:        CL.Name,
			DecisionKind: CL.Spec.Decision.DecisionKind.DecisionKindName,
			Source: closedlooppoocv1.Source{
				Addresse: CL.Spec.Monitoring.MonitoringKind.Source.Addresse,
				Port:     CL.Spec.Monitoring.MonitoringKind.Source.Port,
				Interval: CL.Spec.Monitoring.MonitoringKind.Source.Interval,
			},
		},
	}

	return r.Create(ctx, clcm)

}

// Function use to create the CR Monitoring type 2 (kind : Monitoringv2)
func (r *ClosedLoopDReconciler) createMonitoringType2(ctx context.Context, CL *closedlooppoocv1.ClosedLoopD, l VerbosityLog) error {

	// Try to retrieve the CR to see if it's already within the Cluster
	clcm := &closedlooppoocv1.MonitoringDv2{}
	err := r.Get(ctx, types.NamespacedName{Name: CL.Name + "-monitoringd", Namespace: CL.Namespace}, clcm)

	m := make(map[string]string)
	m["Value"] = "No Value"
	if err == nil {
		l.V(2).Info("monitoringd type 2 Found - No Creation")
		return nil
	}

	l.V(2).Info("Creating monitoringd v2")
	l.V(2).Info("Value Decision kind : ", CL.Spec.Decision.DecisionKind.DecisionKindName)
	//Creating the Monitoringv2 Object with the right Spec
	clcm = &closedlooppoocv1.MonitoringDv2{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: CL.Namespace,
			Name:      CL.Name + "-monitoringd",
		},
		Spec: closedlooppoocv1.MonitoringDv2Spec{
			Affix:        CL.Name,
			Time:         time.Now().String(),
			DecisionKind: CL.Spec.Decision.DecisionKind.DecisionKindName,
			RequestedPod: CL.Spec.Monitoring.MonitoringKind.RequestedPod,
			Data: m,
		},
	}

	//Ask Kubernetes to create it
	return r.Create(ctx, clcm)

}

// Function use to create the CR Decision type 1 (kind : Decision)
func (r *ClosedLoopDReconciler) createDecisionType1(ctx context.Context, CL *closedlooppoocv1.ClosedLoopD, l VerbosityLog) error {

	// Try to retrieve the CR to see if it's already within the Cluster
	clcm := &closedlooppoocv1.DecisionD{}
	err := r.Get(ctx, types.NamespacedName{Name: CL.Name + "-decisiond", Namespace: CL.Namespace}, clcm)

	if err == nil {
		l.V(2).Info("Decisiond type 1Found - No Creation")
		return nil
	}

	l.V(2).Info("Creating decisiond")
	//Creating the Decision Object with the right Spec
	clcm = &closedlooppoocv1.DecisionD{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: CL.Namespace,
			Name:      CL.Name + "-decisiond",
		},
		Spec: closedlooppoocv1.DecisionDSpec{

			Affix:         CL.Name,
			Time:          time.Now().String(),
		},
	}

	//Ask Kubernetes to create it
	return r.Create(ctx, clcm)

}

// Function use to create the CR Execution type 1 (kind : Execution)
func (r *ClosedLoopDReconciler) createExecutionType1(ctx context.Context, CL *closedlooppoocv1.ClosedLoopD, l VerbosityLog) error {

	// Try to retrieve the CR to see if it's already within the Cluster
	clcm := &closedlooppoocv1.ExecutionD{}
	err := r.Get(ctx, types.NamespacedName{Name: CL.Name + "-executiond", Namespace: CL.Namespace}, clcm)

	if err == nil {
		l.V(2).Info("ExecutionD type 1 Found - No Creation")
		return nil
	}

	l.V(2).Info("Creating ExecutionD")
	//Creating the Execution Object with the right Spec
	clcm = &closedlooppoocv1.ExecutionD{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: CL.Namespace,
			Name:      CL.Name + "-executiond",
		},
		Spec: closedlooppoocv1.ExecutionDSpec{

			Affix:           CL.Name,
			Time:            time.Now().String(),
			Action:          CL.Spec.Execution.Spec.Action,
			Metric:			 "No metric",
			ExecutionTypeId: 1,
		},
	}

	//Ask Kubernetes to create it
	return r.Create(ctx, clcm)

}

func (r *ClosedLoopDReconciler) deleteClosedLoopDMonitoringType1(ctx context.Context, name string, namespace string, l VerbosityLog) error {
	// Try to retrieve the CR to see if it's already within the Cluster
	name = name + "-monitoringd"
	clcm := &closedlooppoocv1.MonitoringD{}
	err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, clcm)

	if err != nil {
		if errors.IsNotFound(err) {
			l.Info("MonitoringD type 1 already deleted or no present")
			return nil
		}
		return err
	}
	//Delete Object if present within the CLuster
	err = r.Delete(ctx, clcm)
	if err != nil {
		l.Error(err, "Failed to delete monitoringd")
		return err
	}

	l.Info("Deleted monitoringd")

	return nil

}

func (r *ClosedLoopDReconciler) deleteClosedLoopDMonitoringType2(ctx context.Context, name string, namespace string, l VerbosityLog) error {
	// Try to retrieve the CR to see if it's already within the Cluster
	name = name + "-monitoringd"
	clcm := &closedlooppoocv1.MonitoringDv2{}
	err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, clcm)

	if err != nil {
		if errors.IsNotFound(err) {
			l.Info("MonitoringD type 2 already deleted or no present")
			return nil
		}
		return err
	}
	//Delete Object if present within the CLuster
	err = r.Delete(ctx, clcm)
	if err != nil {
		l.Error(err, "Failed to delete monitoringd")
		return err
	}

	l.Info("Deleted monitoringd")

	return nil

}

func (r *ClosedLoopDReconciler) deleteClosedLoopDDecisionType1(ctx context.Context, name string, namespace string, l VerbosityLog) error {
	// Try to retrieve the CR to see if it's already within the Cluster
	name = name + "-decisiond"
	clcm := &closedlooppoocv1.DecisionD{}
	err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, clcm)

	if err != nil {
		if errors.IsNotFound(err) {
			l.Info("Decision type 1 already deleted or no present")
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

	l.Info("Deleted decision")

	return nil

}

func (r *ClosedLoopDReconciler) deleteClosedLoopDExecutionType1(ctx context.Context, name string, namespace string, l VerbosityLog) error {
	// Try to retrieve the CR to see if it's already within the Cluster
	name = name + "-executiond"
	clcm := &closedlooppoocv1.ExecutionD{}
	err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, clcm)

	if err != nil {
		if errors.IsNotFound(err) {
			l.Info("ExecutionD type 1 already deleted or no present")
			return nil
		}
		return err
	}
	//Delete Object if present within the CLuster
	err = r.Delete(ctx, clcm)
	if err != nil {
		l.Error(err, "Failed to delete executiond")
		return err
	}
	l.Info("Deleted executiond")

	return nil

}

// SetupWithManager sets up the controller with the Manager.
func (r *ClosedLoopDReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&closedlooppoocv1.ClosedLoopD{}).
		Complete(r)
}
