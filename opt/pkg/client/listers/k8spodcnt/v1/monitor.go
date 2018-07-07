/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/k8swatchcrd/opt/pkg/apis/k8spodcnt/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// MonitorLister helps list Monitors.
type MonitorLister interface {
	// List lists all Monitors in the indexer.
	List(selector labels.Selector) (ret []*v1.Monitor, err error)
	// Monitors returns an object that can list and get Monitors.
	Monitors(namespace string) MonitorNamespaceLister
	MonitorListerExpansion
}

// monitorLister implements the MonitorLister interface.
type monitorLister struct {
	indexer cache.Indexer
}

// NewMonitorLister returns a new MonitorLister.
func NewMonitorLister(indexer cache.Indexer) MonitorLister {
	return &monitorLister{indexer: indexer}
}

// List lists all Monitors in the indexer.
func (s *monitorLister) List(selector labels.Selector) (ret []*v1.Monitor, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Monitor))
	})
	return ret, err
}

// Monitors returns an object that can list and get Monitors.
func (s *monitorLister) Monitors(namespace string) MonitorNamespaceLister {
	return monitorNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// MonitorNamespaceLister helps list and get Monitors.
type MonitorNamespaceLister interface {
	// List lists all Monitors in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1.Monitor, err error)
	// Get retrieves the Monitor from the indexer for a given namespace and name.
	Get(name string) (*v1.Monitor, error)
	MonitorNamespaceListerExpansion
}

// monitorNamespaceLister implements the MonitorNamespaceLister
// interface.
type monitorNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Monitors in the indexer for a given namespace.
func (s monitorNamespaceLister) List(selector labels.Selector) (ret []*v1.Monitor, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Monitor))
	})
	return ret, err
}

// Get retrieves the Monitor from the indexer for a given namespace and name.
func (s monitorNamespaceLister) Get(name string) (*v1.Monitor, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("monitor"), name)
	}
	return obj.(*v1.Monitor), nil
}