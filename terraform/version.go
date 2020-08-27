package terraform

import (
	"github.com/truecar-ops/terraform/version"
)

// Deprecated: Providers should use schema.Provider.TerraformVersion instead
func VersionString() string {
	return version.String()
}
