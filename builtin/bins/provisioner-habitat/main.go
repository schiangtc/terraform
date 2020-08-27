package main

import (
	"github.com/schiangtc/terraform/builtin/provisioners/habitat"
	"github.com/schiangtc/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProvisionerFunc: habitat.Provisioner,
	})
}
