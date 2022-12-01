package cgroup

import (
	"github.com/containerd/cgroups/v3/cgroup2"
	"github.com/opencontainers/runtime-spec/specs-go"
)

type Option func(cgroup ICgroup)

// for cgroup v1
func WithName(name string) Option {
	return func(cgroup ICgroup) {
		if c, ok := cgroup.Instance().(*cgroupImplV1); ok {
			c.name = name
		}
	}
}

// for cgroup v2
func WithSlice(slice string) Option {
	return func(cgroup ICgroup) {
		if c, ok := cgroup.Instance().(*cgroupImplV2); ok {
			c.slice = slice
		}
	}
}

// for cgroup v2
func WithGroup(group string) Option {
	return func(cgroup ICgroup) {
		if c, ok := cgroup.Instance().(*cgroupImplV2); ok {
			c.group = group
		}
	}
}

// for cgroup v1,v2
func WithCPULimit(cpuLimit Percent) Option {
	return func(cgroup ICgroup) {
		if c, ok := cgroup.Instance().(*cgroupImplV1); ok {
			if c.res.CPU == nil {
				c.res.CPU = &specs.LinuxCPU{}
			}
			if c.res.CPU.Period == nil {
				c.res.CPU.Period = new(uint64)
				*c.res.CPU.Period = 100000
			}
			c.res.CPU.Quota = new(int64)
			*c.res.CPU.Quota = int64(*c.res.CPU.Period * uint64(cpuLimit) / 100)
		} else if c, ok := cgroup.Instance().(*cgroupImplV2); ok {
			var period *uint64 = new(uint64)
			*period = 100000
			var quota *int64 = new(int64)
			*quota = int64(*period * uint64(cpuLimit) / 100)
			c.res.CPU.Max = cgroup2.NewCPUMax(quota, period)
		}
	}
}

// for cgroup v1,v2
func WithMemoryLimit(memory Memory) Option {
	return func(cgroup ICgroup) {
		if c, ok := cgroup.Instance().(*cgroupImplV1); ok {
			if c.res.Memory == nil {
				c.res.Memory = &specs.LinuxMemory{}
			}
			c.res.Memory.Limit = new(int64)
			*c.res.Memory.Limit = int64(memory)
		} else if c, ok := cgroup.Instance().(*cgroupImplV2); ok {
			var mem *int64 = new(int64)
			*mem = int64(memory)
			c.res.Memory.Max = mem
		}
	}
}

// for cgroup v1,v2
func WithDisableOOMKiller() Option {
	return func(cgroup ICgroup) {
		if c, ok := cgroup.Instance().(*cgroupImplV1); ok {
			c.disableOOMKiller()
		} else if c, ok := cgroup.Instance().(*cgroupImplV2); ok {
			c.disableOOMKiller()
		}
	}
}
