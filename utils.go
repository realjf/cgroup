package cgroup

import "github.com/containerd/cgroups"

func IsCgroupV2() bool {
	return cgroups.Mode() == cgroups.Unified
}
