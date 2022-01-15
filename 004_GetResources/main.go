package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/ghodss/yaml"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog/v2"
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

	depoymentclient := clientset.AppsV1().Deployments(v1.NamespaceDefault)

	var deploy appsv1.Deployment
	yamlfile, err := ioutil.ReadFile("deployment.yaml")
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(yamlfile, &deploy)
	//fmt.Print(deploy)

	fmt.Println("Deploying", deploy.Name)

	result, err := depoymentclient.Create(context.TODO(), &deploy, v1.CreateOptions{})
	if err != nil {
		klog.Fatal(err)
	}

	fmt.Println("Deployed", result.Name)

}
