package main

import (
	"github.com/schiangtc/terraform/builtin/provisioners/remote-exec"
	"github.com/schiangtc/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProvisionerFunc: remoteexec.Provisioner,
	})
}
