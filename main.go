package main

import (
	"context"
	"flag"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var (
	version = "1.1.1"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/ctyun-it/ctyun",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), service.NewCtyunProvider(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
