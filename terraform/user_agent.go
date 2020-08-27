package terraform

import (
	"github.com/truecar-ops/terraform/httpclient"
)

// Generate a UserAgent string
//
// Deprecated: Use httpclient.UserAgent(version) instead
func UserAgentString() string {
	return httpclient.UserAgentString()
}
