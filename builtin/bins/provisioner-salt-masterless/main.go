package main

import (
	"github.com/truecar-ops/terraform/builtin/provisioners/salt-masterless"
	"github.com/truecar-ops/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProvisionerFunc: saltmasterless.Provisioner,
	})
}
