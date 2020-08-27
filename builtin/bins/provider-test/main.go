package main

import (
	"github.com/truecar-ops/terraform/builtin/providers/test"
	"github.com/truecar-ops/terraform/plugin"
	"github.com/truecar-ops/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return test.Provider()
		},
	})
}
