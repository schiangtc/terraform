package main

import (
	"github.com/schiangtc/terraform/builtin/provisioners/salt-masterless"
	"github.com/schiangtc/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProvisionerFunc: saltmasterless.Provisioner,
	})
}
