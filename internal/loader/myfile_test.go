package loader_test

import (
	. "github.com/monopole/mdparse/internal/loader"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestMyFileEqualsEmpty(t *testing.T) {
	f1, f1Prime := NewEmptyFile("f1"), NewEmptyFile("f1")
	f2 := NewEmptyFile("f2")
	assert.True(t, f1.Equals(f1Prime))
	assert.False(t, f1.Equals(f2))
}

func TestMyFileEqualsFull(t *testing.T) {
	f1, f1Prime := NewFile("f1", []byte("f1")), NewFile("f1", []byte("f1"))
	f2 := NewFile("f2", []byte("f2"))
	assert.True(t, f1.Equals(f1Prime))
	assert.False(t, f1.Equals(f2))
}

func TestClean(t *testing.T) {
	// Just documenting behavior
	assert.Equal(t, ".", filepath.Clean(".///"))
	assert.Equal(t, "../..", filepath.Clean("./../../"))
	assert.Equal(t, "hoser", "./hoser"[2:])
}
