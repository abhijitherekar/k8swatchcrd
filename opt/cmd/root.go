package cmd

import (
	"fmt"
	"github.com/k8swatchcrd/opt/controller"
	"github.com/spf13/cobra"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	monitorclient "github.com/k8swatchcrd/opt/pkg/client/clientset/versioned/typed/k8spodcnt/v1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "k8swatcher",
	Short: "watches the pods getting created",
	Long:  `A k8s watcher to watch for the pods getting created`,
	Run:   run_k8swatcher,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

/*	1. Now here at the run_k8swatcher, we need to create a kubernetes client
	this is got from the ~/.kube/config or from the incluster config

	2. We need start the new custom controller which will listen to the
	pods
*/

func run_k8swatcher(cmd *cobra.Command, args []string) {
	var config *rest.Config
	var err error
	kubeConfigPath := os.Getenv("HOME") + "/.kube/config"

	if _, err := os.Stat(kubeConfigPath); err == nil {
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	} else {
		config, err = rest.InClusterConfig()
	}

	if err != nil {
		fmt.Println("err while building")
		return
	}
	//get the k8s clientset to interact to the API server for core resources
	kubeclientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	//get the clientset for the interacting with the CRD creation
	apiextClientSet, err := apiextensionsclient.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	//get the clientset for intercting with our CR
	monitorClientSet, err := monitorclient.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	//Now start the controller which basically watches the core
	//resources and updates the CRD accordingly
	controller.Start(kubeclientset, apiextClientSet, monitorClientSet, false, nil)
	//wait forever
	select {}
}
