package loader

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMyFSplit(t *testing.T) {
	type testC struct {
		arg string
		d   string
		fn  string
	}
	for n, tc := range map[string]testC{
		"t1": {
			arg: "/home/aaa/bbb",
			d:   "/home/aaa",
			fn:  "bbb",
		},
		"t2": {
			arg: "/bbb",
			d:   "",
			fn:  "bbb",
		},
		"t3": {
			arg: "bbb",
			d:   "",
			fn:  "bbb",
		},
		"t4": {
			arg: "",
			d:   "",
			fn:  "",
		},
		"t5": {
			arg: "/",
			d:   "",
			fn:  "",
		},
	} {
		t.Run(n, func(t *testing.T) {
			d, fn := fSplit(tc.arg)
			assert.Equal(t, tc.d, d)
			assert.Equal(t, tc.fn, fn)
		})
	}
}

func TestMyIsOrderFile(t *testing.T) {
	type args struct {
		n string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, MyIsOrderFile(tt.args.n), "MyIsOrderFile(%v)", tt.args.n)
		})
	}
}

func Test_isAnAllowedFile(t *testing.T) {
	type args struct {
		info os.FileInfo
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, isAnAllowedFile(tt.args.info), "isAnAllowedFile(%v)", tt.args.info)
		})
	}
}

func Test_isAnAllowedFolder(t *testing.T) {
	type args struct {
		info os.FileInfo
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, isAnAllowedFolder(tt.args.info), "isAnAllowedFolder(%v)", tt.args.info)
		})
	}
}
