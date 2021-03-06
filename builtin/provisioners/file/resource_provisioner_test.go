package file

import (
	"testing"

	"github.com/truecar-ops/terraform/configs/hcl2shim"
	"github.com/truecar-ops/terraform/helper/schema"
	"github.com/truecar-ops/terraform/terraform"
)

func TestResourceProvisioner_impl(t *testing.T) {
	var _ terraform.ResourceProvisioner = Provisioner()
}

func TestProvisioner(t *testing.T) {
	if err := Provisioner().(*schema.Provisioner).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestResourceProvider_Validate_good_source(t *testing.T) {
	c := testConfig(t, map[string]interface{}{
		"source":      "/tmp/foo",
		"destination": "/tmp/bar",
	})

	warn, errs := Provisioner().Validate(c)
	if len(warn) > 0 {
		t.Fatalf("Warnings: %v", warn)
	}
	if len(errs) > 0 {
		t.Fatalf("Errors: %v", errs)
	}
}

func TestResourceProvider_Validate_good_content(t *testing.T) {
	c := testConfig(t, map[string]interface{}{
		"content":     "value to copy",
		"destination": "/tmp/bar",
	})

	warn, errs := Provisioner().Validate(c)
	if len(warn) > 0 {
		t.Fatalf("Warnings: %v", warn)
	}
	if len(errs) > 0 {
		t.Fatalf("Errors: %v", errs)
	}
}

func TestResourceProvider_Validate_good_unknown_variable_value(t *testing.T) {
	c := testConfig(t, map[string]interface{}{
		"content":     hcl2shim.UnknownVariableValue,
		"destination": "/tmp/bar",
	})

	warn, errs := Provisioner().Validate(c)
	if len(warn) > 0 {
		t.Fatalf("Warnings: %v", warn)
	}
	if len(errs) > 0 {
		t.Fatalf("Errors: %v", errs)
	}
}

func TestResourceProvider_Validate_bad_not_destination(t *testing.T) {
	c := testConfig(t, map[string]interface{}{
		"source": "nope",
	})

	warn, errs := Provisioner().Validate(c)
	if len(warn) > 0 {
		t.Fatalf("Warnings: %v", warn)
	}
	if len(errs) == 0 {
		t.Fatalf("Should have errors")
	}
}

func TestResourceProvider_Validate_bad_no_source(t *testing.T) {
	c := testConfig(t, map[string]interface{}{
		"destination": "/tmp/bar",
	})

	warn, errs := Provisioner().Validate(c)
	if len(warn) > 0 {
		t.Fatalf("Warnings: %v", warn)
	}
	if len(errs) == 0 {
		t.Fatalf("Should have errors")
	}
}

func TestResourceProvider_Validate_bad_to_many_src(t *testing.T) {
	c := testConfig(t, map[string]interface{}{
		"source":      "nope",
		"content":     "value to copy",
		"destination": "/tmp/bar",
	})

	warn, errs := Provisioner().Validate(c)
	if len(warn) > 0 {
		t.Fatalf("Warnings: %v", warn)
	}
	if len(errs) == 0 {
		t.Fatalf("Should have errors")
	}
}

func testConfig(t *testing.T, c map[string]interface{}) *terraform.ResourceConfig {
	return terraform.NewResourceConfigRaw(c)
}
