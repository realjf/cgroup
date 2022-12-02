package main

import (
	"os/user"
	"syscall"

	. "github.com/realjf/cgroup"
	"github.com/realjf/utils"
	"github.com/sirupsen/logrus"
)

func main() {
	type Limiter struct {
		cg  ICgroup
		cmd *utils.Command
	}

	limiter := &Limiter{
		cmd: utils.NewCmd(),
	}
	defer limiter.cmd.Close()

	var err error
	user, err := user.Current()
	if err != nil {
		logrus.Println(err.Error())
		return
	}
	limiter.cmd.SetUser(user)
	attr := syscall.SysProcAttr{
		// Cloneflags:                 syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER | syscall.CLONE_NEWNET,
		// GidMappingsEnableSetgroups: true,
		Setpgid: true,
		// UidMappings: []syscall.SysProcIDMap{
		// 	{
		// 		ContainerID: 0,
		// 		HostID:      0,
		// 		Size:        1,
		// 	},
		// },
		// GidMappings: []syscall.SysProcIDMap{
		// 	{
		// 		ContainerID: 0,
		// 		HostID:      0,
		// 		Size:        1,
		// 	},
		// },
		Pgid:       0,
		Credential: &syscall.Credential{},
	}
	limiter.cmd.SetSysProcAttr(attr)

	limiter.cg, err = NewCgroup(V2, WithSlice("/"), WithGroup("mycgroup.slice"))
	if err != nil {
		logrus.Println(err.Error())
		return
	}
	defer func() {
		err = limiter.cg.Close()
		if err != nil {
			logrus.Println(err.Error())
			return
		}
	}()
	// limit
	limiter.cg.SetOptions(WithCPULimit(80))              // cpu usage limit 80%
	limiter.cg.SetOptions(WithMemoryLimit(8 * Megabyte)) // memory limit 8MB
	limiter.cg.SetOptions(WithDisableOOMKiller())        // disable oom killer

	err = limiter.cg.Create()
	if err != nil {
		logrus.Println(err.Error())
		return
	}

	args := []string{"--cpu", "1", "--vm", "1", "--vm-bytes", "20M", "--timeout", "20s", "--vm-keep"}
	pid, err := limiter.cmd.Command("stress", args...)
	if err != nil {
		logrus.Println(err.Error())
		return
	}
	// limit by pid
	logrus.Printf("limit pid: %d\n", pid)
	err = limiter.cg.LimitPid(pid)
	if err != nil {
		logrus.Println(err.Error())
		return
	}
	wpids, err := limiter.cg.GetLimitPids()
	if err != nil {
		logrus.Println(err.Error())
		return
	}
	logrus.Printf("limit pid now: %v\n", wpids)

	out, err := limiter.cmd.Run()
	if err != nil {
		errout, _ := limiter.cmd.GetStderrOutput()
		logrus.Printf("run cmd stderr:%s\n", errout)
		logrus.Println(err.Error())
		return
	}
	logrus.Printf("%s", out)

	logrus.Println("done")
}
