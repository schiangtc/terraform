package main

import (
	"github.com/schiangtc/terraform/builtin/provisioners/chef"
	"github.com/schiangtc/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProvisionerFunc: chef.Provisioner,
	})
}
