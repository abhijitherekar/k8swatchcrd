/*
	register.go
	In this file we basically, register the CRD schema and declare it to the
	k8s
*/

package v1

import (
	"reflect"
	/*apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	matav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/scheme"
	"k8s.io/code-generator"
	*/
	"github.com/k8swatchcrd/opt/crd"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

//create a new schema builder and pass in the default scehma
var (
	SchemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	localSchemeBuilder = &SchemeBuilder
	AddToScheme        = SchemeBuilder.AddToScheme
)

// schemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: "k8spodcnt.io", Version: "v1"}

var MonitorResource = crd.CustomResource{
	Name:    "monitor",
	Plural:  "monitors",
	Group:   "k8spodcnt.io",
	Version: "v1",
	Scope:   apiextensionsv1beta1.NamespaceScoped,
	Kind:    reflect.TypeOf(Monitor{}).Name(),
}

func Init() {
	// We only register manually written functions here. The registration of the
	// generated functions takes place in the generated files. The separation
	// makes the code compile even when the generated files are missing.
	localSchemeBuilder.Register(addKnownTypes)
}

// Adds the list of known types to api.Scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Monitor{},
		&MonitorList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
