package plans

import (
	"testing"

	"github.com/go-test/deep"

	"github.com/schiangtc/terraform/addrs"
)

func TestProviderAddrs(t *testing.T) {

	plan := &Plan{
		VariableValues: map[string]DynamicValue{},
		Changes: &Changes{
			Resources: []*ResourceInstanceChangeSrc{
				{
					Addr: addrs.Resource{
						Mode: addrs.ManagedResourceMode,
						Type: "test_thing",
						Name: "woot",
					}.Instance(addrs.IntKey(0)).Absolute(addrs.RootModuleInstance),
					ProviderAddr: addrs.AbsProviderConfig{
						Module:   addrs.RootModule,
						Provider: addrs.NewLegacyProvider("test"),
					},
				},
				{
					Addr: addrs.Resource{
						Mode: addrs.ManagedResourceMode,
						Type: "test_thing",
						Name: "woot",
					}.Instance(addrs.IntKey(0)).Absolute(addrs.RootModuleInstance),
					DeposedKey: "foodface",
					ProviderAddr: addrs.AbsProviderConfig{
						Module:   addrs.RootModule,
						Provider: addrs.NewLegacyProvider("test"),
					},
				},
				{
					Addr: addrs.Resource{
						Mode: addrs.ManagedResourceMode,
						Type: "test_thing",
						Name: "what",
					}.Instance(addrs.IntKey(0)).Absolute(addrs.RootModuleInstance),
					ProviderAddr: addrs.AbsProviderConfig{
						Module:   addrs.RootModule.Child("foo"),
						Provider: addrs.NewLegacyProvider("test"),
					},
				},
			},
		},
	}

	got := plan.ProviderAddrs()
	want := []addrs.AbsProviderConfig{
		addrs.AbsProviderConfig{
			Module:   addrs.RootModule.Child("foo"),
			Provider: addrs.NewLegacyProvider("test"),
		},
		addrs.AbsProviderConfig{
			Module:   addrs.RootModule,
			Provider: addrs.NewLegacyProvider("test"),
		},
	}

	for _, problem := range deep.Equal(got, want) {
		t.Error(problem)
	}
}
