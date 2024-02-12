package main

import (
	"context"
	"flag"
	"log"

	"github.com/cloutierMat/terraform-provider-sendbird/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var version = "dev"

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "Set to true to run in debug mode.")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/cloutierMat/sendbird",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
