package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"terraform-provider-ctyun/internal/provider"
)

var (
	version = "1.0.4"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "www.ctyun.cn/ctyun-it/ctyun",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.NewCtyunProvider(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
