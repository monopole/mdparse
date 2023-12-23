package loader_test

import (
	"fmt"
	. "github.com/monopole/mdparse/internal/usegold/loader"
	"os"

	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadFolderFromFs(t *testing.T) {
	type testC struct {
		root string
		path string
	}
	for n, tc := range map[string]testC{
		"t1": {
			root: "/home/jregan/myrepos/github.com/monopole",
			path: "mdparse",
		},
	} {
		t.Run(n, func(t *testing.T) {
			{
				var cwd string
				var err error
				cwd, err = os.Getwd()
				if err != nil {
					return
				}
				fmt.Println("cwd =", cwd)
				//				var folder MyFolder
				//				folder.name = stripTrailingSlash(cwd)
			}

			folder := NewFolder(tc.root)
			l := DefaultFsLoader
			f, err := l.LoadFolderFromFs(folder, tc.path)
			assert.NoError(t, err)
			f.Accept(&VisitorDump{})
		})
	}
}
