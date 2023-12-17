package loader_test

import (
	. "github.com/monopole/mdparse/internal/usegold/loader"
	"github.com/stretchr/testify/assert"
	"io/fs"
	"os"
	"testing"
	"time"
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
			d, fn := FSplit(tc.arg)
			assert.Equal(t, tc.d, d)
			assert.Equal(t, tc.fn, fn)
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
	fi      *mockFileInfo
	allowed bool
}

func TestIsAnAllowedFile(t *testing.T) {
	for n, tc := range map[string]tCase{
		"t1": {
			fi: &mockFileInfo{
				name: "aDirectory.md",
				mode: fs.ModeDir,
			},
		},
		"t2": {
			fi: &mockFileInfo{
				name: "notMarkdown",
			},
		},
		"t3": {
			fi: &mockFileInfo{
				name: "aFile.md",
			},
			allowed: true,
		},
	} {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, tc.allowed, IsAnAllowedFile(tc.fi))
		})
	}
}

func TestIsAnAllowedFolder(t *testing.T) {
	for n, tc := range map[string]tCase{
		"t1": {
			fi: &mockFileInfo{
				name: ".git",
			},
		},
		"t2": {
			fi: &mockFileInfo{
				name: "./",
			},
			allowed: true,
		},
		"t3": {
			fi: &mockFileInfo{
				name: "..",
			},
			allowed: true,
		},
	} {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, tc.allowed, IsAnAllowedFolder(tc.fi))
		})
	}
}
