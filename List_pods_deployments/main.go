package main

import (
	"context"
	"flag"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := flag.String("Kubeconfig", "/Users/bkasa724/.kube/config", "location to kube config file")
	namespace := flag.String("namespace", "default", "namespace name")
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	ns := *namespace
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	pods, err := clientset.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{})

	if err != nil {
		panic(err.Error())
	}
	fmt.Println(pods.Items[0])
	for _, pod := range pods.Items {
		fmt.Println(pod.Name)
	}

	deployments, err := clientset.AppsV1().Deployments(ns).List(ctx, metav1.ListOptions{})
	fmt.Println("")
	fmt.Println("")
	fmt.Println("deployments are")
	fmt.Println("")

	for _, deployment := range deployments.Items {
		fmt.Println(deployment.Name)
	}

}
