package loader

import (
	"github.com/monopole/mdrip/base"
	"testing"
)

func Test_codeBlock_HasLabel(t *testing.T) {
	tests := map[string]struct {
		cb    CodeBlock
		label base.Label
		want  bool
	}{
		"t1": {
			cb: CodeBlock{
				labels: nil,
			},
			label: "sss",
			want:  false,
		},
		"t2": {
			cb: CodeBlock{
				labels: []base.Label{"protein", base.SleepLabel},
			},
			label: "protein",
			want:  true,
		},
		"t3": {
			cb: CodeBlock{
				labels: []base.Label{base.WildCardLabel, base.SleepLabel},
			},
			label: base.WildCardLabel,
			want:  true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := tc.cb.HasLabel(tc.label); got != tc.want {
				t.Errorf("HasLabel(%s) = %v, want %v", tc.label, got, tc.want)
			}
		})
	}
}
