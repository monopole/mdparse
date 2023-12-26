package loader_test

import (
	. "github.com/monopole/mdparse/internal/usegold/loader"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

// A test to demonstrate the difference between
//
//	dir, base = filepath.Dir(tc.arg), filepath.Base(tc.arg)
//	dir, base = filepath.Split
//	dir, base = FSplit
func TestFsSplit(t *testing.T) {
	type result struct {
		dir  string
		base string
	}
	type testC struct {
		arg string
		r0  *result // filepath.Dir(tc.arg), filepath.Base(tc.arg)
		r1  *result // filepath.Split
		r2  *result // FSplit
	}
	for n, tc := range map[string]testC{
		"t0": { // good
			arg: "../aaa/bbb/ccc",
			r0: &result{
				dir:  "../aaa/bbb",
				base: "ccc",
			},
			r1: &result{
				dir:  "../aaa/bbb/",
				base: "ccc",
			},
			r2: &result{
				dir:  "../aaa/bbb",
				base: "ccc",
			},
		},
		"t1": { // good
			arg: "/aaa/bbb/ccc",
			r0: &result{
				dir:  "/aaa/bbb",
				base: "ccc",
			},
			r1: &result{
				dir:  "/aaa/bbb/",
				base: "ccc",
			},
			r2: &result{
				dir:  "/aaa/bbb", // no trailing slash
				base: "ccc",
			},
		},
		"t2": { // good
			arg: "/bbb",
			r0: &result{
				dir:  "/",
				base: "bbb",
			},
			// r2==r1==r0
		},
		"t3": { // good
			arg: "bbb",
			r0: &result{
				dir:  ".",
				base: "bbb",
			},
			r1: &result{
				dir:  "",
				base: "bbb",
			},
			// r2 == r1
		},
		"t4": { // good but we need to know it
			arg: "",
			r0: &result{
				dir:  ".",
				base: ".",
			},
			r1: &result{
				dir:  "",
				base: "",
			},
			// r2 == r1
		},
		"t5": { // good but we need to know it
			arg: "/",
			r0: &result{
				dir:  "/",
				base: "/",
			},
			r1: &result{
				dir:  "/",
				base: "",
			},
			// r2 == r1
		},
		"t6": { // good
			arg: "./bob/sally",
			r0: &result{
				dir:  "bob",
				base: "sally",
			},
			r1: &result{
				dir:  "./bob/",
				base: "sally",
			},
			r2: &result{
				dir:  "bob", // no dot or trailing slash
				base: "sally",
			},
		},
		"t7": { // good
			arg: "./bob",
			r0: &result{
				dir:  ".",
				base: "bob",
			},
			r1: &result{
				dir:  "./",
				base: "bob",
			},
			r2: &result{
				dir:  "", // no dot or trailing slash
				base: "bob",
			},
		},
		"t8": { // good but we need to know it
			arg: ".",
			r0: &result{
				dir:  ".",
				base: ".",
			},
			r1: &result{
				dir:  "",
				base: ".", // why does it do this?
			},
			r2: &result{
				dir:  "",
				base: "", // no dot
			},
		},
		"t9": { // good but we need to know it
			arg: "./",
			r0: &result{
				dir:  ".",
				base: ".",
			},
			r1: &result{
				dir:  "./",
				base: "",
			},
			r2: &result{
				dir:  "", // no dot, no trailing slash
				base: "",
			},
		},
	} {
		t.Run(n, func(t *testing.T) {
			var dir, base string
			dir, base = filepath.Dir(tc.arg), filepath.Base(tc.arg)
			assert.Equal(t, tc.r0.dir, dir)
			assert.Equal(t, tc.r0.base, base)
			dir, base = filepath.Split(tc.arg)
			if tc.r1 == nil {
				tc.r1 = tc.r0 // same result
			}
			assert.Equal(t, tc.r1.dir, dir)
			assert.Equal(t, tc.r1.base, base)
			dir, base = FSplit(tc.arg)
			if tc.r2 == nil {
				tc.r2 = tc.r1 // same result
			}
			assert.Equal(t, tc.r2.dir, dir)
			assert.Equal(t, tc.r2.base, base)

		})
	}
}
