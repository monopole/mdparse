package loader_test

import (
	. "github.com/monopole/mdparse/internal/usegold/loader"

	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	f1      = NewFile("f1")
	f1Prime = NewFile("f1")
	f2      = NewFile("f2")
	f2Prime = NewFile("f2")
)

func TestMyFolderNaming(t *testing.T) {
	f1 = NewFile("f1")
	f2 = NewFile("f2")

	d1 := NewFolder("d1")
	d1.AddFileObject(f1)
	d1.AddFileObject(f2)

	assert.Equal(t, "d1", d1.Name())
	assert.Equal(t, "/d1", d1.FullName())
	assert.Equal(t, "/", d1.DirName())
	assert.Equal(t, "f1", f1.Name())
	assert.Equal(t, "/d1/f1", f1.FullName())
	assert.Equal(t, "/d1", f1.DirName())
	assert.Equal(t, "f2", f2.Name())
	assert.Equal(t, "/d1/f2", f2.FullName())
	assert.Equal(t, "/d1", f2.DirName())

	d2 := NewFolder("d2")
	d2.AddFolderObject(d1)

	assert.Equal(t, "d2", d2.Name())
	assert.Equal(t, "/d2", d2.FullName())
	assert.Equal(t, "/", d2.DirName())
	assert.Equal(t, "f2", f2.Name())
	assert.Equal(t, "/d2/d1/f2", f2.FullName())
	assert.Equal(t, "/d2/d1", f2.DirName())
}

func TestMyFolderEquals(t *testing.T) {

	d1 := NewFolder("d1")
	d1.AddFileObject(f1)
	d1.AddFileObject(f2)

	d1Prime := NewFolder("d1")
	d1Prime.AddFileObject(f1Prime)
	d1Prime.AddFileObject(f2Prime)

	assert.True(t, d1.Equals(d1))
	assert.True(t, d1.Equals(d1Prime))

	d2 := NewFolder("d2")
	d2.AddFolderObject(d1)
	d2Prime := NewFolder("d2")
	d2Prime.AddFolderObject(d1Prime)

	assert.True(t, d2.Equals(d2))
	assert.True(t, d2.Equals(d2Prime))
	assert.False(t, d2.Equals(d1))
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
			assert.Equal(t, "mdparse", f.Name())
		})
	}
}
