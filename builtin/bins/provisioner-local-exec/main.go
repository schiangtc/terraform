package main

import (
	"github.com/schiangtc/terraform/builtin/provisioners/local-exec"
	"github.com/schiangtc/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProvisionerFunc: localexec.Provisioner,
	})
}
