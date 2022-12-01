package cgroup

import (
	"fmt"
)

type ICgroup interface {
	Create() error
	Version() CgroupVersion
	SetOptions(options ...Option)
	Instance() any
	Close() error
	Load() error
	LimitPid(pid int) error
	GetLimitPids() ([]uint64, error)
	WaitForEvents()
}

type cgroupImpl struct {
	version CgroupVersion
	cg      ICgroup
	ch      chan bool
}

func NewCgroup(version CgroupVersion, options ...Option) (ICgroup, error) {
	cg := &cgroupImpl{
		version: version,
		ch:      make(chan bool),
	}

	switch version {
	case V1:
		cg.cg = NewCgroupImplV1()
	case V2:
		cg.cg = NewCgroupImplV2()
	default:
		return nil, fmt.Errorf("unsupported cgroup version")
	}

	cg.SetOptions(options...)

	return cg, nil
}

func (c *cgroupImpl) Version() CgroupVersion {
	return c.version
}

func (c *cgroupImpl) Create() error {
	return c.cg.Create()
}

func (c *cgroupImpl) SetOptions(options ...Option) {
	c.cg.SetOptions(options...)
}

func (c *cgroupImpl) Instance() any {
	return c.cg.Instance()
}

func (c *cgroupImpl) Close() error {
	close(c.ch)
	return c.cg.Close()
}

func (c *cgroupImpl) Load() error {
	return c.cg.Load()
}

func (c *cgroupImpl) LimitPid(pid int) error {
	return c.cg.LimitPid(pid)
}

func (c *cgroupImpl) WaitForEvents() {
	for {
		select {
		case <-c.ch:
			return
		default:
			c.cg.WaitForEvents()
		}
	}
}

func (c *cgroupImpl) GetLimitPids() ([]uint64, error) {
	return c.cg.GetLimitPids()
}
