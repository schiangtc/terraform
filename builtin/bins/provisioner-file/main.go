package main

import (
	"github.com/schiangtc/terraform/builtin/provisioners/file"
	"github.com/schiangtc/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProvisionerFunc: file.Provisioner,
	})
}
