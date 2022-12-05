# cgroup
a go library to control cgroup


### Usage

**`import`**
```sh
go get github.com/realjf/cgroup
```

**`cgroup v1 example`**
```go
import (
 "fmt"
 "os"
 "sync"

 "github.com/realjf/utils"
 "github.com/sirupsen/logrus"

 "github.com/realjf/cgroup"
)

type Limiter struct {
 cg cgroup.ICgroup
 wg sync.WaitGroup
}

func main() {
 limiter := &Limiter{
  wg: sync.WaitGroup{},
 }

 var err error

 limiter.cg, err = cgroup.NewCgroup(cgroup.V1, cgroup.WithName("test"))
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
 limiter.cg.SetOptions(cgroup.WithCPULimit(80))                      // cpu usage limit 80%
 limiter.cg.SetOptions(cgroup.WithMemoryLimit(20 * cgroup.Megabyte)) // memory limit 8MB
 limiter.cg.SetOptions(cgroup.WithDisableOOMKiller())                // disable oom killer

 err = limiter.cg.Create()
 if err != nil {
  logrus.Println(err.Error())
  return
 }
 // limit by pid
 pid := os.Getpid()
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
 limiter.wg.Add(1)
 go limiter.command()

 limiter.wg.Wait()

 logrus.Println("done")
}

func (l *Limiter) command() {
 defer func() {
  l.wg.Done()
 }()
 cmd := utils.NewCmd().SetDebug(true)
 defer cmd.Close()

 var err error
 err = cmd.SetUsername(os.Getenv("SUDO_USER"))
 if err != nil {
  fmt.Println(err.Error())
  return
 }
 cmd.SetNoSetGroups(true)

 args := []string{"--cpu", "1", "--vm", "1", "--vm-bytes", "20M", "--timeout", "10s", "--vm-keep"}
 _, err = cmd.Command("stress", args...)
 if err != nil {
  fmt.Println(err.Error())
  return
 }

 out, err := cmd.Run()
 if err != nil {
  fmt.Println(err.Error())
  return
 }
 fmt.Printf("%s\n", out)
}



```

**`cgroup v2 example`**

```go
import (
 "fmt"
 "os"
 "sync"

 "github.com/realjf/utils"
 "github.com/sirupsen/logrus"

 "github.com/realjf/cgroup"
)

type Limiter struct {
 cg cgroup.ICgroup
 wg sync.WaitGroup
}

func main() {
 limiter := &Limiter{
  wg: sync.WaitGroup{},
 }

 var err error

 limiter.cg, err = cgroup.NewCgroup(cgroup.V2, cgroup.WithSlice("/"), cgroup.WithGroup("mycgroup.slice"))
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
 limiter.cg.SetOptions(cgroup.WithCPULimit(80))                      // cpu usage limit 80%
 limiter.cg.SetOptions(cgroup.WithMemoryLimit(20 * cgroup.Megabyte)) // memory limit 8MB
 limiter.cg.SetOptions(cgroup.WithDisableOOMKiller())                // disable oom killer

 err = limiter.cg.Create()
 if err != nil {
  logrus.Println(err.Error())
  return
 }
 // limit by pid
 pid := os.Getpid()
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
 limiter.wg.Add(1)
 go limiter.command()

 limiter.wg.Wait()

 logrus.Println("done")
}

func (l *Limiter) command() {
 defer func() {
  l.wg.Done()
 }()
 cmd := utils.NewCmd().SetDebug(true)
 defer cmd.Close()

 var err error
 err = cmd.SetUsername(os.Getenv("SUDO_USER"))
 if err != nil {
  fmt.Println(err.Error())
  return
 }
 cmd.SetNoSetGroups(true)

 args := []string{"--cpu", "1", "--vm", "1", "--vm-bytes", "20M", "--timeout", "10s", "--vm-keep"}
 _, err = cmd.Command("stress", args...)
 if err != nil {
  fmt.Println(err.Error())
  return
 }

 out, err := cmd.Run()
 if err != nil {
  fmt.Println(err.Error())
  return
 }
 fmt.Printf("%s\n", out)
}
```
