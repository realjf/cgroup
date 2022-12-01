package cgroup

import (
	"errors"

	v1 "github.com/containerd/cgroups/stats/v1"
	"github.com/containerd/cgroups/v3/cgroup2/stats"
)

// Percent is a percentage value. It is used to specify CPU rate limits.
type Percent uint

// Memory represents a number of bytes
type Memory uint

const (
	Byte     Memory = 1
	Kilobyte        = 1024 * Byte
	Megabyte        = 1024 * Kilobyte
	Gigabyte        = 1024 * Megabyte
)

type CgroupVersion string

type Metrics interface {
	*stats.Metrics | *v1.Metrics
}

const (
	V1 CgroupVersion = "v1"
	V2 CgroupVersion = "v2"
)

// String is used both by fmt.Print and by Cobra in help textv
func (e *CgroupVersion) String() string {
	return string(*e)
}

// Set must have pointer receiver so it doesn't change the value of a copy
func (e *CgroupVersion) Set(v string) error {
	switch v {
	case "v1":
		*e = CgroupVersion(v)
		return nil
	case "v2":
		*e = CgroupVersion(v)
		return nil
	default:
		return errors.New(`must be one of ["v1", "v2"]`)
	}
}

// Type is only used in help text
func (e *CgroupVersion) Type() string {
	return "cgroup version"
}
