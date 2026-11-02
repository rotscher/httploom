package weave

import "testing"

func TestWrap(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		width int
		want  string
	}{
		{
			name:  "wraps at width",
			text:  "The quick brown fox jumps over the lazy dog",
			width: 20,
			want:  "The quick brown fox\njumps over the lazy\ndog",
		},
		{
			name:  "keeps text that fits",
			text:  "short text",
			width: 80,
			want:  "short text",
		},
		{
			name:  "preserves explicit newlines",
			text:  "one two\nthree four",
			width: 80,
			want:  "one two\nthree four",
		},
		{
			name:  "preserves blank lines",
			text:  "one\n\ntwo",
			width: 80,
			want:  "one\n\ntwo",
		},
		{
			name:  "does not split long words",
			text:  "abcdefghij",
			width: 5,
			want:  "abcdefghij",
		},
		{
			name:  "uses default width for non-positive values",
			text:  "a b",
			width: 0,
			want:  "a b",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Wrap(tt.text, tt.width); got != tt.want {
				t.Fatalf("Wrap(%q, %d) = %q, want %q", tt.text, tt.width, got, tt.want)
			}
		})
	}
}
