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
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"io"
	"net/http"

	closedlooppoocv1 "closedloop/api/v1"
)

// MonitoringDReconciler reconciles a MonitoringD object
type MonitoringDReconciler struct {
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
// the Monitoring object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *MonitoringDReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	verbosityLog := new (VerbosityLog)
	verbosityLog.SetMaxLevel(1)
	l := verbosityLog.FromContext(ctx)
	l.V(1).Info(">>>>>>>>>>>>>>>>>>>>>>>>>")
	l.V(1).Info("Enter Reconcile MonitoringD", "req", req)
	//Retreiving ClosedLoop Object who triggered the Reconciler
	MonitoringD := &closedlooppoocv1.MonitoringD{}
	err := r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, MonitoringD)
	if err != nil {
		l.V(2).Info("MonitoringD Not found")
		return ctrl.Result{}, nil
	}
	// Update Status if it's not already done
	// Here we don't use Status for our logic, so the section is not needed but only to show how to use it
	if MonitoringD.Name != MonitoringD.Status.Affix {
		MonitoringD.Status.Affix = MonitoringD.Name
		if err := r.Status().Update(ctx, MonitoringD); err != nil {
			l.Error(err, "Failed to update MonitoringD status")
			return ctrl.Result{}, err
		}
		l.V(2).Info("Enter Reconcile View Spec & Status", "spec", MonitoringD.Spec, "status", MonitoringD.Status)
	}

	/* ---------------------------------- START Monitoring Part ---------------------------------- */

	// Retreive Data

	url := "http://" + MonitoringD.Spec.Source.Addresse
	response, err := http.Get(url)
	if err != nil {
		l.V(2).Info("Error on GET : %s\n", err)
		return ctrl.Result{}, err
	}
	defer response.Body.Close()

	// Read Data
	body, err := io.ReadAll(response.Body)
	if err != nil {
		l.V(2).Info("Error on ReadAll(response.Body) : %s\n", err)
		return ctrl.Result{}, err
	}

	// Convert to string Debug
	content := string(body)
	l.V(2).Info(content)

	var data []map[string]int

	err = json.Unmarshal([]byte(content), &data)
	if err != nil {
		l.V(2).Info("Error on JSON convertion :", err)
		return ctrl.Result{}, err
	}

	message := "No Message"
	ObjectData := -1
	message = ""

	if len(data) > 0 {
		for key, value := range MonitoringD.Spec.MonitoringPolicies.Data {
			l.V(2).Info("Data Key:", key, " Value: ", value)

			lastObj := data[len(data)-1]
			ObjectData = lastObj[value]
			l.V(2).Info("Values Of the last Object", value, " = ", ObjectData)
			if ObjectData == -1 {

				l.V(2).Info("No Data Found for the monitoring Spec Given : ", value)
				return ctrl.Result{}, nil
			}

			thresholdvalue, err := strconv.Atoi(MonitoringD.Spec.MonitoringPolicies.TresholdValue[key+"-thresholdvalue"])
			l.V(2).Info(MonitoringD.Spec.MonitoringPolicies.TresholdValue[key+"-thresholdvalue"])

			if err != nil {
				l.V(2).Info("Error during conversion")
				l.Error(err, "Failed Cast")
				return ctrl.Result{}, err
			}

			switch MonitoringD.Spec.MonitoringPolicies.TresholdKind[key+"-thresholdkind"] {
			case "inferior":
				l.V(2).Info("Inferior")
				if ObjectData < thresholdvalue {
					message = message + ",Low " + value
				}
			case "superior":
				l.V(2).Info("Superior")
				if ObjectData > thresholdvalue {
					message = message + ",High " + value

				}
			case "equal":
				l.V(2).Info("equal")
				if ObjectData == thresholdvalue {
					message = message + ",Equal " + value

				}

			}
		}
	} else {
		l.V(2).Info("Object List is Empty.")
	}

	if message != "" && string(message[0]) == "," {
		message = strings.TrimPrefix(message, string(message[0]))
	} else {
		message = "No Event"
	}

	/* ---------------------------------- END  Monitoring Part ---------------------------------- */

	l.V(1).Info("message Send to Decision " + message)

	/* -------------------------------- Apply modification on Decision --------------------- */

	switch MonitoringD.Spec.DecisionKind {
	case "Decision":
		if err := r.ApplyDecision(ctx, MonitoringD, l, message); err != nil {
			l.Error(err, "Failed to ApplyDecision")
			return ctrl.Result{}, err
		}

	}

	return ctrl.Result{
		Requeue:      true,
		RequeueAfter: time.Duration(MonitoringD.Spec.Source.Interval) * time.Second,
	}, nil
}

/*func fetchURL(url string) (string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	response, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}*/

func (r *MonitoringDReconciler) ApplyDecision(ctx context.Context, monitoringd *closedlooppoocv1.MonitoringD, l VerbosityLog, Message string) error {
	// Try to retrieve the CR that we want to update
	DecisionD := &closedlooppoocv1.DecisionD{}
	r.Get(ctx, types.NamespacedName{Name: monitoringd.Spec.Affix + "-decision", Namespace: monitoringd.Namespace}, DecisionD)

	l.V(1).Info("Update Message on Decision")
	//Update it's field with the variable message
	//DecisionD.Spec.Message = Message
	return r.Update(ctx, DecisionD)

}

// SetupWithManager sets up the controller with the Manager.
func (r *MonitoringDReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&closedlooppoocv1.MonitoringD{}).
		Complete(r)
}
