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

// Demo the difference between filepath.Split and FSplit.
func TestFsSplit(t *testing.T) {
	type result struct {
		d  string
		fn string
	}
	type testC struct {
		arg string
		r1  *result
		r2  *result
	}
	for n, tc := range map[string]testC{
		"t1": {
			arg: "/aaa/bbb/ccc",
			r1: &result{
				d:  "/aaa/bbb/",
				fn: "ccc",
			},
			r2: &result{
				d:  "/aaa/bbb",
				fn: "ccc",
			},
		},
		"t2": {
			arg: "/bbb",
			r1: &result{
				d:  "/",
				fn: "bbb",
			},
		},
		"t3": {
			arg: "bbb",
			r1: &result{
				d:  "",
				fn: "bbb",
			},
		},
		"t4": {
			arg: "",
			r1: &result{
				d:  "",
				fn: "",
			},
		},
		"t5": {
			arg: "/",
			r1: &result{
				d:  "/",
				fn: "",
			},
		},
		"t6": {
			arg: "./bob/sally",
			r1: &result{
				d:  "./bob/",
				fn: "sally",
			},
			r2: &result{
				d:  "bob",
				fn: "sally",
			},
		},
		"t7": {
			arg: "./bob",
			r1: &result{
				d:  "./",
				fn: "bob",
			},
			r2: &result{
				d:  "",
				fn: "bob",
			},
		},
		"t8": {
			arg: ".",
			r1: &result{
				d:  "",
				fn: ".", // odd
			},
			r2: &result{
				d:  "",
				fn: "",
			},
		},
		"t9": {
			arg: "./",
			r1: &result{
				d:  "./",
				fn: "",
			},
			r2: &result{
				d:  "",
				fn: "",
			},
		},
	} {
		t.Run(n, func(t *testing.T) {
			d, fn := filepath.Split(tc.arg)
			assert.Equal(t, tc.r1.d, d)
			assert.Equal(t, tc.r1.fn, fn)
			d, fn = FSplit(tc.arg)
			if tc.r2 == nil {
				tc.r2 = tc.r1 // same result
			}
			assert.Equal(t, tc.r2.d, d)
			assert.Equal(t, tc.r2.fn, fn)
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
