package main

import (
	"github.com/truecar-ops/terraform/builtin/provisioners/chef"
	"github.com/truecar-ops/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProvisionerFunc: chef.Provisioner,
	})
}
