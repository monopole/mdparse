package loader_test

import (
	. "github.com/monopole/mdparse/internal/usegold/loader"
	"github.com/stretchr/testify/assert"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// Demo the difference between
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
		r0  *result
		r1  *result
		r2  *result
	}
	for n, tc := range map[string]testC{
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
			r2: &result{
				dir:  "/aaa/bbb", // no trailing slash
				base: "ccc",
			},
		},
		"t2": {
			arg: "/bbb",
			r0: &result{
				dir:  "/",
				base: "bbb",
			},
			// r2==r1==r0
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
			// r2 == r1
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
			// r2 == r1
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
			// r2 == r1
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
			r2: &result{
				dir:  "bob", // no dot or trailing slash
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
			r2: &result{
				dir:  "", // no dot or trailing slash
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
				base: ".", // why does it do this?
			},
			r2: &result{
				dir:  "",
				base: "", // no dot
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

type mockFileInfo struct {
	name string
	mode fs.FileMode
}

func (m *mockFileInfo) Name() string       { return m.name }
func (m *mockFileInfo) Size() int64        { return 0 }
func (m *mockFileInfo) Mode() fs.FileMode  { return m.mode }
func (m *mockFileInfo) ModTime() time.Time { panic("didn't think ModTime was needed") }
func (m *mockFileInfo) IsDir() bool        { panic("didn't think IsDir was needed") }
func (m *mockFileInfo) Sys() any           { panic("didn't think Sys was needed") }

var _ os.FileInfo = &mockFileInfo{}

type tCase struct {
	fi  *mockFileInfo
	err error
}

func TestIsMarkDownFile(t *testing.T) {
	for n, tc := range map[string]tCase{
		"t1": {
			fi: &mockFileInfo{
				name: "aDirectory.md",
				mode: fs.ModeDir,
			},
			err: NotMarkDownErr,
		},
		"t2": {
			fi: &mockFileInfo{
				name: "notMarkdown",
			},
			err: NotMarkDownErr,
		},
		"t3": {
			fi: &mockFileInfo{
				name: "aFileButIrregular.md",
				mode: fs.ModeIrregular,
			},
			err: NotMarkDownErr,
		},
		"t4": {
			fi: &mockFileInfo{
				name: "aFile.md",
			},
		},
	} {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, tc.err, IsMarkDownFile(tc.fi))
		})
	}
}

func TestIsNotADotDir(t *testing.T) {
	for n, tc := range map[string]tCase{
		"t1": {
			fi: &mockFileInfo{
				name: ".git",
			},
			err: IsADotDirErr,
		},
		"t2": {
			fi: &mockFileInfo{
				name: "./",
			},
		},
		"t3": {
			fi: &mockFileInfo{
				name: "..",
			},
		},
	} {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, tc.err, IsNotADotDir(tc.fi))
		})
	}
}
