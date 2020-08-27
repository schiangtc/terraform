package terraform

import (
	"testing"

	"github.com/schiangtc/terraform/providers"

	backendInit "github.com/schiangtc/terraform/backend/init"
)

var testAccProviders map[string]*Provider
var testAccProvider *Provider

func init() {
	// Initialize the backends
	backendInit.Init(nil)

	testAccProvider = NewProvider()
	testAccProviders = map[string]*Provider{
		"terraform": testAccProvider,
	}
}

func TestProvider_impl(t *testing.T) {
	var _ providers.Interface = NewProvider()
}

func testAccPreCheck(t *testing.T) {
}
