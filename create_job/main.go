package main

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetes "k8s.io/client-go/kubernetes"
	clientcmd "k8s.io/client-go/tools/clientcmd"
)

func connectToK8s() *kubernetes.Clientset {
	home, exists := os.LookupEnv("HOME")
	if !exists {
		// get the user from CL options, Build the User's home dire using sprintf instead of hardcoding.
		home = "/root/"
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

func launchK8sJob(clientset *kubernetes.Clientset, jobName *string, image *string, cmd *string, ns *string) {
	jobs := clientset.BatchV1().Jobs(*ns)
	var backOffLimit int32 = 0

	jobSpec := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      *jobName,
			Namespace: *ns,
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:    *jobName,
							Image:   *image,
							Command: strings.Split(*cmd, " "),
						},
					},
					RestartPolicy: v1.RestartPolicyNever,
				},
			},
			BackoffLimit: &backOffLimit,
		},
	}

	_, err := jobs.Create(context.Background(), jobSpec, metav1.CreateOptions{})
	if err != nil {
		log.Fatalln("Failed to create K8s job")
	}

	//print job details
	log.Println("Created K8s job successfully")
}

func main() {
	jobName := flag.String("jobname", "test_job", "The name of the job")
	containerImage := flag.String("image", "centos:latest", "Name of the container image")
	entryCommand := flag.String("command", "ls", "The command to run inside the container")
	ns := flag.String("ns", "default", "namespace name")

	flag.Parse()

	clientset := connectToK8s()
	launchK8sJob(clientset, jobName, containerImage, entryCommand, ns)
}
