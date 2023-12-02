package usegold

import (
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
		want []string
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
			want: []string{"aa", "b", "ccc"},
		},
		"t5": {
			data: "  @aa @b  @   @@ccc @@@ @@@d ",
			want: []string{"aa", "b", "ccc", "d"},
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
