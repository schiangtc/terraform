package command

import (
	"context"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mitchellh/cli"

	svchost "github.com/hashicorp/terraform-svchost"
	"github.com/hashicorp/terraform-svchost/disco"
	"github.com/truecar-ops/terraform/command/cliconfig"
	oauthserver "github.com/truecar-ops/terraform/command/testdata/login-oauth-server"
	tfeserver "github.com/truecar-ops/terraform/command/testdata/login-tfe-server"
	"github.com/truecar-ops/terraform/command/webbrowser"
	"github.com/truecar-ops/terraform/httpclient"
	"github.com/truecar-ops/terraform/version"
)

func TestLogin(t *testing.T) {
	// oauthserver.Handler is a stub OAuth server implementation that will,
	// on success, always issue a bearer token named "good-token".
	s := httptest.NewServer(oauthserver.Handler)
	defer s.Close()

	// tfeserver.Handler is a stub TFE API implementation which will respond
	// to ping and current account requests, when requests are authenticated
	// with token "good-token"
	ts := httptest.NewServer(tfeserver.Handler)
	defer ts.Close()

	loginTestCase := func(test func(t *testing.T, c *LoginCommand, ui *cli.MockUi)) func(t *testing.T) {
		return func(t *testing.T) {
			t.Helper()
			workDir, err := ioutil.TempDir("", "terraform-test-command-login")
			if err != nil {
				t.Fatalf("cannot create temporary directory: %s", err)
			}
			defer os.RemoveAll(workDir)

			// We'll use this context to avoid asynchronous tasks outliving
			// a single test run.
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// Do not use the NewMockUi initializer here, as we want to delay
			// the call to init until after setting up the input mocks
			ui := new(cli.MockUi)

			browserLauncher := webbrowser.NewMockLauncher(ctx)
			creds := cliconfig.EmptyCredentialsSourceForTests(filepath.Join(workDir, "credentials.tfrc.json"))
			svcs := disco.NewWithCredentialsSource(creds)
			svcs.SetUserAgent(httpclient.TerraformUserAgent(version.String()))

			svcs.ForceHostServices(svchost.Hostname("app.terraform.io"), map[string]interface{}{
				"login.v1": map[string]interface{}{
					// On app.terraform.io we use password-based authorization.
					// That's the only hostname that it's permitted for, so we can't
					// use a fake hostname here.
					"client":      "terraformcli",
					"token":       s.URL + "/token",
					"grant_types": []interface{}{"password"},
				},
			})
			svcs.ForceHostServices(svchost.Hostname("example.com"), map[string]interface{}{
				"login.v1": map[string]interface{}{
					// For this fake hostname we'll use a conventional OAuth flow,
					// with browser-based consent that we'll mock away using a
					// mock browser launcher below.
					"client": "anything-goes",
					"authz":  s.URL + "/authz",
					"token":  s.URL + "/token",
				},
			})
			svcs.ForceHostServices(svchost.Hostname("tfe.acme.com"), map[string]interface{}{
				// This represents a Terraform Enterprise instance which does not
				// yet support the login API, but does support the TFE tokens API.
				"tfe.v2":   ts.URL + "/api/v2",
				"tfe.v2.1": ts.URL + "/api/v2",
				"tfe.v2.2": ts.URL + "/api/v2",
			})
			svcs.ForceHostServices(svchost.Hostname("unsupported.example.net"), map[string]interface{}{
				// This host intentionally left blank.
			})

			c := &LoginCommand{
				Meta: Meta{
					Ui:              ui,
					BrowserLauncher: browserLauncher,
					Services:        svcs,
				},
			}

			test(t, c, ui)
		}
	}

	t.Run("defaulting to app.terraform.io with password flow", loginTestCase(func(t *testing.T, c *LoginCommand, ui *cli.MockUi) {
		defer testInputMap(t, map[string]string{
			"approve":  "yes",
			"username": "foo",
			"password": "bar",
		})()
		status := c.Run(nil)
		if status != 0 {
			t.Fatalf("unexpected error code %d\nstderr:\n%s", status, ui.ErrorWriter.String())
		}

		credsSrc := c.Services.CredentialsSource()
		creds, err := credsSrc.ForHost(svchost.Hostname("app.terraform.io"))
		if err != nil {
			t.Errorf("failed to retrieve credentials: %s", err)
		}
		if got, want := creds.Token(), "good-token"; got != want {
			t.Errorf("wrong token %q; want %q", got, want)
		}
	}))

	t.Run("example.com with authorization code flow", loginTestCase(func(t *testing.T, c *LoginCommand, ui *cli.MockUi) {
		// Enter "yes" at the consent prompt.
		defer testInputMap(t, map[string]string{
			"approve": "yes",
		})()
		status := c.Run([]string{"example.com"})
		if status != 0 {
			t.Fatalf("unexpected error code %d\nstderr:\n%s", status, ui.ErrorWriter.String())
		}

		credsSrc := c.Services.CredentialsSource()
		creds, err := credsSrc.ForHost(svchost.Hostname("example.com"))
		if err != nil {
			t.Errorf("failed to retrieve credentials: %s", err)
		}
		if got, want := creds.Token(), "good-token"; got != want {
			t.Errorf("wrong token %q; want %q", got, want)
		}
	}))

	t.Run("TFE host without login support", loginTestCase(func(t *testing.T, c *LoginCommand, ui *cli.MockUi) {
		// Enter "yes" at the consent prompt, then paste a token with some
		// accidental whitespace.
		defer testInputMap(t, map[string]string{
			"approve": "yes",
			"token":   "  good-token ",
		})()
		status := c.Run([]string{"tfe.acme.com"})
		if status != 0 {
			t.Fatalf("unexpected error code %d\nstderr:\n%s", status, ui.ErrorWriter.String())
		}

		credsSrc := c.Services.CredentialsSource()
		creds, err := credsSrc.ForHost(svchost.Hostname("tfe.acme.com"))
		if err != nil {
			t.Errorf("failed to retrieve credentials: %s", err)
		}
		if got, want := creds.Token(), "good-token"; got != want {
			t.Errorf("wrong token %q; want %q", got, want)
		}
	}))

	t.Run("TFE host without login support, incorrectly pasted token", loginTestCase(func(t *testing.T, c *LoginCommand, ui *cli.MockUi) {
		// Enter "yes" at the consent prompt, then paste an invalid token.
		defer testInputMap(t, map[string]string{
			"approve": "yes",
			"token":   "good-tok",
		})()
		status := c.Run([]string{"tfe.acme.com"})
		if status != 1 {
			t.Fatalf("unexpected error code %d\nstderr:\n%s", status, ui.ErrorWriter.String())
		}

		credsSrc := c.Services.CredentialsSource()
		creds, err := credsSrc.ForHost(svchost.Hostname("tfe.acme.com"))
		if err != nil {
			t.Errorf("failed to retrieve credentials: %s", err)
		}
		if creds != nil {
			t.Errorf("wrong token %q; should have no token", creds.Token())
		}
	}))

	t.Run("host without login or TFE API support", loginTestCase(func(t *testing.T, c *LoginCommand, ui *cli.MockUi) {
		status := c.Run([]string{"unsupported.example.net"})
		if status == 0 {
			t.Fatalf("successful exit; want error")
		}

		if got, want := ui.ErrorWriter.String(), "Error: Host does not support Terraform tokens API"; !strings.Contains(got, want) {
			t.Fatalf("missing expected error message\nwant: %s\nfull output:\n%s", want, got)
		}
	}))

	t.Run("answering no cancels", loginTestCase(func(t *testing.T, c *LoginCommand, ui *cli.MockUi) {
		// Enter "no" at the consent prompt
		defer testInputMap(t, map[string]string{
			"approve": "no",
		})()
		status := c.Run(nil)
		if status != 1 {
			t.Fatalf("unexpected error code %d\nstderr:\n%s", status, ui.ErrorWriter.String())
		}

		if got, want := ui.ErrorWriter.String(), "Login cancelled"; !strings.Contains(got, want) {
			t.Fatalf("missing expected error message\nwant: %s\nfull output:\n%s", want, got)
		}
	}))

	t.Run("answering y cancels", loginTestCase(func(t *testing.T, c *LoginCommand, ui *cli.MockUi) {
		// Enter "y" at the consent prompt
		defer testInputMap(t, map[string]string{
			"approve": "y",
		})()
		status := c.Run(nil)
		if status != 1 {
			t.Fatalf("unexpected error code %d\nstderr:\n%s", status, ui.ErrorWriter.String())
		}

		if got, want := ui.ErrorWriter.String(), "Login cancelled"; !strings.Contains(got, want) {
			t.Fatalf("missing expected error message\nwant: %s\nfull output:\n%s", want, got)
		}
	}))
}
