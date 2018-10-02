package main

import (
	"github.com/Brightspace/terraform-provider-evident/evident"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{ProviderFunc: evident.Provider})
}
