package spec

import (
	_ "embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace/pkg/sandbox"
	"gopkg.in/yaml.v3"
)

//go:embed example/spec_nested.yaml
var specNestedYaml string

func TestParseSpecInput(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantPath   string
		wantOffset int
		wantErr    bool
	}{
		{
			name:       "path only",
			input:      "example.yaml",
			wantPath:   "example.yaml",
			wantOffset: 0,
			wantErr:    false,
		},
		{
			name:       "path with offset",
			input:      "example.yaml, 2",
			wantPath:   "example.yaml",
			wantOffset: 2,
			wantErr:    false,
		},
		{
			name:       "path with offset no space",
			input:      "example.yaml,2",
			wantPath:   "example.yaml",
			wantOffset: 2,
			wantErr:    false,
		},
		{
			name:       "path with offset multiple spaces",
			input:      "example.yaml,   3",
			wantPath:   "example.yaml",
			wantOffset: 3,
			wantErr:    false,
		},
		{
			name:       "path with zero offset",
			input:      "example.yaml, 0",
			wantPath:   "example.yaml",
			wantOffset: 0,
			wantErr:    false,
		},
		{
			name:       "path with variable",
			input:      "dev-support/run-scripts/${C_ARG0}/spec.yaml",
			wantPath:   "dev-support/run-scripts/${C_ARG0}/spec.yaml",
			wantOffset: 0,
			wantErr:    false,
		},
		{
			name:       "path with variable and offset",
			input:      "dev-support/run-scripts/${C_ARG0}/spec.yaml, 1",
			wantPath:   "dev-support/run-scripts/${C_ARG0}/spec.yaml",
			wantOffset: 1,
			wantErr:    false,
		},
		{
			name:       "path with comma in directory name",
			input:      "some,dir/spec.yaml",
			wantPath:   "some,dir/spec.yaml",
			wantOffset: 0,
			wantErr:    false,
		},
		{
			name:       "path with comma in directory name and offset",
			input:      "some,dir/spec.yaml, 1",
			wantPath:   "some,dir/spec.yaml",
			wantOffset: 1,
			wantErr:    false,
		},
		// Paths with spaces
		{
			name:       "path with spaces",
			input:      "path with spaces/spec.yaml",
			wantPath:   "path with spaces/spec.yaml",
			wantOffset: 0,
			wantErr:    false,
		},
		{
			name:       "path with spaces and offset",
			input:      "path with spaces/spec.yaml, 2",
			wantPath:   "path with spaces/spec.yaml",
			wantOffset: 2,
			wantErr:    false,
		},
		// Comma followed by letters (not a valid offset)
		{
			name:       "comma followed by letters",
			input:      "spec,extra.yaml",
			wantPath:   "spec,extra.yaml",
			wantOffset: 0,
			wantErr:    false,
		},
		{
			name:       "comma space letters",
			input:      "spec, extra.yaml",
			wantPath:   "spec, extra.yaml",
			wantOffset: 0,
			wantErr:    false,
		},
		{
			name:       "comma followed by mixed",
			input:      "spec,2ab.yaml",
			wantPath:   "spec,2ab.yaml",
			wantOffset: 0,
			wantErr:    false,
		},
		// Unicode numerics should NOT be treated as offset
		{
			name:       "unicode fullwidth digit",
			input:      "spec.yaml, \uff12", // U+FF12 is fullwidth digit two
			wantPath:   "spec.yaml, \uff12",
			wantOffset: 0,
			wantErr:    false,
		},
		{
			name:       "unicode arabic-indic digit",
			input:      "spec.yaml, \u0662", // U+0662 is Arabic-Indic digit two
			wantPath:   "spec.yaml, \u0662",
			wantOffset: 0,
			wantErr:    false,
		},
		{
			name:       "unicode superscript digit",
			input:      "spec.yaml, \u00b2", // U+00B2 is superscript two
			wantPath:   "spec.yaml, \u00b2",
			wantOffset: 0,
			wantErr:    false,
		},
		// Multiple commas
		{
			name:       "multiple commas no offset",
			input:      "a,b,c/spec.yaml",
			wantPath:   "a,b,c/spec.yaml",
			wantOffset: 0,
			wantErr:    false,
		},
		{
			name:       "multiple commas with offset",
			input:      "a,b,c/spec.yaml, 5",
			wantPath:   "a,b,c/spec.yaml",
			wantOffset: 5,
			wantErr:    false,
		},
		// Large offset
		{
			name:       "large offset",
			input:      "spec.yaml, 100",
			wantPath:   "spec.yaml",
			wantOffset: 100,
			wantErr:    false,
		},
		// Trailing content after number should not match
		{
			name:       "number followed by text",
			input:      "spec.yaml, 2 extra",
			wantPath:   "spec.yaml, 2 extra",
			wantOffset: 0,
			wantErr:    false,
		},
		{
			name:       "number followed by more path",
			input:      "spec.yaml, 2/more.yaml",
			wantPath:   "spec.yaml, 2/more.yaml",
			wantOffset: 0,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPath, gotOffset, err := parseSpecInput(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseSpecInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPath != tt.wantPath {
				t.Errorf("parseSpecInput() gotPath = %v, want %v", gotPath, tt.wantPath)
			}
			if gotOffset != tt.wantOffset {
				t.Errorf("parseSpecInput() gotOffset = %v, want %v", gotOffset, tt.wantOffset)
			}
		})
	}
}

