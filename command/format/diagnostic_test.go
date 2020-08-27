package format

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/hcl/v2"
	"github.com/mitchellh/colorstring"

	"github.com/schiangtc/terraform/tfdiags"
)

func TestDiagnosticWarningsCompact(t *testing.T) {
	var diags tfdiags.Diagnostics
	diags = diags.Append(tfdiags.SimpleWarning("foo"))
	diags = diags.Append(tfdiags.SimpleWarning("foo"))
	diags = diags.Append(tfdiags.SimpleWarning("bar"))
	diags = diags.Append(&hcl.Diagnostic{
		Severity: hcl.DiagWarning,
		Summary:  "source foo",
		Detail:   "...",
		Subject: &hcl.Range{
			Filename: "source.tf",
			Start:    hcl.Pos{Line: 2, Column: 1, Byte: 5},
			End:      hcl.Pos{Line: 2, Column: 1, Byte: 5},
		},
	})
	diags = diags.Append(&hcl.Diagnostic{
		Severity: hcl.DiagWarning,
		Summary:  "source foo",
		Detail:   "...",
		Subject: &hcl.Range{
			Filename: "source.tf",
			Start:    hcl.Pos{Line: 3, Column: 1, Byte: 7},
			End:      hcl.Pos{Line: 3, Column: 1, Byte: 7},
		},
	})
	diags = diags.Append(&hcl.Diagnostic{
		Severity: hcl.DiagWarning,
		Summary:  "source bar",
		Detail:   "...",
		Subject: &hcl.Range{
			Filename: "source2.tf",
			Start:    hcl.Pos{Line: 1, Column: 1, Byte: 1},
			End:      hcl.Pos{Line: 1, Column: 1, Byte: 1},
		},
	})

	// ConsolidateWarnings groups together the ones
	// that have source location information and that
	// have the same summary text.
	diags = diags.ConsolidateWarnings(1)

	// A zero-value Colorize just passes all the formatting
	// codes back to us, so we can test them literally.
	got := DiagnosticWarningsCompact(diags, &colorstring.Colorize{})
	want := `[bold][yellow]Warnings:[reset]

- foo
- foo
- bar
- source foo
  on source.tf line 2 (and 1 more)
- source bar
  on source2.tf line 1
`
	if got != want {
		t.Errorf(
			"wrong result\ngot:\n%s\n\nwant:\n%s\n\ndiff:\n%s",
			got, want, cmp.Diff(want, got),
		)
	}
}

// Test case via https://github.com/schiangtc/terraform/issues/21359
func TestDiagnostic_nonOverlappingHighlightContext(t *testing.T) {
	var diags tfdiags.Diagnostics

	diags = diags.Append(&hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  "Some error",
		Detail:   "...",
		Subject: &hcl.Range{
			Filename: "source.tf",
			Start:    hcl.Pos{Line: 1, Column: 5, Byte: 5},
			End:      hcl.Pos{Line: 1, Column: 5, Byte: 5},
		},
		Context: &hcl.Range{
			Filename: "source.tf",
			Start:    hcl.Pos{Line: 1, Column: 5, Byte: 5},
			End:      hcl.Pos{Line: 4, Column: 2, Byte: 60},
		},
	})
	sources := map[string][]byte{
		"source.tf": []byte(`x = somefunc("testing", {
  alpha = "foo"
  beta  = "bar"
})
`),
	}
	color := &colorstring.Colorize{
		Colors:  colorstring.DefaultColors,
		Reset:   true,
		Disable: true,
	}
	expected := `
Error: Some error

  on source.tf line 1:
   1: x = somefunc("testing", {
   2:   alpha = "foo"
   3:   beta  = "bar"
   4: })

...
`
	output := Diagnostic(diags[0], sources, color, 80)

	if output != expected {
		t.Fatalf("unexpected output: got:\n%s\nwant\n%s\n", output, expected)
	}
}

func TestDiagnostic_emptyOverlapHighlightContext(t *testing.T) {
	var diags tfdiags.Diagnostics

	diags = diags.Append(&hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  "Some error",
		Detail:   "...",
		Subject: &hcl.Range{
			Filename: "source.tf",
			Start:    hcl.Pos{Line: 3, Column: 10, Byte: 38},
			End:      hcl.Pos{Line: 4, Column: 1, Byte: 39},
		},
		Context: &hcl.Range{
			Filename: "source.tf",
			Start:    hcl.Pos{Line: 2, Column: 13, Byte: 27},
			End:      hcl.Pos{Line: 4, Column: 1, Byte: 39},
		},
	})
	sources := map[string][]byte{
		"source.tf": []byte(`variable "x" {
  default = {
    "foo"
  }
`),
	}
	color := &colorstring.Colorize{
		Colors:  colorstring.DefaultColors,
		Reset:   true,
		Disable: true,
	}
	expected := `
Error: Some error

  on source.tf line 3, in variable "x":
   2:   default = {
   3:     "foo"
   4:   }

...
`
	output := Diagnostic(diags[0], sources, color, 80)

	if output != expected {
		t.Fatalf("unexpected output: got:\n%s\nwant\n%s\n", output, expected)
	}
}
