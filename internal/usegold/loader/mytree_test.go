package loader

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTreeItem(t *testing.T) {
	var empty *myTreeItem
	assert.Equal(t, "", empty.Name())
	assert.Equal(t, "/", empty.FullName())
	assert.Equal(t, "{ERROR}", empty.DirName())

	bob := &myTreeItem{
		name: "bob",
	}
	assert.Equal(t, "bob", bob.Name())
	assert.Equal(t, "/bob", bob.FullName())
	assert.Equal(t, "/", bob.DirName())
	joe := myTreeItem{
		parent: bob,
		name:   "joe",
	}
	assert.Equal(t, "joe", joe.Name())
	assert.Equal(t, "/bob/joe", joe.FullName())
	assert.Equal(t, "/bob", joe.DirName())
}
