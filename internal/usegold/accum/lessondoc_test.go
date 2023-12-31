package accum

import (
	_ "embed"
	"github.com/monopole/mdrip/base"
	"github.com/stretchr/testify/assert"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"slices"
	"testing"
)

func Test_commentBody(t *testing.T) {
	tests := map[string]struct {
		data string
		want string
	}{
		"t1": {
			data: "<!--hello-->\n",
			want: "hello",
		},
		"t2": {
			data: "<!-- hello -->",
			want: " hello ",
		},
		"t3": {
			data: "<!- hello -->",
			want: "",
		},
		"t4": {
			data: "<!-- hello ->",
			want: "",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := commentBody(tc.data); got != tc.want {
				t.Errorf("got = %v, want %v", got, tc.want)
			}
		})
	}
}

func Test_parseLabels(t *testing.T) {
	tests := map[string]struct {
		data string
		want []base.Label
	}{
		"t1": {
			data: "",
			want: nil,
		},
		"t2": {
			data: "    ",
			want: nil,
		},
		"t3": {
			data: "   aaa ",
			want: nil,
		},
		"t4": {
			data: "  @aa @b     @ccc ",
			want: []base.Label{"aa", "b", "ccc"},
		},
		"t5": {
			data: "  @aa @b  @   @@ccc @@@ @@@d ",
			want: []base.Label{"aa", "b", "ccc", "d"},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := parseLabels(tc.data); !slices.Equal(got, tc.want) {
				t.Errorf("got = %v, want %v", got, tc.want)
			}
		})
	}
}

func makeParser() parser.Parser {
	return goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(),
		),
	).Parser()
}

//go:embed testdata/small.md
var mds string

func TestNewLessonDocFromPath(t *testing.T) {
	ld, err := NewLessonDocFromBytes(makeParser(), []byte(mds))
	assert.NoError(t, err)
	assert.NotNil(t, ld)
	// TODO: this is dumb
	assert.Empty(t, string(ld.Path()))
}
