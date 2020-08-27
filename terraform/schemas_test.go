package terraform

import (
	"github.com/truecar-ops/terraform/addrs"
	"github.com/truecar-ops/terraform/configs/configschema"
)

func simpleTestSchemas() *Schemas {
	provider := simpleMockProvider()
	provisioner := simpleMockProvisioner()
	return &Schemas{
		Providers: map[addrs.Provider]*ProviderSchema{
			addrs.NewDefaultProvider("test"): provider.GetSchemaReturn,
		},
		Provisioners: map[string]*configschema.Block{
			"test": provisioner.GetSchemaResponse.Provisioner,
		},
	}
}
