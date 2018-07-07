package controller

import (
	"fmt"
	"github.com/k8swatchcrd/opt/config"
	"github.com/k8swatchcrd/opt/crd"
	monitor "github.com/k8swatchcrd/opt/pkg/apis/k8spodcnt/v1"
	monitorclient "github.com/k8swatchcrd/opt/pkg/client/clientset/versioned/typed/k8spodcnt/v1"
	api_v1 "k8s.io/api/core/v1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"os"
	"os/signal"
	"syscall"
	"time"
)

/*	for the controller structure we need the following:

	1. A Queue to process the updates on k8s resources
	2. A k8s-clientset to access all the resources of the cluster
	3. A shared index informer to listen on the resources
	4. A logger to log the activities of the controller
*/

type Controller struct {
	KubeClientset    kubernetes.Interface
	ApiExtClientSet  apiextensionsclient.Interface
	MonitorClientSet monitorclient.K8spodcntV1Interface
	PodQueue         workqueue.RateLimitingInterface
	PodInformer      cache.SharedIndexInformer
	//PodIndexer      cache.Indexer
	K8sConfig         *config.Config
	monitorrestclient monitorclient.MonitorInterface
}

/*
	The main func which does the following:
	1. Checks which resource needs to be monitored
	2. Creates a CRD and waits for its init
	3. Then a resource of kind:CRD to watch for pods.
	4. Start the main controller which watches the K8s CORE resources.
*/

func Start(clientset kubernetes.Interface,
	apiclientset apiextensionsclient.Interface,
	monitorclientset monitorclient.K8spodcntV1Interface, fake bool,
	startpod chan bool) {
	k8sconfig := &config.Config{}
	var err error
	//load the config which tells which config to watch for
	if !fake {
		k8sconfig.Resource.Pod = false
		k8sconfig, err = config.New()
		if err != nil {
			fmt.Println("\n error reading the config file")
			panic(err.Error())
		}
	}
	//if pod is true start a PodInformer
	if k8sconfig.Resource.Pod || fake {
		c := NewPodController(clientset, apiclientset,
			monitorclientset)
		stopch := make(chan struct{})

		// 1st Create and wait for CRD resources
		fmt.Println("Registering the monitor resource to the kubernetes")
		resources := []crd.CustomResource{monitor.MonitorResource}
		err = crd.CreateCustomResources(c.ApiExtClientSet, resources, fake)
		if err != nil {
			fmt.Printf("failed to create custom resource. %+v\n", err)
			os.Exit(1)
		}
		// create signals to stop watching the resources
		signalChan := make(chan os.Signal, 1)
		stopChan := make(chan struct{})
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

		// Create the resource of type Monitor
		fmt.Println("Create the monitor resource")
		example := &monitor.Monitor{
			ObjectMeta: meta_v1.ObjectMeta{
				Name:      "podcount",
				Namespace: "default",
			},
			Spec: monitor.MonitorSpec{
				MonitorName: "podcrd",
			},
			Status: monitor.MonitorStatus{
				Allpodcnt:  0,
				Currpodcnt: 0,
			},
		}
		//create a resource of type monitor to watch over the k8s pods
		c.monitorrestclient = monitorclientset.Monitors("default")
		//RESTClient()
		result, err := c.monitorrestclient.Create(example)
		if err == nil {
			fmt.Printf("CREATED: %#v\n", result)
		} else if apierrors.IsAlreadyExists(err) {
			fmt.Printf("ALREADY EXISTS: %#v\n", result)
		} else {
			panic(err)
		}

		//start the podInformer and send 2 channels:
		// 1. stopch to stop the pod queues and informer
		// 2. startpod, to basically signal the test to start creating pods.
		go c.Run(stopch, startpod, fake)
		done := make(chan bool)
		go func() {
			for {
				select {
				case <-signalChan:
					fmt.Println("shutdown signal received, exiting...")
					close(stopChan)
					done <- true
					return
				}
			}
		}()

		<-done
	}
}

/*
	Creates the LIST and WATCH shared informer for the pods.
	Creates a new Controller/informer
	Registers eventHandlers
*/
func NewPodController(client kubernetes.Interface,
	apiclient apiextensionsclient.Interface,
	monitorclient monitorclient.K8spodcntV1Interface) *Controller {
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	informer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
				return client.CoreV1().Pods(meta_v1.NamespaceAll).List(options)
			},
			WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
				return client.CoreV1().Pods(meta_v1.NamespaceAll).Watch(options)
			},
		},
		&api_v1.Pod{},
		0, //Skip resync
		cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
	)

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
		UpdateFunc: func(old, new interface{}) {
			_, err := cache.MetaNamespaceKeyFunc(new)
			if err == nil {
				//queue.Add(key), not needed as we dont want updates.
			}
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
	})
	c := &Controller{
		KubeClientset:    client,
		ApiExtClientSet:  apiclient,
		MonitorClientSet: monitorclient,
		PodQueue:         queue,
		PodInformer:      informer,
		//		PodIndexer:  indexer,
	}

	return c
}

// Run starts the k8swatchcrd controller
func (c *Controller) Run(stopCh <-chan struct{}, startpod chan bool,
	fake bool) {
	defer utilruntime.HandleCrash()
	defer c.PodQueue.ShutDown()

	fmt.Println("Starting  controller")

	go c.PodInformer.Run(stopCh)

	if !cache.WaitForCacheSync(stopCh, c.PodInformer.HasSynced) {
		utilruntime.HandleError(fmt.Errorf("Timed out waiting for cache-syn"))
		return
	}

	fmt.Println("Controller synced and ready")
	if fake {
		startpod <- true
	}
	//wait until we receive a stopch
	wait.Until(c.runWorker, time.Second, stopCh)
}

func (c *Controller) runWorker() {
	for c.ProcessItem() {
	}
}

func (c *Controller) ProcessItem() bool {
	key, quit := c.PodQueue.Get()
	if quit {
		return false
	}
	defer c.PodQueue.Done(key)
	err := c.processPod(key.(string))
	if err == nil {
		c.PodQueue.Forget(key)
	} else {
		c.PodQueue.AddRateLimited(key)
	}
	return true
}

/*
	the main function which basically interacts with
	the Kubeclient and the CRD client to update the status
*/

func (c *Controller) processPod(key string) error {
	_, present, err := c.PodInformer.GetIndexer().GetByKey(key)
	if err != nil {
		return err
	}
	//Once the K8s controller informs of the POD update in the system
	// we need to update the CRD count.
	result, err1 := c.monitorrestclient.Get("podcount", meta_v1.GetOptions{})
	if err1 == nil {
		if !present {
			fmt.Println("Pod deleted with key: ", key)
			result.Status.Currpodcnt--
			_, e := c.monitorrestclient.Update(result)
			if e != nil {
				panic(e)
			}
			return nil
		}
		fmt.Println("Pod added:", key)
		result.Status.Allpodcnt++
		result.Status.Currpodcnt++

		_, e := c.monitorrestclient.Update(result)
		if e != nil {
			panic(e)
		}
	} else {
		fmt.Println("Error while gtting example")
	}

	return nil
}
