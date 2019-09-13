package jenkinspipeline

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	jakubbujnyv1alpha1 "github.com/jakubbujny/jenkins-pipeline-operator/pkg/apis/jakubbujny/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"os"
	"net/http"
	b64 "encoding/base64"

	coreErrors "errors"
	logr "github.com/go-logr/logr"
)

var log = logf.Log.WithName("controller_jenkinspipeline")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new JenkinsPipeline Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileJenkinsPipeline{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("jenkinspipeline-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource JenkinsPipeline
	err = c.Watch(&source.Kind{Type: &jakubbujnyv1alpha1.JenkinsPipeline{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}


	return nil
}

// blank assignment to verify that ReconcileJenkinsPipeline implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileJenkinsPipeline{}

// ReconcileJenkinsPipeline reconciles a JenkinsPipeline object
type ReconcileJenkinsPipeline struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a JenkinsPipeline object and makes changes based on the state read
// and what is in the JenkinsPipeline.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.

func (r *ReconcileJenkinsPipeline) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling JenkinsPipeline")

	// Fetch the JenkinsPipeline instance
	instance := &jakubbujnyv1alpha1.JenkinsPipeline{}
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

	resp, err := getSeedJob()

	if err != nil {
		reqLogger.Error(err, "Failed to get seed config to check whether job exists")
		return reconcile.Result{}, err
	}

	if resp.StatusCode == 404 {
		reqLogger.Info("Seed job not found so must be created for microservice "+instance.Spec.Microservice)

		resp, err := createSeedJob()
		err = handleResponse(resp, err, reqLogger, "create seed job")
		if err != nil {
			return reconcile.Result{}, err
		}

		resp, err = updateSeedJob(instance.Spec.Microservice)
		err = handleResponse(resp, err, reqLogger, "update seed job")
		if err != nil {
			return reconcile.Result{}, err
		}
	} else if resp.StatusCode == 200 {
		reqLogger.Info("Seed job found so must be updated for microservice "+instance.Spec.Microservice)
		resp, err = updateSeedJob(instance.Spec.Microservice)
		err = handleResponse(resp, err, reqLogger, "update seed job")
		if err != nil {
			return reconcile.Result{}, err
		}
	} else {
		err = coreErrors.New(fmt.Sprintf("Received invalid response from Jenkins %s",resp.Status))
		reqLogger.Error(err, "Failed to get seed config to check whether job exists")
		return reconcile.Result{}, err
	}

	resp, err = triggerSeedJob()
	err = handleResponse(resp, err, reqLogger, "trigger seed job")
	if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func handleResponse( resp *http.Response, err error, reqLogger logr.Logger, action string) error {
	if err != nil {
		reqLogger.Error(err, "Failed to "+action)
		return err
	}

	if resp.StatusCode != 200 {
		err = coreErrors.New(fmt.Sprintf("Received invalid response from Jenkins %s",resp.Status))
		reqLogger.Error(err, "Failed to"+action)
		return err
	}
	return nil
}

func decorateRequestToJenkinsWithAuth(req *http.Request) {
	jenkinsApiToken := os.Getenv("JENKINS_API_TOKEN")
	req.Header.Add("Authorization", "Basic "+ b64.StdEncoding.EncodeToString([]byte("admin:"+jenkinsApiToken)))
}

func getSeedJob() (*http.Response, error) {
	req, err := http.NewRequest("GET", os.Getenv("JENKINS_URL")+"/job/seed/config.xml", nil)
	if err != nil {
		return nil, err
	}
	decorateRequestToJenkinsWithAuth(req)
	return (&http.Client{}).Do(req)
}

func createSeedJob() (*http.Response, error) {
	seedFileData, err := ioutil.ReadFile("/opt/seed.xml")

	req, err := http.NewRequest("POST", os.Getenv("JENKINS_URL")+"/createItem?name=seed", bytes.NewBuffer(seedFileData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-type", "text/xml")
	decorateRequestToJenkinsWithAuth(req)
	return (&http.Client{}).Do(req)
}

func updateSeedJob(microservice string) (*http.Response, error) {
	resp, err := getSeedJob()
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}
	seedXml := buf.String()

	r := regexp.MustCompile(`<defaultValue>(.+)<\/defaultValue>`)
	foundMicroservices := r.FindStringSubmatch(seedXml)

	toReplace := ""
	if strings.Contains(foundMicroservices[1], microservice) {
		return nil,nil
	} else {
		if foundMicroservices[1] == "default" {
			toReplace = microservice
		} else {
			toReplace = foundMicroservices[1] + "," + microservice
		}
	}

	toUpdate := r.ReplaceAllString(seedXml, fmt.Sprintf("<defaultValue>%s</defaultValue>", toReplace))

	req, err := http.NewRequest("POST", os.Getenv("JENKINS_URL")+"/job/seed/config.xml", bytes.NewBuffer([]byte(toUpdate)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-type", "text/xml")
	decorateRequestToJenkinsWithAuth(req)
	return (&http.Client{}).Do(req)
}

func triggerSeedJob() (*http.Response, error) {
	req, err := http.NewRequest("POST", os.Getenv("JENKINS_URL")+"/job/seed/buildWithParameters", nil)
	if err != nil {
		return nil, err
	}
	decorateRequestToJenkinsWithAuth(req)
	return (&http.Client{}).Do(req)
}