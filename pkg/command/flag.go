package command

import "fmt"

type Flag struct {
	Longhand   string
	Shorthand  string
	Usage      string
	Repeatable bool
	Optarg     bool
	Value      bool
	Hidden     bool
	Required   bool
	Persistent bool
	Nargs      int
}

func (f Flag) format() string {
	var s string

	if f.Shorthand != "" {
		s += f.Shorthand
		if f.Longhand != "" {
			s += ", "
		}
	}

	if f.Longhand != "" {
		s += f.Longhand
	}

	switch {
	case f.Optarg:
		s += "?"
	case f.Value:
		s += "="
	}

	if f.Repeatable {
		s += "*"
	}

	if f.Required {
		s += "!"
	}

	if f.Hidden {
		s += "&"
	}

	if f.Nargs != 0 {
		s += fmt.Sprintf("{%v}", f.Nargs)
	}

	return s
}
