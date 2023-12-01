package usegold

import "testing"

func Test_commentBody(t *testing.T) {
	tests := map[string]struct {
		data string
		want string
	}{
		"hoser": {
			data: "<!--hello-->",
			want: "hello",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := commentBody(tc.data); got != tc.want {
				t.Errorf("commentBody() = %v, want %v", got, tc.want)
			}
		})
	}
}
