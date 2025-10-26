package spec

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"go/format"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/carapace-sh/carapace-spec/internal/pflagfork"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type codegenCmd struct {
	cmd *cobra.Command
}

func (s codegenCmd) formatHeader() string {
	return `package cmd
import (
	"github.com/carapace-sh/carapace"
	"github.com/spf13/cobra"
)
`
}

func (s codegenCmd) formatGroups() string {
	if len(s.cmd.Groups()) == 0 {
		return ""
	}

	groups := make([]string, 0)
	for _, group := range s.cmd.Groups() {
		groups = append(groups, fmt.Sprintf(`&cobra.Group{ID: %#v, Title: %#v},`, group.ID, group.Title))
	}
	return fmt.Sprintf("%vCmd.AddGroup(\n%v\n)\n", cmdVarName(s.cmd), strings.Join(groups, "\n"))
}

func (s codegenCmd) formatCommand() string {
	snippet := fmt.Sprintf(
		`var %vCmd = &cobra.Command{
	Use:     %#v,
	Short:   %#v,
	GroupID: %#v,
	Aliases: []string{"%v"},
	Hidden:  %v,
	Run:     func(cmd *cobra.Command, args []string) {},
}
`, cmdVarName(s.cmd), strings.SplitN(s.cmd.Use, "\n", 2)[0], s.cmd.Short, s.cmd.GroupID, strings.Join(s.cmd.Aliases, `", "`), s.cmd.Hidden)

	if s.cmd.GroupID == "" {
		re := regexp.MustCompile("(?m)\n\tGroupID:.*$")
		snippet = re.ReplaceAllString(snippet, "")

	}

	if len(s.cmd.Aliases) == 0 {
		re := regexp.MustCompile("(?m)\n\t+Aliases:.*$")
		snippet = re.ReplaceAllString(snippet, "")
	}

	if !s.cmd.Hidden {
		re := regexp.MustCompile("(?m)\n\tHidden:.*$")
		snippet = re.ReplaceAllString(snippet, "")
	}

	return snippet
}

func (s codegenCmd) formatExecute() string {
	if s.cmd.HasParent() {
		return ""
	}
	return `func Execute() error {
	return rootCmd.Execute()
}
`
}
func Codegen(cmd *cobra.Command) error {
	dir, err := os.MkdirTemp(os.TempDir(), "carapace-codegen-")
	if err != nil {
		return err
	}

	return codegen(cmd, dir)
}

func codegen(cmd *cobra.Command, tmpDir string) error {
	out := &bytes.Buffer{}
	fmt.Fprintln(out, codegenCmd{cmd}.formatHeader())
	fmt.Fprintln(out, codegenCmd{cmd}.formatCommand())
	fmt.Fprintln(out, codegenCmd{cmd}.formatExecute())

	fmt.Fprintf(out, `func init() {
	carapace.Gen(%vCmd).Standalone()
%v
`, cmdVarName(cmd), codegenCmd{cmd}.formatGroups())

	cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
		if f.Deprecated != "" {
			return
		}

		persistentPrefix := ""
		if cmd.PersistentFlags().Lookup(f.Name) != nil {
			persistentPrefix = "Persistent"
		}

		switch (pflagfork.Flag{Flag: f}).Mode() {
		case pflagfork.ShorthandOnly:
			fmt.Fprintf(out, `	%vCmd.%vFlags().%vS("%v", "%v", %v, %v)`+"\n", cmdVarName(cmd), persistentPrefix, flagType(f), f.Name, f.Shorthand, flagValue(f), formatUsage(f.Usage))
		case pflagfork.NameAsShorthand:
			fmt.Fprintf(out, `	%vCmd.%vFlags().%vN("%v", "%v", %v, %v)`+"\n", cmdVarName(cmd), persistentPrefix, flagType(f), f.Name, f.Shorthand, flagValue(f), formatUsage(f.Usage))
		case pflagfork.Default:
			switch {
			case f.Shorthand != "" && f.Value.Type() == "count":
				fmt.Fprintf(out, `	%vCmd.%vFlags().%vP("%v", "%v", %v)`+"\n", cmdVarName(cmd), persistentPrefix, flagType(f), f.Name, f.Shorthand, formatUsage(f.Usage))
			case f.Shorthand != "" && f.Value.Type() != "count":
				fmt.Fprintf(out, `	%vCmd.%vFlags().%vP("%v", "%v", %v, %v)`+"\n", cmdVarName(cmd), persistentPrefix, flagType(f), f.Name, f.Shorthand, flagValue(f), formatUsage(f.Usage))
			case f.Value.Type() == "count":
				fmt.Fprintf(out, `	%vCmd.%vFlags().%v("%v", %v)`+"\n", cmdVarName(cmd), persistentPrefix, flagType(f), f.Name, formatUsage(f.Usage))
			default:
				fmt.Fprintf(out, `	%vCmd.%vFlags().%v("%v", %v, %v)`+"\n", cmdVarName(cmd), persistentPrefix, flagType(f), f.Name, flagValue(f), formatUsage(f.Usage))
			}
		}
	})

	cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
		if f.Deprecated != "" {
			return
		}

		if f.Value.Type() != "bool" && f.Value.Type() != "count" && f.NoOptDefVal != "" {
			fmt.Fprintf(out, `    %vCmd.Flag("%v").NoOptDefVal = "%v"`+"\n", cmdVarName(cmd), f.Name, f.NoOptDefVal)
		}

		if f.Hidden {
			fmt.Fprintf(out, `    %vCmd.Flag("%v").Hidden = true`+"\n", cmdVarName(cmd), f.Name)
		}

		if annotation := f.Annotations[cobra.BashCompOneRequiredFlag]; len(annotation) == 1 && annotation[0] == "true" {
			fmt.Fprintf(out, `    %vCmd.MarkFlagRequired("%v")`+"\n", cmdVarName(cmd), f.Name)
		}
	})

	if cmd.HasParent() {
		fmt.Fprintf(out, `	%vCmd.AddCommand(%vCmd)`+"\n", cmdVarName(cmd.Parent()), cmdVarName(cmd))
	}

	fmt.Fprintln(out, "}")

	filename := fmt.Sprintf(`%v/%v.go`, tmpDir, cmdVarName(cmd))

	println(filename)
	formatted, err := format.Source(out.Bytes())
	if err != nil {
		unformatted := strings.Split(out.String(), "\n")
		if line, err := strconv.Atoi(strings.SplitN(err.Error(), ":", 2)[0]); err == nil {
			unformatted[line-1] = "\033[31m" + unformatted[line-1] + "\033[2;37m"
		}
		println("\033[2;37m" + strings.Join(unformatted, "\n") + "\033[0m")
		return err
	}

	os.WriteFile(filename, formatted, 0644)

	for _, subcmd := range cmd.Commands() {
		if subcmd.Deprecated == "" && subcmd.Name() != "_carapace" {
			if err := codegen(subcmd, tmpDir); err != nil {
				return err
			}
		}
	}
	return nil
}

