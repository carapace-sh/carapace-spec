package shebang

import (
	"errors"
	"regexp"
	"strings"

	"github.com/carapace-sh/carapace-shlex"
)

type shebang struct {
	Command string   // interpreter
	Args    []string // optional arguments (deriving from the standard and allowing more than one)
	Script  string   // script without shebang header for compability
}

func Parse(s string) (*shebang, error) {
	firstLine, script, ok := strings.Cut(s, "\n")
	if !ok {
		return nil, errors.New("missing shebang header")
	}

	re := regexp.MustCompile(`^#!(?P<command>[^ ]+)( (?P<arg>.*))?$`)
	matches := re.FindStringSubmatch(strings.TrimSpace(firstLine))
	if matches == nil {
		return nil, errors.New("invalid shebang header")
	}

	shebang := &shebang{
		Command: matches[1],
		Args:    []string{},
		Script:  script,
	}
	if matches[3] != "" {
		tokens, err := shlex.Split(matches[3])
		if err != nil {
			return nil, err
		}
		shebang.Args = tokens.Words().Strings() // optional args
	}

	return shebang, nil
}
