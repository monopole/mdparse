package accum

import (
	"bytes"
	_ "embed"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTxtPrinter(t *testing.T) {
	ld, err := NewLessonDocFromBytes(makeParser(), []byte(mds))
	assert.NoError(t, err)
	assert.NotNil(t, ld)
	var b bytes.Buffer
	v := NewTutorialTxtPrinter(&b)
	ld.Accept(v)
	assert.Equal(t, `NO IDEA WHAT THIS IS CALLED
  one --- a := one...
  clickToCopy --- c := two...
  four --- c := three...
  clickToCopy --- c := four...
  blimp --- c := five d := six...
`, b.String())
}
