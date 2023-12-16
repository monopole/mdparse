package loader

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTreeItem(t *testing.T) {
	var empty *myTreeItem
	assert.Equal(t, "", empty.Name())
	assert.Equal(t, "/", empty.AbsPath())
	assert.Equal(t, "{ERROR}", empty.DirName())

	bob := &myTreeItem{
		name: "bob",
	}
	assert.Equal(t, "bob", bob.Name())
	assert.Equal(t, "/bob", bob.AbsPath())
	assert.Equal(t, "/", bob.DirName())
	joe := myTreeItem{
		parent: bob,
		name:   "joe",
	}
	assert.Equal(t, "joe", joe.Name())
	assert.Equal(t, "/bob/joe", joe.AbsPath())
	assert.Equal(t, "/bob", joe.DirName())
}

func TestMyFolder(t *testing.T) {
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
	assert.Equal(t, "/bob", bob.AbsPath())
	assert.Equal(t, "/", bob.DirName())
	assert.Equal(t, "f1", f1.Name())
	assert.Equal(t, "/bob/f1", f1.AbsPath())
	assert.Equal(t, "/bob", f1.DirName())
	assert.Equal(t, "f2", f2.Name())
	assert.Equal(t, "/bob/f2", f2.AbsPath())
	assert.Equal(t, "/bob", f2.DirName())

	joe := MyFolder{
		myTreeItem: myTreeItem{
			name: "joe",
		},
	}
	joe.AddFolder(bob)

	assert.Equal(t, "joe", joe.Name())
	assert.Equal(t, "/joe", joe.AbsPath())
	assert.Equal(t, "/", joe.DirName())
	assert.Equal(t, "f2", f2.Name())
	assert.Equal(t, "/joe/bob/f2", f2.AbsPath())
	assert.Equal(t, "/joe/bob", f2.DirName())
}
