package loader_test

import (
	. "github.com/monopole/mdparse/internal/usegold/loader"
	"log/slog"
	"os"

	"github.com/stretchr/testify/assert"
	"testing"
)

func turnOnDebugging() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey || a.Key == slog.LevelKey {
				a.Value = slog.StringValue("")
			}
			return a
		},
	})))
}

func TestMakeTreeItemErrors(t *testing.T) {
	type testC struct {
		arg    string
		errMsg string
	}
	for n, tc := range map[string]testC{
		"t1": {
			arg:    "/etc/passwd",
			errMsg: "not a markdown file",
		},
		"t2": {
			arg:    "/etc",
			errMsg: "unable to read folder",
		},
	} {
		t.Run(n, func(t *testing.T) {
			_, err := MakeTreeItem(DefaultFsLoader, tc.arg)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.errMsg)
		})
	}
}

func TestMakeTreeItemHappy(t *testing.T) {
	type testC struct {
		arg     string
		topName string
	}
	for n, tc := range map[string]testC{
		"t1": {
			arg:     "fart.md",
			topName: "/",
		},
		"t2": {
			arg:     "/home/jregan/myrepos/github.com/monopole/mdparse/internal/usegold/loader/fart.md",
			topName: "/",
		},
		"t3": {
			arg:     "/home/jregan/myrepos/github.com/monopole/mdparse",
			topName: "/",
		},
		"t4": {
			arg:     "/home/jregan/myrepos/github.com/monopole/mdrip",
			topName: "/",
		},
		"t5": {
			arg:     "/home/jregan/myrepos/github.com/monopole/mdrip/README.md",
			topName: "/",
		},
		"t6": {
			arg:     "/home/jregan/myrepos/github.com/monopole/mdparse",
			topName: "/",
		},
		"t7": {
			arg:     ".",
			topName: "/home/jregan/myrepos/github.com/monopole/mdparse/internal/usegold/loader",
		},
	} {
		t.Run(n, func(t *testing.T) {
			f, err := MakeTreeItem(DefaultFsLoader, tc.arg)
			assert.NoError(t, err)
			f.Accept(&VisitorDump{})
			assert.Equal(t, tc.topName, f.Name())
		})
	}
}

func TestMakeTreeItemRepo(t *testing.T) {
	// turnOnDebugging()
	type testC struct {
		arg     string
		topName string
	}
	for n, tc := range map[string]testC{
		"gh1": {
			arg:     "git@github.com:monopole/mdrip.git",
			topName: "/",
		},
		"gh2": {
			arg:     "git@github.com:monopole/mdrip.git/data",
			topName: "/",
		},
	} {
		t.Run(n, func(t *testing.T) {
			f, err := MakeTreeItem(DefaultFsLoader, tc.arg)
			assert.NoError(t, err)
			f.Accept(&VisitorDump{})
			assert.Equal(t, tc.topName, f.Name())
			// f.Cleanup()
		})
	}
}
