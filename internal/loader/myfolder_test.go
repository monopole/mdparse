package loader_test

import (
	. "github.com/monopole/mdparse/internal/loader"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMyFolderNaming(t *testing.T) {
	f1 := NewEmptyFile("f1")
	f2 := NewEmptyFile("f2")
	d1 := NewFolder("d1").AddFileObject(f1).AddFileObject(f2)
	assert.Equal(t, "d1", d1.Name())
	assert.Equal(t, "d1", d1.FullName())
	assert.Equal(t, "f1", f1.Name())
	assert.Equal(t, "d1/f1", f1.FullName())
	assert.Equal(t, "f2", f2.Name())
	assert.Equal(t, "d1/f2", f2.FullName())

	d2 := NewFolder("d2").AddFolderObject(d1)
	assert.Equal(t, "d2", d2.Name())
	assert.Equal(t, "d2", d2.FullName())
	assert.Equal(t, "f2", f2.Name())
	assert.Equal(t, "d2/d1/f2", f2.FullName())
}

func TestMyFolderEquals(t *testing.T) {
	f1, f1Prime := NewEmptyFile("f1"), NewEmptyFile("f1")
	f2, f2Prime := NewEmptyFile("f2"), NewEmptyFile("f2")

	d1 := NewFolder("d1").AddFileObject(f1).AddFileObject(f2)
	d1Prime := NewFolder("d1").AddFileObject(f1Prime).AddFileObject(f2Prime)

	assert.True(t, d1.Equals(d1))
	assert.True(t, d1.Equals(d1Prime))

	d2 := NewFolder("d2").AddFolderObject(d1)
	d2Prime := NewFolder("d2").AddFolderObject(d1Prime)

	assert.True(t, d2.Equals(d2))
	assert.True(t, d2.Equals(d2Prime))
	assert.False(t, d2.Equals(d1))
}
