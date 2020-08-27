package main

import (
	"github.com/schiangtc/terraform/builtin/provisioners/puppet"
	"github.com/schiangtc/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProvisionerFunc: puppet.Provisioner,
	})
}
