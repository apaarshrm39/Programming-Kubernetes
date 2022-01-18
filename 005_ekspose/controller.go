package main

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	depInfromer "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type controller struct {
	clientset      *kubernetes.Clientset           // to interact with the k8s resources
	deplLister     v1.DeploymentLister             // deployment lister
	depCacheSynced cache.InformerSynced            // tocheck if cache has synced
	queue          workqueue.RateLimitingInterface //queue to add objects to
}

func newController(clientset kubernetes.Clientset, depInformer depInfromer.DeploymentInformer) *controller {
	c := &controller{
		clientset:      &clientset,
		deplLister:     depInformer.Lister(),
		depCacheSynced: depInformer.Informer().HasSynced,
		// Initialising queue
		queue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "myqueue"),
	}

	depInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.handleAdd,
		DeleteFunc: c.handleDelete,
	})
	return c
}

func (c *controller) run(ch <-chan struct{}) {
	fmt.Println("Starting controller")
	// Make sure informer cache is synced succesfully
	// we need to pass it a channel of struct{}
	// if this is not done then something went wrong
	if !cache.WaitForCacheSync(ch, c.depCacheSynced) {
		fmt.Print("\n error waiting for cache to be synced")
	}

	// Until runs untill stop channel is closed
	go wait.Until(c.worker, 1*time.Second, ch)

	// tty stays waiting for input from channel which never comes
	<-ch
}

func (c *controller) worker() {
	// if process item returns bool
	for c.processItem() {

	}
}

func (c *controller) processItem() bool {
	item, shutdown := c.queue.Get()
	if shutdown {
		return false
	}

	key, err := cache.MetaNamespaceKeyFunc(item)
	if err != nil {
		fmt.Println("getting key from cache", err)
	}

	ns, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		fmt.Println("splitting key from cache", err)
	}

	err = c.syncDeployment(ns, name)
	if err != nil {
		// retry
		fmt.Println("syncing deployment", err)
		return false
	}
	c.queue.Done(item)
	return true
}

func (c *controller) syncDeployment(ns string, name string) error {
	if ns == "default" {
		dep, err := c.deplLister.Deployments(ns).Get(name)
		if err != nil {
			fmt.Println(err)
		}
		labels := dep.Labels
		// create service
		svc := corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:   name + "svc",
				Labels: labels,
			},
			Spec: corev1.ServiceSpec{
				Selector: labels,
				Ports: []corev1.ServicePort{
					corev1.ServicePort{
						Name: "http",
						Port: 80,
					},
				},
			},
		}
		_, err = c.clientset.CoreV1().Services(ns).Create(context.Background(), &svc, metav1.CreateOptions{})
		if err != nil {
			fmt.Println(err)
		}
		return nil
	} else {
		fmt.Println("ns is" + ns + "Not making svc")
		return nil
	}
}

func (c *controller) handleAdd(obj interface{}) {
	fmt.Println("Add was called")
	// add to the work queue
	c.queue.Add(obj)

}

func (c *controller) handleDelete(obj interface{}) {
	fmt.Println("Delete was called")
	// add object to queue
	c.queue.Add(obj)
}
