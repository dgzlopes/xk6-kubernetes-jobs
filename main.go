package kubernetes

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"

	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetes "k8s.io/client-go/kubernetes"
	clientcmd "k8s.io/client-go/tools/clientcmd"

	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
)

// Register the extension on module initialization, available to
// import from JS as "k6/x/kubernetes-jobs".
func init() {
	modules.Register("k6/x/kubernetes-jobs", new(Job))
}

// Job is the k6 extension for interacting with Kubernetes jobs.
type Job struct{}

// Client is the Kubernetes client wrapper.
type Client struct {
	Client    *kubernetes.Clientset
	Namespace string
}

// XClient represents the Client constructor (i.e. `new kubernetes.Client()`) and
// returns a new Kubernetes client object.
func (r *Job) XClient(ctxPtr *context.Context) interface{} {
	rt := common.GetRuntime(*ctxPtr)
	return common.Bind(rt, &Client{Client: clientFromConfig(), Namespace: "observability"}, ctxPtr)
}

func clientFromConfig() *kubernetes.Clientset {
	home, exists := os.LookupEnv("HOME")
	if !exists {
		home = "/root"
	}

	configPath := filepath.Join(home, ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		log.Fatalln("failed to create K8s config")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln("Failed to create K8s clientset")
	}

	return clientset
}

// Create a new Job
func (c *Client) Create(name string, image string, cmd string) string {
	jobs := c.Client.BatchV1().Jobs(c.Namespace)
	var ttlAfterFinished int32 = 0
	var backOffLimit int32 = 0

	jobSpec := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: c.Namespace,
			Labels:    map[string]string{"job-type": "k6"},
		},
		Spec: batchv1.JobSpec{
			TTLSecondsAfterFinished: &ttlAfterFinished,
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:    name,
							Image:   image,
							Command: strings.Split(cmd, " "),
						},
					},
					RestartPolicy: v1.RestartPolicyNever,
				},
			},
			BackoffLimit: &backOffLimit,
		},
	}

	created, err := jobs.Create(context.TODO(), jobSpec, metav1.CreateOptions{})
	if err != nil {
		log.Fatalln(err)
	}
	return created.GetName()
}

// Delete an existing job
func (c *Client) Delete(name string) {
	jobs := c.Client.BatchV1().Jobs(c.Namespace)
	err := jobs.Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		log.Fatalln(err)
	}
}

// DeleteAll the existing jobs
func (c *Client) DeleteAll() {
	jobs := c.Client.BatchV1().Jobs(c.Namespace)
	allTheJobs, _ := jobs.List(context.TODO(), metav1.ListOptions{
		LabelSelector: "job-type=k6",
	})
	for _, s := range allTheJobs.Items {
		err := jobs.Delete(context.TODO(), s.GetName(), metav1.DeleteOptions{})
		if err != nil {
			log.Fatalln(err)
		}
	}
}

// List the names of the existing jobs
func (c *Client) List() []string {
	jobs := c.Client.BatchV1().Jobs(c.Namespace)
	jobList := []string{}
	allTheJobs, _ := jobs.List(context.TODO(), metav1.ListOptions{
		LabelSelector: "job-type=k6",
	})
	for _, s := range allTheJobs.Items {
		jobList = append(jobList, s.GetName())
	}
	return jobList
}

// Get a certain job
func (c *Client) Get(name string) *batchv1.Job {
	jobs := c.Client.BatchV1().Jobs(c.Namespace)
	job, _ := jobs.Get(context.TODO(), name, metav1.GetOptions{})
	return job
}
