package cgroup

import (
	"errors"
	"strconv"

	"github.com/containerd/cgroups/v3/cgroup2"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/sirupsen/logrus"
)

type CgV2 struct {
}

type CgroupV2 interface {
	ICgroup
}

type cgroupImplV2 struct {
	cg      *cgroup2.Manager
	slice   string
	group   string
	res     *cgroup2.Resources
	version CgroupVersion
	oomkill bool // whether to use oomkill
	cevent  <-chan cgroup2.Event
	cerr    <-chan error
}

func NewCgroupImplV2() *cgroupImplV2 {
	return &cgroupImplV2{
		version: V2,
		res: &cgroup2.Resources{
			CPU:     &cgroup2.CPU{},
			Memory:  &cgroup2.Memory{},
			Pids:    &cgroup2.Pids{},
			IO:      &cgroup2.IO{},
			RDMA:    &cgroup2.RDMA{},
			HugeTlb: &cgroup2.HugeTlb{},
			Devices: make([]specs.LinuxDeviceCgroup, 0),
		},
		cevent:  make(<-chan cgroup2.Event),
		cerr:    make(<-chan error, 1),
		oomkill: false,
	}
}

func (c *cgroupImplV2) Version() CgroupVersion {
	return c.version
}

func (c *cgroupImplV2) SetOptions(options ...Option) {
	for _, opt := range options {
		opt(c)
	}
}

func (c *cgroupImplV2) Close() error {
	return c.cg.DeleteSystemd()
}

func (c *cgroupImplV2) Load() error {
	var err error
	c.cg, err = cgroup2.LoadSystemd(c.slice, c.group)
	return err
}

func (c *cgroupImplV2) Instance() any {
	return c
}

func (c *cgroupImplV2) Create() error {
	if c.slice == "" {
		return errors.New("slice is empty")
	}

	if c.group == "" {
		return errors.New("group is empty")
	}

	// dummy PID of -1 is used for creating a "general slice" to be used as a parent cgroup.
	// see https://github.com/containerd/cgroups/blob/1df78138f1e1e6ee593db155c6b369466f577651/v2/manager.go#L732-L735
	// for example: slice="/" group="hello.slice"
	var err error
	c.cg, err = cgroup2.NewSystemd(c.slice, c.group, -1, c.res)
	return err
}

func (c *cgroupImplV2) LimitPid(pid int) error {
	pid_u64, err := strconv.ParseUint(strconv.Itoa(pid), 10, 64)
	if err != nil {
		return err
	}
	return c.cg.AddProc(pid_u64)
}

func (c *cgroupImplV2) GetLimitPids() ([]uint64, error) {
	return c.cg.Procs(true)
}

func (c *cgroupImplV2) disableOOMKiller() {
	c.oomkill = true
}

func (c *cgroupImplV2) handleOOMKillEvent(cgroup2.Event) {
	logrus.Info("oom killed")
}

func (c *cgroupImplV2) WaitForEvents() {
	c.cevent, c.cerr = c.cg.EventChan()
	err := <-c.cerr
	logrus.Infof("e:%#v", err)
	if err != nil {
		logrus.Error(err)
		return
	}
	event := <-c.cevent
	logrus.Infof("e:%#v", event)
	if c.oomkill {
		if event.OOM > 0 {
			c.handleOOMKillEvent(event)
		}
		if event.OOMKill > 0 {
			c.handleOOMKillEvent(event)
		}
	}
}

func (c *cgroupImplV2) Stats() (any, error) {
	// return *stats.Metrics
	return c.cg.Stat()
}
