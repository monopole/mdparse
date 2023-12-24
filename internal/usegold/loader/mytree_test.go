package loader_test

import (
	"fmt"
	. "github.com/monopole/mdparse/internal/usegold/loader"
	"github.com/spf13/afero"
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
			errMsg: "not a simple markdown file",
		},
		"t2": {
			arg:    "/etc",
			errMsg: "unable to read folder",
		},
	} {
		t.Run(n, func(t *testing.T) {
			_, err := MakeTreeItem(NewFsLoader(afero.NewOsFs()), tc.arg)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.errMsg)
		})
	}
}

func TestMakeTreeItemHappy(t *testing.T) {
	{
		var cwd string
		var err error
		cwd, err = os.Getwd()
		if err != nil {
			return
		}
		fmt.Println("cwd of test =", cwd)
	}
	type testC struct {
		arg     string
		topName string
	}
	for n, tc := range map[string]testC{
		"t1": {
			arg:     "fart.md",
			topName: ".",
		},
		"t2": {
			arg:     "/home/jregan/myrepos/github.com/monopole/mdparse/internal/usegold/loader/fart.md",
			topName: "/home/jregan/myrepos/github.com/monopole/mdparse/internal/usegold/loader",
		},
		"t3": {
			arg:     "/home/jregan/myrepos/github.com/monopole/mdparse",
			topName: "/home/jregan/myrepos/github.com/monopole",
		},
		"t4": {
			arg:     "/home/jregan/myrepos/github.com/monopole/mdrip",
			topName: "/home/jregan/myrepos/github.com/monopole",
		},
		"t5": {
			arg:     "/home/jregan/myrepos/github.com/monopole/mdrip/README.md",
			topName: "/home/jregan/myrepos/github.com/monopole/mdrip",
		},
		"t7": {
			arg:     ".",
			topName: ".",
		},
	} {
		t.Run(n, func(t *testing.T) {
			fsl := NewFsLoader(afero.NewOsFs())
			f, err := MakeTreeItem(fsl, tc.arg)
			if err == nil {
				fmt.Println("no error!")
			} else {
				fmt.Println("err: ", err.Error())
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.topName, f.Name())
			f.Accept(NewVisitorDump(fsl))
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
			fsl := NewFsLoader(afero.NewOsFs())
			f, err := MakeTreeItem(fsl, tc.arg)
			assert.NoError(t, err)
			f.Accept(NewVisitorDump(fsl))
			assert.Equal(t, tc.topName, f.Name())
			// f.Cleanup()
		})
	}
}
