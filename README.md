# cgroup
a go library to control cgroup


### Usage

**`import`**
```sh
go get github.com/realjf/cgroup
```

**`cgroup v1 example`**
```go
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

 limiter.cg, err = NewCgroup(V1, WithName("test"))
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


```

**`cgroup v2 example`**

```go
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

```
