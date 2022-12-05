package cgroup

import (
	"fmt"
	"sync"
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
	Stats() (any, error)
}

type cgroupImpl struct {
	version CgroupVersion
	cg      ICgroup
	ch      chan bool
	lock    sync.Mutex
}

func NewCgroup(version CgroupVersion, options ...Option) (ICgroup, error) {
	cg := &cgroupImpl{
		version: version,
		ch:      make(chan bool),
		lock:    sync.Mutex{},
	}

	switch version {
	case V1:
		cg.cg = newCgroupImplV1()
	case V2:
		cg.cg = newCgroupImplV2()
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
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.cg.Create()
}

func (c *cgroupImpl) SetOptions(options ...Option) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.cg.SetOptions(options...)
}

func (c *cgroupImpl) Instance() any {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.cg.Instance()
}

func (c *cgroupImpl) Close() error {
	close(c.ch)
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.cg.Close()
}

func (c *cgroupImpl) Load() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.cg.Load()
}

func (c *cgroupImpl) LimitPid(pid int) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.cg.LimitPid(pid)
}

func (c *cgroupImpl) GetLimitPids() ([]uint64, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.cg.GetLimitPids()
}

func (c *cgroupImpl) Stats() (any, error) {
	return c.cg.Stats()
}
