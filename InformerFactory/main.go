package main

import (
	"flag"
	"fmt"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"time"

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
		fmt.Printf("error is %s building config from flags\n", err.Error())
		// Inclusterconfig uses default SA which mounts on each to get the cluster config file.
		config, err = rest.InClusterConfig()
		// this will give authorization issues luke cannot list resource pods because pod uses default SA, that SA does not have  required permisions.
		if err != nil {
			fmt.Printf("error is %s building config from flags\n", err.Error())
		}

	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	informerfactory := informers.NewSharedInformerFactory(clientset, 10*time.Minute)

	/*
		informerfactory := informers.NewFilteredSharedInformerFactory(clientset, 10*time.Minute, ns, func(lo *metav1.ListOptions) {
			lo.LabelSelector = "app"
		})

	*/

	podinformation := informerfactory.Core().V1().Pods()
	podinformation.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(new interface{}) {
			fmt.Println("Pod added and  trigger something")
		},
		UpdateFunc: func(old, new interface{}) {
			fmt.Println("Updated trigger logic")

		},
		DeleteFunc: func(obj interface{}) {
			fmt.Println("Updated trigger logic")
		},
	})

	informerfactory.Start(wait.NeverStop)
	informerfactory.WaitForCacheSync(wait.NeverStop)
	pod, err := podinformation.Lister().Pods(ns).Get(ns)
	fmt.Println(pod)
}
