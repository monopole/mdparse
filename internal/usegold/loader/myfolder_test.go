package loader

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestMyFolderWhiteBox(t *testing.T) {
	f1 := &MyFile{
		myTreeItem: myTreeItem{
			name: "f1",
		},
	}
	f2 := &MyFile{
		myTreeItem: myTreeItem{
			name: "f2",
		},
	}

	bob := &MyFolder{
		myTreeItem: myTreeItem{
			name: "bob",
		},
	}
	bob.AddFileObject(f1)
	bob.AddFileObject(f2)

	assert.Equal(t, "bob", bob.Name())
	assert.Equal(t, "/bob", bob.FullName())
	assert.Equal(t, "/", bob.DirName())
	assert.Equal(t, "f1", f1.Name())
	assert.Equal(t, "/bob/f1", f1.FullName())
	assert.Equal(t, "/bob", f1.DirName())
	assert.Equal(t, "f2", f2.Name())
	assert.Equal(t, "/bob/f2", f2.FullName())
	assert.Equal(t, "/bob", f2.DirName())

	joe := MyFolder{
		myTreeItem: myTreeItem{
			name: "joe",
		},
	}
	joe.AddFolderObject(bob)

	assert.Equal(t, "joe", joe.Name())
	assert.Equal(t, "/joe", joe.FullName())
	assert.Equal(t, "/", joe.DirName())
	assert.Equal(t, "f2", f2.Name())
	assert.Equal(t, "/joe/bob/f2", f2.FullName())
	assert.Equal(t, "/joe/bob", f2.DirName())
}

func TestClean(t *testing.T) {
	// Just documenting behavior
	assert.Equal(t, ".", filepath.Clean(".///"))
	assert.Equal(t, "../..", filepath.Clean("./../../"))
	assert.Equal(t, "hoser", "./hoser"[2:])
}

func TestLoadFolder(t *testing.T) {
	type testC struct {
		arg string
	}
	for n, tc := range map[string]testC{
		//"t1": {
		//	arg: "/home/jregan/myrepos/github.com/monopole/mdparse",
		//},
		//"t2": {
		//	arg: ".",
		//},
	} {
		t.Run(n, func(t *testing.T) {
			f, err := LoadFolder(DefaultFsLoader, tc.arg)
			assert.NoError(t, err)
			f.Accept(&VisitorDump{})
			assert.Equal(t, "mdparse", f.name)
		})
	}
}
