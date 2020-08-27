package main

import (
	"github.com/truecar-ops/terraform/builtin/provisioners/habitat"
	"github.com/truecar-ops/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProvisionerFunc: habitat.Provisioner,
	})
}
