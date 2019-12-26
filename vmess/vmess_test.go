package vmess

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestNewQuanVmess(t *testing.T) {

	Must := func(v interface{}, err error) interface{} {
		if err != nil {
			t.Error(err)
		}

		return v
	}

	tests := []struct {
		name    string
		args    string
		want    *VmessLink
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"1",
			"vmess://",
			Must(NewVnVmess("")).(*VmessLink),
			false,
		},
		{
			"2",
			"",
			Must(NewVnVmess("")).(*VmessLink),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewQuanVmess(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewQuanVmess() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.want.Aid = "0"
			tt.want.OrigLink = got.OrigLink

			if d := cmp.Diff(got, tt.want); d != "" {
				t.Error(d)
			}
		})
	}
}
