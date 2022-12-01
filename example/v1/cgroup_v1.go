package main

import (
	"log"
	"os"

	. "github.com/realjf/cgroup"
)

func main() {
	cg, err := NewCgroup(V1, WithName("test"))
	if err != nil {
		log.Println(err.Error())
		return
	}

	cg.SetOptions(WithCPULimit(80))              // cpu usage limit 80%
	cg.SetOptions(WithMemoryLimit(8 * Megabyte)) // memory limit 8MB

	err = cg.Create()
	if err != nil {
		log.Println(err.Error())
		return
	}

	// limit by pid
	err = cg.LimitPid(os.Getpid())
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = cg.Close()
	if err != nil {
		log.Println(err.Error())
		return
	}
}
