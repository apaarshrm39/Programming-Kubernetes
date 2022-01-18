package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "kubeconfig")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "kubeconfig")
	}
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		fmt.Printf("could not create config from kubeconfig %s\n", err)
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Fatal(err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	// Creating a shared informer factory

	informerfactory := informers.NewSharedInformerFactory(clientset, 10*time.Minute) //Sync the in memory cache with kubernetes cluster state
	// You can also create fltered information factory if you want
	informers.NewFilteredSharedInformerFactory(clientset, 10*time.Minute, "default", func(lo *v1.ListOptions) {
		lo.Kind = "Pods"
	})
	// create Informer for Pods
	podinformer := informerfactory.Core().V1().Pods()

	// Add Event handlers to Pod informer
	podinformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			fmt.Println("Add was called")
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			fmt.Println("Update was called")
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Println("Delete was called")
		},
	})

	// Start the informerfactory, this initializes the In memory store
	informerfactory.Start(wait.NeverStop)
	// Kickstart the cache with the list request
	informerfactory.WaitForCacheSync(wait.NeverStop)
	pod, err := podinformer.Lister().Pods("default").Get("k9s")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(pod)
}