// TestSpecWithOffset tests that $spec with offset correctly shifts arguments
func TestSpecWithOffset(t *testing.T) {
	// Create a temp directory with the nested spec
	tempDir, err := os.MkdirTemp("", "spec-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Write the nested spec to temp dir
	nestedPath := filepath.Join(tempDir, "nested.yaml")
	if err := os.WriteFile(nestedPath, []byte(specNestedYaml), 0644); err != nil {
		t.Fatal(err)
	}

	// Create parent spec that references the nested spec with offset
	parentSpec := `
name: specparent
description: parent spec for testing shift
completion:
  positional:
    - [cmd1, cmd2, cmd3]
  positionalany:
    - "$spec(` + nestedPath + `, 1)"
`

	var command Command
	if err := yaml.Unmarshal([]byte(parentSpec), &command); err != nil {
		t.Fatal(err)
	}

	sandbox.Command(t, command.ToCobra)(func(s *sandbox.Sandbox) {
		// First positional should complete from parent spec
		s.Run("").
			Expect(carapace.ActionValues(
				"cmd1",
				"cmd2",
				"cmd3",
			))

		// After first arg, nested spec takes over with shift=1
		// So nested spec sees position 0 (first, second, third)
		s.Run("cmd1", "").
			Expect(carapace.ActionValues(
				"first",
				"second",
				"third",
			))

		// After two args, nested spec sees position 1 -> positionalany
		s.Run("cmd1", "first", "").
			Expect(carapace.ActionValues(
				"any1",
				"any2",
				"any3",
			))

		// After three args, still positionalany
		s.Run("cmd1", "first", "any1", "").
			Expect(carapace.ActionValues(
				"any1",
				"any2",
				"any3",
			))
	})
}

// TestSpecWithoutOffset verifies backward compatibility - $spec without offset works as before
func TestSpecWithoutOffset(t *testing.T) {
	// Create a temp directory with the nested spec
	tempDir, err := os.MkdirTemp("", "spec-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Write the nested spec to temp dir
	nestedPath := filepath.Join(tempDir, "nested.yaml")
	if err := os.WriteFile(nestedPath, []byte(specNestedYaml), 0644); err != nil {
		t.Fatal(err)
	}

	// Create parent spec that references the nested spec WITHOUT offset
	parentSpec := `
name: specparent
description: parent spec for testing without offset
completion:
  positional:
    - [cmd1, cmd2, cmd3]
  positionalany:
    - "$spec(` + nestedPath + `)"
`

	var command Command
	if err := yaml.Unmarshal([]byte(parentSpec), &command); err != nil {
		t.Fatal(err)
	}

	sandbox.Command(t, command.ToCobra)(func(s *sandbox.Sandbox) {
		// First positional should complete from parent spec
		s.Run("").
			Expect(carapace.ActionValues(
				"cmd1",
				"cmd2",
				"cmd3",
			))

		// After first arg, nested spec takes over WITHOUT shift
		// So nested spec sees TWO args (cmd1, "") at position 1 -> positionalany
		// because positionalany in parent means we already have 1 arg when entering nested
		s.Run("cmd1", "").
			Expect(carapace.ActionValues(
				"any1",
				"any2",
				"any3",
			))
	})
}
