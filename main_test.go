package main

import (
	"os"
	"testing"
)

func Test_main(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "init",
			args: []string{"kk", "init", "--modulepath=github.com/yourworkspace/CoolModuleName"},
		},
		{
			name: "add",
			args: []string{"kk", "add", "resource", "customer"},
		},
	}
	for _, tt := range tests {
		os.Args = tt.args
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}
