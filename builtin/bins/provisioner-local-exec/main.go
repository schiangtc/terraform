package main

import (
	"github.com/truecar-ops/terraform/builtin/provisioners/local-exec"
	"github.com/truecar-ops/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProvisionerFunc: localexec.Provisioner,
	})
}
