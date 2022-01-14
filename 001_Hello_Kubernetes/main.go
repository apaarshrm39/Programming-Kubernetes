package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	//_, err = clientset.CoreV1().Pods("booksapp").Get(context.TODO(), "authors-b7bbfb747-cnxkz", metav1.GetOptions{})
	pod, err := clientset.AppsV1().Deployments("booksapp").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	for _, dep := range pod.Items {
		fmt.Println(dep.Name)
	}

	st := []string{"Hi", "Hello"}

	stuff(st...)
}

// accepts O or more parameter and reference them as slice
func stuff(st ...string) {
	fmt.Println(st)
}
