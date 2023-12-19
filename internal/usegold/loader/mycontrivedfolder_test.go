package loader_test

import (
	. "github.com/monopole/mdparse/internal/usegold/loader"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"os"
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

var ts = DefaultFsLoader

func TestMakeTreeItem(t *testing.T) {
	type testC struct {
		arg     string
		topName string
	}
	for n, tc := range map[string]testC{
		"t1": {
			arg:     "/home/jregan/myrepos/github.com/monopole/mdparse",
			topName: "/",
		},
		"t2": {
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

func TestMyContrivedFolderErrors(t *testing.T) {
	var (
		err error
		f   MyContrivedFolder
	)
	err = f.Initialize("/etc/passwd", ts)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not a markdown file")
	err = f.Initialize("/etc", ts)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unable to read folder")
}

func TestMyContrivedFolderHappy(t *testing.T) {
	// turnOnDebugging()
	type testC struct {
		arg string
	}
	for n, tc := range map[string]testC{
		"t1": {
			arg: "fart.md",
		},
		"t2": {
			arg: "/home/jregan/myrepos/github.com/monopole/mdparse/internal/usegold/loader/fart.md",
		},
		"t3": {
			arg: "/home/jregan/myrepos/github.com/monopole/mdparse",
		},
		"t4": {
			arg: "/home/jregan/myrepos/github.com/monopole/mdrip",
		},
		"t5": {
			arg: "/home/jregan/myrepos/github.com/monopole/mdrip/README.md",
		},
		"gh1": {
			arg: "git@github.com:monopole/mdrip.git",
		},
		"gh2": {
			arg: "git@github.com:monopole/mdrip.git/data",
		},
	} {
		t.Run(n, func(t *testing.T) {
			var f MyContrivedFolder
			err := f.Initialize(tc.arg, ts)
			assert.NoError(t, err)
			f.Accept(&VisitorDump{})
			f.Cleanup()
		})
	}
}
