package local

import (
	"context"
	"fmt"
	"testing"

	"github.com/schiangtc/terraform/addrs"
	"github.com/schiangtc/terraform/backend"
	"github.com/schiangtc/terraform/configs/configschema"
	"github.com/schiangtc/terraform/internal/initwd"
	"github.com/schiangtc/terraform/providers"
	"github.com/schiangtc/terraform/states"
	"github.com/schiangtc/terraform/terraform"

	"github.com/zclconf/go-cty/cty"
)

func TestLocal_refresh(t *testing.T) {
	b, cleanup := TestLocal(t)
	defer cleanup()

	p := TestLocalProvider(t, b, "test", refreshFixtureSchema())
	testStateFile(t, b.StatePath, testRefreshState())

	p.ReadResourceFn = nil
	p.ReadResourceResponse = providers.ReadResourceResponse{NewState: cty.ObjectVal(map[string]cty.Value{
		"id": cty.StringVal("yes"),
	})}

	op, configCleanup := testOperationRefresh(t, "./testdata/refresh")
	defer configCleanup()

	run, err := b.Operation(context.Background(), op)
	if err != nil {
		t.Fatalf("bad: %s", err)
	}
	<-run.Done()

	if !p.ReadResourceCalled {
		t.Fatal("ReadResource should be called")
	}

	checkState(t, b.StateOutPath, `
test_instance.foo:
  ID = yes
  provider = provider["registry.terraform.io/hashicorp/test"]
	`)

	// the backend should be unlocked after a run
	assertBackendStateUnlocked(t, b)
}

func TestLocal_refreshNoConfig(t *testing.T) {
	b, cleanup := TestLocal(t)
	defer cleanup()
	p := TestLocalProvider(t, b, "test", refreshFixtureSchema())
	testStateFile(t, b.StatePath, testRefreshState())
	p.ReadResourceFn = nil
	p.ReadResourceResponse = providers.ReadResourceResponse{NewState: cty.ObjectVal(map[string]cty.Value{
		"id": cty.StringVal("yes"),
	})}

	op, configCleanup := testOperationRefresh(t, "./testdata/empty")
	defer configCleanup()

	run, err := b.Operation(context.Background(), op)
	if err != nil {
		t.Fatalf("bad: %s", err)
	}
	<-run.Done()

	if !p.ReadResourceCalled {
		t.Fatal("ReadResource should be called")
	}

	checkState(t, b.StateOutPath, `
test_instance.foo:
  ID = yes
  provider = provider["registry.terraform.io/hashicorp/test"]
	`)
}

// GH-12174
func TestLocal_refreshNilModuleWithInput(t *testing.T) {
	b, cleanup := TestLocal(t)
	defer cleanup()
	p := TestLocalProvider(t, b, "test", refreshFixtureSchema())
	testStateFile(t, b.StatePath, testRefreshState())
	p.ReadResourceFn = nil
	p.ReadResourceResponse = providers.ReadResourceResponse{NewState: cty.ObjectVal(map[string]cty.Value{
		"id": cty.StringVal("yes"),
	})}

	b.OpInput = true

	op, configCleanup := testOperationRefresh(t, "./testdata/empty")
	defer configCleanup()

	run, err := b.Operation(context.Background(), op)
	if err != nil {
		t.Fatalf("bad: %s", err)
	}
	<-run.Done()

	if !p.ReadResourceCalled {
		t.Fatal("ReadResource should be called")
	}

	checkState(t, b.StateOutPath, `
test_instance.foo:
  ID = yes
  provider = provider["registry.terraform.io/hashicorp/test"]
	`)
}

