package main

import (
	"github.com/truecar-ops/terraform/builtin/provisioners/file"
	"github.com/truecar-ops/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProvisionerFunc: file.Provisioner,
	})
}
