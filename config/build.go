package config

import (
	"fmt"
	"strings"
)

type Build struct {
	Version string
	Number  string
}

func (build *Build) String() string {
	return fmt.Sprintf("%v (%v)", build.Version, build.Number)
}

func ParseBuild(s string) (*Build, error) {
	components := strings.Split(s, " ")
	if len(components) != 2 {
		return nil, fmt.Errorf("Provided string does not have required components")
	}

	version := components[0]
	number := strings.Trim(components[1], "() ")
	build := new(Build)
	build.Version = version
	build.Number = number

	return build, nil
}
