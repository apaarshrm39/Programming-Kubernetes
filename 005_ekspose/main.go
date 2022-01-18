package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		fmt.Sprintf("Could not build config %s\n", err)
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Fatal("could not build config", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal("Could not build clientset", err)
	}

	informerfactory := informers.NewSharedInformerFactory(clientset, 10*time.Minute)

	//podinformer := informerfactory.Core().V1().Pods()
	//informerfactory.Apps().V1().Deployments()

	c := newController(*clientset, informerfactory.Apps().V1().Deployments())
	ch := make(chan struct{})
	informerfactory.Start(ch)
	c.run(ch)
}
