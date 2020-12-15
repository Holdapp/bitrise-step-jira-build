package config

import (
	"fmt"
)

type Build struct {
	Version string
	Number  string
}

func (build *Build) String() string {
	return fmt.Sprintf("%v (%v)", build.Version, build.Number)
}
