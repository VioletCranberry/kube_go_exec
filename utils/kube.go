package utils

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"regexp"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

type KubeClient struct {
	clientset  *kubernetes.Clientset
	restconfig *rest.Config
}

func KubeClientFromConfig(bearer_token, master_url string) (*KubeClient, error) {
	var client = new(KubeClient)
	var err error

	if bearer_token == "" && master_url == "" {
		path := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		if _, err := os.Stat(path); err != nil {
			client.restconfig, err = rest.InClusterConfig()
			if err != nil {
				return nil, err
			}
		} else {
			client.restconfig, err = clientcmd.BuildConfigFromFlags("", path)
			if err != nil {
				return nil, err
			}
		}
	}

	if bearer_token != "" && master_url != "" {
		client.restconfig, err = clientcmd.BuildConfigFromFlags(master_url, "")
		if err != nil {
			return nil, err
		}
		client.restconfig.BearerToken = bearer_token
		client.restconfig.Insecure = true
	}

	client.clientset, err = kubernetes.NewForConfig(client.restconfig)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (kube *KubeClient) GetPodByFilter(filter, namespace string) (*v1.Pod, error) {
	podClient := kube.clientset.CoreV1().Pods(namespace)
	pods, err := podClient.List(context.TODO(), metav1.ListOptions{
		LabelSelector: filter,
	})
	if err != nil {
		return nil, err
	}

	var pod v1.Pod
	for _, p := range pods.Items {
		if p.Status.Phase == "Running" {
			pod = p
			break
		}
	}
	return &pod, nil
}

func (kube *KubeClient) ExecInPod(pod *v1.Pod, container string, commands string) (string, string, error) {
	restClient := kube.clientset.CoreV1().RESTClient()

	request := restClient.Post().
		Namespace(pod.Namespace).
		Resource("pods").
		Name(pod.Name).
		SubResource("exec").
		Param("container", container).
		Param("stdout", "true").
		Param("stderr", "true")

	r := regexp.MustCompile(`[^\s"']+|"([^"]*)"|'([^']*)`)
	for _, command := range r.FindAllString(commands, -1) {

		if len(command) >= 2 {
			if command[0] == '"' && command[len(command)-1] == '"' {
				request.Param("command", command[1:len(command)-1])
			} else {
				request.Param("command", command)
			}
		}
	}

	executor, err := remotecommand.NewSPDYExecutor(kube.restconfig, "POST", request.URL())
	if err != nil {
		return "", "", err
	}
	var stdout, stderr bytes.Buffer
	err = executor.Stream(remotecommand.StreamOptions{
		Stdin:  nil,
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	})
	return stdout.String(), stderr.String(), err
}
