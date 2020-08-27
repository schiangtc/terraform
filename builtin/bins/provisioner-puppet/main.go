package main

import (
	"github.com/truecar-ops/terraform/builtin/provisioners/puppet"
	"github.com/truecar-ops/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProvisionerFunc: puppet.Provisioner,
	})
}
