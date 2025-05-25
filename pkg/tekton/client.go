package tekton

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	tektonv1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	tektonClient "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// Client represents a Tekton client
type Client struct {
	tektonClientset *tektonClient.Clientset
	namespace       string
}

// NewClient creates a new Tekton client
func NewClient(namespace string) (*Client, error) {
	// Get kubeconfig path
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = filepath.Join(home, ".kube", "config")
		} else {
			return nil, fmt.Errorf("kubeconfig not found and $HOME is not defined")
		}
	}

	// Build config from kubeconfig file
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to build config: %v", err)
	}

	// Create Tekton clientset
	tektonClientset, err := tektonClient.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create tekton clientset: %v", err)
	}

	return &Client{
		tektonClientset: tektonClientset,
		namespace:       namespace,
	}, nil
}

// ListTasks returns a list of all Tasks in the namespace
func (c *Client) ListTasks() (*tektonv1.TaskList, error) {
	return c.tektonClientset.TektonV1().Tasks(c.namespace).List(context.Background(), metav1.ListOptions{})
}

// ListTaskRuns returns a list of all TaskRuns in the namespace
func (c *Client) ListTaskRuns() (*tektonv1.TaskRunList, error) {
	return c.tektonClientset.TektonV1().TaskRuns(c.namespace).List(context.Background(), metav1.ListOptions{})
}

// CreateTaskRun creates a new TaskRun for the specified Task
func (c *Client) CreateTaskRun(taskName string, params map[string]string) (*tektonv1.TaskRun, error) {
	// Create TaskRun spec
	taskRun := &tektonv1.TaskRun{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("%s-run-", taskName),
			Namespace:    c.namespace,
		},
		Spec: tektonv1.TaskRunSpec{
			TaskRef: &tektonv1.TaskRef{
				Name: taskName,
			},
		},
	}

	// Add parameters if provided
	if len(params) > 0 {
		taskRunParams := []tektonv1.Param{}
		for name, value := range params {
			taskRunParams = append(taskRunParams, tektonv1.Param{
				Name: name,
				Value: tektonv1.ParamValue{
					Type:      tektonv1.ParamTypeString,
					StringVal: value,
				},
			})
		}
		taskRun.Spec.Params = taskRunParams
	}

	// Create the TaskRun
	return c.tektonClientset.TektonV1().TaskRuns(c.namespace).Create(context.Background(), taskRun, metav1.CreateOptions{})
}

// CreateTaskRunFromTaskRef creates a new TaskRun from a task reference
func (c *Client) CreateTaskRunFromTaskRef(taskName string) error {
	_, err := c.CreateTaskRun(taskName, nil)
	return err
}
