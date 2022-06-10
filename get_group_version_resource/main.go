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
	// WithDeprecatedPasswordFlag() enables usernmae, pwd(empty ones)
	configFlag := genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag()
	//ConfigFlags composes the set of values necessary for obtaining a REST client

	// To get the RestMapperInterface , we have RESTClientGetter Interface.
	//RESTClientGetter Interface has ToRESTMapper() method that will return RestMapper interface.
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
