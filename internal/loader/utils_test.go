package loader_test

import (
	. "github.com/monopole/mdparse/internal/loader"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

// Demonstrate the difference between
//
//	dir, base = DirBase(tc.arg)
//	dir, base = filepath.Split
func TestDirBase(t *testing.T) {
	type result struct {
		dir  string
		base string
	}
	type testC struct {
		arg string
		r0  *result // DirBase(tc.arg)
		r1  *result // filepath.Split
	}
	for n, tc := range map[string]testC{
		"t0": {
			arg: "../aaa/bbb/ccc",
			r0: &result{
				dir:  "../aaa/bbb",
				base: "ccc",
			},
			r1: &result{
				dir:  "../aaa/bbb/",
				base: "ccc",
			},
		},
		"t1": {
			arg: "/aaa/bbb/ccc",
			r0: &result{
				dir:  "/aaa/bbb",
				base: "ccc",
			},
			r1: &result{
				dir:  "/aaa/bbb/",
				base: "ccc",
			},
		},
		"t2": {
			arg: "/bbb",
			r0: &result{
				dir:  "/",
				base: "bbb",
			},
			// r1==r0
		},
		"t3": {
			arg: "bbb",
			r0: &result{
				dir:  ".",
				base: "bbb",
			},
			r1: &result{
				dir:  "",
				base: "bbb",
			},
		},
		"t4": {
			arg: "",
			r0: &result{
				dir:  ".",
				base: ".",
			},
			r1: &result{
				dir:  "",
				base: "",
			},
		},
		"t5": {
			arg: "/",
			r0: &result{
				dir:  "/",
				base: "/",
			},
			r1: &result{
				dir:  "/",
				base: "",
			},
		},
		"t6": {
			arg: "./bob/sally",
			r0: &result{
				dir:  "bob",
				base: "sally",
			},
			r1: &result{
				dir:  "./bob/",
				base: "sally",
			},
		},
		"t7": {
			arg: "./bob",
			r0: &result{
				dir:  ".",
				base: "bob",
			},
			r1: &result{
				dir:  "./",
				base: "bob",
			},
		},
		"t8": {
			arg: ".",
			r0: &result{
				dir:  ".",
				base: ".",
			},
			r1: &result{
				dir:  "",
				base: ".",
			},
		},
		"t9": {
			arg: "./",
			r0: &result{
				dir:  ".",
				base: ".",
			},
			r1: &result{
				dir:  "./",
				base: "",
			},
		},
	} {
		t.Run(n, func(t *testing.T) {
			var dir, base string
			dir, base = DirBase(tc.arg)
			assert.Equal(t, tc.r0.dir, dir)
			assert.Equal(t, tc.r0.base, base)
			dir, base = filepath.Split(tc.arg)
			if tc.r1 == nil {
				tc.r1 = tc.r0 // same result
			}
			assert.Equal(t, tc.r1.dir, dir)
			assert.Equal(t, tc.r1.base, base)
		})
	}
}
