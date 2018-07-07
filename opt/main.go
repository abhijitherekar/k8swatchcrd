// Package main for the CRD and k8swatch controller
package main

import (
	"fmt"
	"github.com/k8swatchcrd/opt/cmd"
)

func main() {
	fmt.Println("Starting k8s-watcher")
	cmd.Execute()
}