func TestLocal_refreshInput(t *testing.T) {
	b, cleanup := TestLocal(t)
	defer cleanup()
	p := TestLocalProvider(t, b, "test", refreshFixtureSchema())
	testStateFile(t, b.StatePath, testRefreshState())

	p.GetSchemaReturn = &terraform.ProviderSchema{
		Provider: &configschema.Block{
			Attributes: map[string]*configschema.Attribute{
				"value": {Type: cty.String, Optional: true},
			},
		},
		ResourceTypes: map[string]*configschema.Block{
			"test_instance": {
				Attributes: map[string]*configschema.Attribute{
					"foo": {Type: cty.String, Optional: true},
					"id":  {Type: cty.String, Optional: true},
				},
			},
		},
	}
	p.ReadResourceFn = nil
	p.ReadResourceResponse = providers.ReadResourceResponse{NewState: cty.ObjectVal(map[string]cty.Value{
		"id": cty.StringVal("yes"),
	})}
	p.ConfigureFn = func(c *terraform.ResourceConfig) error {
		if v, ok := c.Get("value"); !ok || v != "bar" {
			return fmt.Errorf("no value set")
		}

		return nil
	}

	// Enable input asking since it is normally disabled by default
	b.OpInput = true
	b.ContextOpts.UIInput = &terraform.MockUIInput{InputReturnString: "bar"}

	op, configCleanup := testOperationRefresh(t, "./testdata/refresh-var-unset")
	defer configCleanup()
	op.UIIn = b.ContextOpts.UIInput

	run, err := b.Operation(context.Background(), op)
	if err != nil {
		t.Fatalf("bad: %s", err)
	}
	<-run.Done()

	if !p.ReadResourceCalled {
		t.Fatal("ReadResource should be called")
	}

	checkState(t, b.StateOutPath, `
test_instance.foo:
  ID = yes
  provider = provider["registry.terraform.io/hashicorp/test"]
	`)
}

func TestLocal_refreshValidate(t *testing.T) {
	b, cleanup := TestLocal(t)
	defer cleanup()
	p := TestLocalProvider(t, b, "test", refreshFixtureSchema())
	testStateFile(t, b.StatePath, testRefreshState())
	p.ReadResourceFn = nil
	p.ReadResourceResponse = providers.ReadResourceResponse{NewState: cty.ObjectVal(map[string]cty.Value{
		"id": cty.StringVal("yes"),
	})}

	// Enable validation
	b.OpValidation = true

	op, configCleanup := testOperationRefresh(t, "./testdata/refresh")
	defer configCleanup()

	run, err := b.Operation(context.Background(), op)
	if err != nil {
		t.Fatalf("bad: %s", err)
	}
	<-run.Done()

	if !p.PrepareProviderConfigCalled {
		t.Fatal("Prepare provider config should be called")
	}

	checkState(t, b.StateOutPath, `
test_instance.foo:
  ID = yes
  provider = provider["registry.terraform.io/hashicorp/test"]
	`)
}

// This test validates the state lacking behavior when the inner call to
// Context() fails
func TestLocal_refresh_context_error(t *testing.T) {
	b, cleanup := TestLocal(t)
	defer cleanup()
	testStateFile(t, b.StatePath, testRefreshState())
	op, configCleanup := testOperationRefresh(t, "./testdata/apply")
	defer configCleanup()

	// we coerce a failure in Context() by omitting the provider schema

	run, err := b.Operation(context.Background(), op)
	if err != nil {
		t.Fatalf("bad: %s", err)
	}
	<-run.Done()
	if run.Result == backend.OperationSuccess {
		t.Fatal("operation succeeded; want failure")
	}
	assertBackendStateUnlocked(t, b)
}

func testOperationRefresh(t *testing.T, configDir string) (*backend.Operation, func()) {
	t.Helper()

	_, configLoader, configCleanup := initwd.MustLoadConfigForTests(t, configDir)

	return &backend.Operation{
		Type:         backend.OperationTypeRefresh,
		ConfigDir:    configDir,
		ConfigLoader: configLoader,
		LockState:    true,
	}, configCleanup
}

// testRefreshState is just a common state that we use for testing refresh.
func testRefreshState() *states.State {
	state := states.NewState()
	root := state.EnsureModule(addrs.RootModuleInstance)
	root.SetResourceInstanceCurrent(
		mustResourceInstanceAddr("test_instance.foo").Resource,
		&states.ResourceInstanceObjectSrc{
			Status:    states.ObjectReady,
			AttrsJSON: []byte(`{"id":"bar"}`),
		},
		mustProviderConfig(`provider["registry.terraform.io/hashicorp/test"]`),
	)
	return state
}

// refreshFixtureSchema returns a schema suitable for processing the
// configuration in testdata/refresh . This schema should be
// assigned to a mock provider named "test".
func refreshFixtureSchema() *terraform.ProviderSchema {
	return &terraform.ProviderSchema{
		ResourceTypes: map[string]*configschema.Block{
			"test_instance": {
				Attributes: map[string]*configschema.Attribute{
					"ami": {Type: cty.String, Optional: true},
					"id":  {Type: cty.String, Computed: true},
				},
			},
		},
	}
}
