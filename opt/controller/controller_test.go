package controller

import (
	"fmt"
	"k8s.io/apimachinery/pkg/util/wait"
	"testing"
	"time"

	"github.com/k8swatchcrd/opt/pkg/client/clientset/versioned/fake"
	"k8s.io/api/core/v1"
	apiextfake "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/fake"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sfake "k8s.io/client-go/kubernetes/fake"
)

type fixture struct {
	t *testing.T

	client       *fake.Clientset
	kubeclient   *k8sfake.Clientset
	apiextclient *apiextfake.Clientset
}

func (f *fixture) newController() error {
	f.client = fake.NewSimpleClientset()
	f.kubeclient = k8sfake.NewSimpleClientset()
	f.apiextclient = apiextfake.NewSimpleClientset()

	fmt.Println("starting the controller")
	startpod := make(chan bool)
	go Start(f.kubeclient, f.apiextclient, f.client.K8spodcntV1(), true, startpod)
	<-startpod
	fmt.Println("starting the 1st pod")
	p := &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "my-pod"}}
	_, err := f.kubeclient.Core().Pods("default").Create(p)
	if err != nil {
		f.t.Errorf("error injecting pod add: %v", err)
	}
	fmt.Println("starting the 2nd pod")
	p1 := &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "my-pod1"}}
	_, err = f.kubeclient.Core().Pods("default").Create(p1)
	if err != nil {
		f.t.Errorf("error injecting pod add: %v", err)
	}
	if wait := waitForPodcnt(f.client, "podcount", 2, 2); wait != nil {
		return wait
	}
	fmt.Println("Deleting the 2nd pod")
	err = f.kubeclient.Core().Pods("default").Delete("my-pod1", &metav1.DeleteOptions{})
	if err != nil {
		f.t.Errorf("error deleting pod add: %v", err)
	}
	if wait := waitForPodcnt(f.client, "podcount", 2, 1); wait != nil {
		return wait
	}
	return nil
}

func waitForPodcnt(client *fake.Clientset, name string, allcnt int,
	currcnt int) error {
	return wait.Poll(500*time.Millisecond, 60*time.Second, func() (bool, error) {
		crd, err := client.K8spodcntV1().Monitors("default").Get(name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		if allcnt == crd.Status.Allpodcnt && currcnt == crd.Status.Currpodcnt {
			fmt.Println("got the expected pod cnt")
			return true, nil
		}
		return false, nil
	})
}

func (f *fixture) run() {
	f.runController(true, false)
}

func (f *fixture) runController(startInformers bool, expectError bool) {
	f.newController()
}

func newFixture(t *testing.T) *fixture {
	f := &fixture{}
	f.t = t
	return f
}

func TestCRD(t *testing.T) {
	f := newFixture(t)
	f.run()
}
