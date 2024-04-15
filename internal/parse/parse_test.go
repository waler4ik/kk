package parse_test

import (
	"testing"

	"github.com/waler4ik/kk/internal/parse"
)

func TestRouterType(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "chi router type",
			args: args{
				path: "./data/chi_test.go.tmpl",
			},
			want:    "chi",
			wantErr: false,
		},
		{
			name: "gin router type",
			args: args{
				path: "./data/gin_test.go.tmpl",
			},
			want:    "gin",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parse.RouterType(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("RouterType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RouterType() = %v, want %v", got, tt.want)
			}
		})
	}
}
