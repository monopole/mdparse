package loader

import (
	"github.com/stretchr/testify/assert"
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
