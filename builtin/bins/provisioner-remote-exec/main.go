package main

import (
	"github.com/truecar-ops/terraform/builtin/provisioners/remote-exec"
	"github.com/truecar-ops/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProvisionerFunc: remoteexec.Provisioner,
	})
}
