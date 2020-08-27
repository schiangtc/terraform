package main

import (
	"github.com/schiangtc/terraform/builtin/providers/test"
	"github.com/schiangtc/terraform/plugin"
	"github.com/schiangtc/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return test.Provider()
		},
	})
}
