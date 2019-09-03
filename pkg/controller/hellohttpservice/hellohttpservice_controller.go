package hellohttpservice

import (
	"context"

	hellohttpv1alpha1 "github.com/raelga/k8s-hello-operator/pkg/apis/hellohttp/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_hellohttpservice")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new HelloHttpService Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileHelloHttpService{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("hellohttpservice-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource HelloHttpService
	err = c.Watch(&source.Kind{Type: &hellohttpv1alpha1.HelloHttpService{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner HelloHttpService
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &hellohttpv1alpha1.HelloHttpService{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileHelloHttpService implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileHelloHttpService{}

// ReconcileHelloHttpService reconciles a HelloHttpService object
type ReconcileHelloHttpService struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a HelloHttpService object and makes changes based on the state read
// and what is in the HelloHttpService.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileHelloHttpService) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling HelloHttpService")

	// Fetch the HelloHttpService instance
	instance := &hellohttpv1alpha1.HelloHttpService{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// // Define a new ConfigMap object
	// cm := newConfigMapForCR(instance)

	// // Check if this ConfigMap already exists
	// configMapFound := &corev1.ConfigMap{}
	// err = r.client.Get(context.TODO(), types.NamespacedName{Name: cm.Name, Namespace: cm.Namespace}, configMapFound)
	// if err != nil && errors.IsNotFound(err) {
	// 	reqLogger.Info("Creating a new ConfigMap", "ConfigMap.Namespace", cm.Namespace, "ConfigMap.Name", cm.Name)
	// 	err = r.client.Create(context.TODO(), cm)
	// 	if err != nil {
	// 		return reconcile.Result{}, err
	// 	}

	// 	// ConfigMap created successfully - don't requeue
	// 	return reconcile.Result{}, nil
	// } else if err != nil {
	// 	return reconcile.Result{}, err
	// }

	// // ConfigMap already exists - don't requeue
	// reqLogger.Info("Skip reconcile: ConfigMap already exists", "ConfigMap.Namespace", configMapFound.Namespace, "ConfigMap.Name", configMapFound.Name)

	rs := newReplicaSetCR(instance)

	// Set HelloHttpService instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, rs, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Pod already exists
	rsFound := &appsv1.ReplicaSet{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: rs.Name, Namespace: rs.Namespace}, rsFound)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new ReplicaSet", "ReplicaSet.Namespace", rs.Namespace, "ReplicaSet.Name", rs.Name)
		err = r.client.Create(context.TODO(), rs)
		if err != nil {
			return reconcile.Result{}, err
		}

		// ReplicaSet created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// ReplicaSet already exists - don't requeue
	reqLogger.Info("Skip reconcile: ReplicaSet already exists", "ReplicaSet.Namespace", rsFound.Namespace, "ReplicaSet.Name", rsFound.Name)
	return reconcile.Result{}, nil
}

// newConfigMapForCR returns a hello-http config map with the same name/namespace as the cr
// func newConfigMapForCR(cr *hellohttpv1alpha1.HelloHttpService) *corev1.ConfigMap {
// 	labels := map[string]string{
// 		"app": cr.Name,
// 	}
// 	return &corev1.ConfigMap{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      cr.Name + "-cmd",
// 			Namespace: cr.Namespace,
// 			Labels:    labels,
// 		},
// 		Data: map[string]string{"HELLO_NAME": cr.Spec.Subject},
// 	}
// }

// newReplicaSetCR returns a hello-http rs with the same name/namespace as the cr
func newReplicaSetCR(cr *hellohttpv1alpha1.HelloHttpService) *appsv1.ReplicaSet {
	labels := map[string]string{
		"app": cr.Name,
	}
	return &appsv1.ReplicaSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-rs",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.ReplicaSetSpec{
			Replicas: cr.Spec.Size,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      cr.Name + "-pod",
					Namespace: cr.Namespace,
					Labels:    labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "app",
							Image: "raelga/hello-http:latest",
							Env: []corev1.EnvVar{
								{Name: "HELLO_NAME", Value: cr.Spec.Subject},
							},
						},
					},
				},
			},
		},
	}
}
