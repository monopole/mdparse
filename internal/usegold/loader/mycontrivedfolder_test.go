package loader

import (
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

func TestMyContrivedFolderDebug(t *testing.T) {
	var (
		err error
		f   MyContrivedFolder
	)
	err = f.Initialize(
		[]string{"/etc/passwd"}, IsAnAllowedFile, IsAnAllowedFolder)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not a markdown file")

	err = f.Initialize([]string{
		"fart.md",
		"/etc",
		"/home/jregan/myrepos/github.com/monopole/mdparse/internal/usegold/loader/fart.md",
		"/home/jregan/myrepos/github.com/monopole/mdrip",
		"/home/jregan/myrepos/github.com/monopole/mdrip/README.md",
	}, IsAnAllowedFile, IsAnAllowedFolder)
	assert.NoError(t, err)
	assert.NotNil(t, f)

	f.Accept(&VisitorDump{})
}