func formatUsage(usage string) string {
	return fmt.Sprintf("%q", strings.Split(usage, "\n")[0])
}

func cmdVarName(cmd *cobra.Command) string {
	if !cmd.HasParent() {
		return "root"
	}
	return strings.TrimPrefix(fmt.Sprintf(`%v_%v`, cmdVarName(cmd.Parent()), normalizeVarName(cmd.Name())), "root_")
}

func normalizeVarName(s string) string {
	normalized := make([]string, 0)
	capitalize := false

	for _, c := range s {
		switch {
		case c == '-' || c == ':':
			capitalize = true
			continue
		case capitalize:
			normalized = append(normalized, strings.ToUpper(string(c)))
			capitalize = false
		default:
			normalized = append(normalized, string(c))
		}
	}
	return strings.Join(normalized, "")
}

func flagType(f *pflag.Flag) string {
	switch f.Value.Type() {
	case
		"bool",
		"boolSlice",
		"bytesBase64",
		"bytesHex",
		"count",
		"custom",
		"duration",
		"durationSlice",
		"flagVar",
		"float32",
		"float32Slice",
		"float64",
		"float64Slice",
		"int",
		"int16",
		"int32",
		"int32Slice",
		"int64",
		"int64Slice",
		"int8",
		"intSlice",
		"ip",
		"ipMask",
		"ipNet",
		"ipNetSlice",
		"ipSlice",
		"string",
		"stringArray",
		"strings",
		"stringSlice",
		"uint",
		"uint16",
		"uint32",
		"uint64",
		"uint8",
		"uintSlice",
		"version":
		return strings.ToUpper(f.Value.Type()[:1]) + f.Value.Type()[1:]
	case
		"stringToInt",
		"stringToInt64",
		"stringToString":
		return "String"
	default:
		return "String"
	}
}

func flagValue(f *pflag.Flag) string {
	if strings.HasSuffix(f.Value.Type(), "Slice") ||
		strings.HasSuffix(f.Value.Type(), "Array") {
		if strings.HasPrefix(f.Value.Type(), "string") {
			if f.Value.String() == "[]" {
				return "nil"
			}

			vals, _ := csv.NewReader(strings.NewReader(f.Value.String()[1 : len(f.Value.String())-1])).Read()
			formatted := strings.Join(vals, `", "`)
			if len(formatted) > 0 {
				formatted = fmt.Sprintf(`"%v"`, formatted)
			}
			return fmt.Sprintf(`[]string{%v}`, formatted)
		}
		return fmt.Sprintf(`[]%v{%v}`, strings.TrimSuffix(strings.TrimSuffix(f.Value.Type(), "Slice"), "Array"), f.Value.String()[1:len(f.Value.String())-1])
	}

	switch f.Value.Type() {
	case "string":
		return fmt.Sprintf(`"%v"`, f.Value.String())
	case
		"duration",
		"float32",
		"float64",
		"int",
		"int8",
		"int16",
		"int32",
		"int64",
		"uint",
		"uint16",
		"uint32",
		"uint64",
		"uint8":
		return "0"
	case
		"bool",
		"boolSlice",
		"bytesBase64",
		"bytesHex",
		"count",
		"custom",
		"durationSlice",
		"flagVar",
		"float32Slice",
		"float64Slice",
		"int32Slice",
		"int64Slice",
		"intSlice",
		"ip",
		"ipMask",
		"ipNet",
		"ipNetSlice",
		"ipSlice",
		"stringArray",
		"strings",
		"stringSlice",
		"uintSlice",
		"version":
		return f.Value.String()
	case
		"stringToInt",
		"stringToInt64",
		"stringToString":
		return `""`
	default:
		return `""`
	}
}
