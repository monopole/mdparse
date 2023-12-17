package loader

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMyFolderSilly(t *testing.T) {
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
	bob.AddFile(f1)
	bob.AddFile(f2)

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
	joe.AddFolder(bob)

	assert.Equal(t, "joe", joe.Name())
	assert.Equal(t, "/joe", joe.FullName())
	assert.Equal(t, "/", joe.DirName())
	assert.Equal(t, "f2", f2.Name())
	assert.Equal(t, "/joe/bob/f2", f2.FullName())
	assert.Equal(t, "/joe/bob", f2.DirName())
}

func TestMyFolderAbsorb(t *testing.T) {

}
