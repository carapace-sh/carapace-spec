package spec

import (
	"fmt"
	"maps"
	"os/exec"
	"regexp"
	"strings"
)

// TODO experimental - internal use
func ScanMacros(pkgs ...string) (MacroMap, error) {
	result := make(MacroMap)
	for _, pkg := range pkgs {
		output, err := exec.Command("go", "list", pkg+"/...").Output()
		if err != nil {
			if err, ok := err.(*exec.ExitError); ok {
				println(string(err.Stderr))
			}
			return nil, err
		}

		for _, subPkg := range strings.Split(string(output), "\n") {
			if pkg == subPkg {
				println("skipping " + pkg)
				continue // TODO re-enable
			}
			prefix := strings.TrimPrefix(subPkg, pkg)
			prefix = strings.TrimLeft(prefix, "/")
			prefix = strings.ReplaceAll(prefix, "/", ".")
			pkgMacros, err := scan(subPkg, prefix)
			if err != nil {
				println("scan failed")
				return nil, err
			}
			maps.Copy(result, pkgMacros)
		}
	}
	return result, nil
}

func scan(pkg, prefix string) (MacroMap, error) {
	result := make(MacroMap, 0)
	output, err := exec.Command("go", "doc", "-all", pkg).Output()
	if err != nil {
		if err, ok := err.(*exec.ExitError); ok {
			println(string(err.Stderr))
		}
		return result, nil
		// return nil, err // TODO some packages just don't return results - need to differentiate between real error (if possible)
	}

	rPackage := regexp.MustCompile(`^package (?P<name>[^ ]+) \/\/ import "(?P<package>.+)"$`)
	rFunction := regexp.MustCompile(`^func Action(?P<name>[^(]+)\((?P<args>.*)\) carapace.Action$`)

	var functionPrefix string

	var currentMacro *Macro
	for index, line := range strings.Split(string(output), "\n") {
		if index == 0 {
			matches := rPackage.FindStringSubmatch(line)
			if matches == nil {
				return result, nil
				// return nil, errors.New("failed to determine package") // TODO sometimes returns no package?
			}
			functionPrefix = matches[2]
		}
		if matches := rFunction.FindStringSubmatch(line); matches != nil {
			if currentMacro != nil {
				result[currentMacro.Name] = *currentMacro
			}
			currentMacro = &Macro{
				Name:     matches[1],
				Function: fmt.Sprintf("%s#Action%s", functionPrefix, matches[1]),
				Args:     matches[2],
			}
			if prefix != "" {
				currentMacro.Name = fmt.Sprintf("%s.%s", prefix, currentMacro.Name)
			}
			continue
		} else if currentMacro != nil && line != "" && !strings.HasPrefix(line, " ") {
			result[currentMacro.Name] = *currentMacro
			currentMacro = nil
		}

		if currentMacro == nil {
			continue
		}

		switch {
		case strings.HasPrefix(line, "        "):
			if currentMacro.Example != "" {
				currentMacro.Example += "\n"
			}
			currentMacro.Example += strings.TrimPrefix(line, "        ")

		case strings.HasPrefix(line, "    "):
			if currentMacro.Description != "" {
				currentMacro.Description += " "
			}
			currentMacro.Description += strings.TrimLeft(line, " ")
		}
	}
	if currentMacro != nil {
		result[currentMacro.Name] = *currentMacro
	}
	return result, nil
}
