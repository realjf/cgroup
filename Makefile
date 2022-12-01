test:
	@go test -v ./...


push:
	@git add -A && git commit -m "update" && git push origin master


build_v1:
	@rm -rf build && mkdir build
	@go build -ldflags='-s -w' -o build/cgroup_v1  ./example/v1/cgroup_v1.go

run_v1:
	@sudo ./build/cgroup_v1

build_v2:
	@rm -rf build && mkdir build
	@go build -ldflags='-s -w' -o build/cgroup_v2 ./example/v2/cgroup_v2.go

run_v2:
	@sudo ./build/cgroup_v2

.PHONY: test push build
