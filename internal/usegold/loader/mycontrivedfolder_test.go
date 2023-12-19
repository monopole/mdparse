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

func TestMyContrivedFolderErrors(t *testing.T) {
	var (
		err error
		f   MyContrivedFolder
	)
	err = f.Initialize([]string{"/etc/passwd"}, ts)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not a markdown file")
	err = f.Initialize([]string{"/etc"}, ts)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unable to read folder")
}

func TestMyContrivedFolderHappy(t *testing.T) {
	var (
		err error
		f   MyContrivedFolder
	)
	// turnOnDebugging()
	err = f.Initialize([]string{
		"fart.md",
		"/home/jregan/myrepos/github.com/monopole/mdparse/internal/usegold/loader/fart.md",
		"/home/jregan/myrepos/github.com/monopole/mdparse",
		"/home/jregan/myrepos/github.com/monopole/mdrip",
		"/home/jregan/myrepos/github.com/monopole/mdrip/README.md",
	}, ts)
	assert.NoError(t, err)
	f.Accept(&VisitorDump{})
	f.Cleanup()
}

func xxx_disabled_TestMyContrivedFolderGit(t *testing.T) {
	var (
		err error
		f   MyContrivedFolder
	)
	// turnOnDebugging()
	err = f.Initialize([]string{
		"fart.md",
		//"git@github.com:monopole/mdrip.git",
		"git@github.com:monopole/mdrip.git/data",
	}, ts)
	assert.NoError(t, err)
	f.Accept(&VisitorDump{})
	f.Cleanup()
}
