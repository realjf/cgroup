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

```

**`cgroup v2 example`**

```go
 cg, err := NewCgroup(V2, WithSlice("/"), WithGroup("hello.slice"))
 if err != nil {
  log.Println(err.Error())
  return
 }
 // limit
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

 args := []string{"--cpu", "4", "--vm", "2", "--vm-bytes", "120M", "--timeout", "20s"}
 out, err := utils.NewCmd().RunCommand("stress", args...)
 if err != nil {
  log.Println(err.Error())
  return
 }
 log.Printf("%s", out)

 err = cg.Close()
 if err != nil {
  log.Println(err.Error())
  return
 }
 log.Println("done")
```
