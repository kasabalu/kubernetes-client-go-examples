package main

import (
	"flag"
	"fmt"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmduttil "k8s.io/kubectl/pkg/cmd/util"
)

func main() {
	var res string
	flag.StringVar(&res, "res", "", "resource name ")
	flag.Parse()

	//NewConfigFlags returns ConfigFlags with default values set
	configFlag := genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag()
	matchVersionFlags := cmduttil.NewMatchVersionFlags(configFlag)
	m, err := cmduttil.NewFactory(matchVersionFlags).ToRESTMapper()
	if err != nil {
		fmt.Printf("Error is %s", err.Error())
		return
	}
	gvr, err := m.ResourceFor(schema.GroupVersionResource{
		Resource: res,
	})
	if err != nil {
		fmt.Printf("Error while getting GVR  %s", err.Error())
		return
	}
	fmt.Printf("group %s, version %s, resource %s", gvr.Group, gvr.Version, gvr.Resource)

	fmt.Println()
}
