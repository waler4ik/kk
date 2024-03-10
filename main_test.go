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
			args: []string{"kk", "init", "--modulename=github.com/yourworkspace/CoolModuleName"},
		},
	}
	for _, tt := range tests {
		os.Args = tt.args
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}
