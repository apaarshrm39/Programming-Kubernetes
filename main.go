package main

import (
	"flag"

	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := flag.String("kubeconfig", "~/.kube/config", "kubeconfig")
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
}
