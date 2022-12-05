package cgroup

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/realjf/utils"
	"github.com/stretchr/testify/assert"
)

func TestCgroup(t *testing.T) {
	cases := map[string]struct {
		name    string
		slice   string
		group   string
		version CgroupVersion
		cpu     Percent
		mem     Memory
	}{
		"v1": {
			name:    "test",
			version: V1,
			cpu:     80,
			mem:     8 * Megabyte,
		},
		"v2": {
			version: V2,
			slice:   "/",
			group:   "hello.slice",
			cpu:     80,
			mem:     8 * Megabyte,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			cg, err := NewCgroup(tc.version, WithName(tc.name), WithSlice(tc.slice), WithGroup(tc.group), WithCPULimit(tc.cpu), WithMemoryLimit(tc.mem))
			assert.NoError(t, err)
			if assert.NotNil(t, cg) {
				err = cg.Create()
				if assert.NoError(t, err) {
					err = cg.Close()
					assert.NoError(t, err)
				}

			}

		})
	}
}

func TestV2(t *testing.T) {
	cmd := utils.NewCmd()
	defer cmd.Close()
	args := []string{"-c", "$(echo -1000 > /proc/" + strconv.Itoa(os.Getpid()) + "/oom_score_adj)"}
	_, err := cmd.RunCommand("/bin/bash", args...)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	} else {
		fmt.Println("done")
	}
}
