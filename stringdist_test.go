package structarg

import "testing"

func TestFindSimilar(t *testing.T) {
	cases := []struct{
		niddle string
		candidates []string
		want []string
	} {
		{"a", []string{"ab", "a"}, []string{"a", "ab"}},
		{"abc", []string{"ab", "abc", "abcd", "abcdef", "xyz"}, []string{"abc", "ab", "abcd"}},
		{"abc", []string{"abcd", "abc", "ab", "abcdef"}, []string{"abc", "ab", "abcd"}},
	}
	for _, tt := range cases {
		got := FindSimilar(tt.niddle, tt.candidates, -1, 0.7)
		t.Logf("%#v", got)
	}
}

func TestChoicesString(t *testing.T) {
	cases := []struct{
		candidates []string
		want string
	} {
		{[]string{"ab", "a"}, "ab or a"},
		{[]string{"ab", "abc", "abcd", "abcdef"}, "ab abc abcd or abcdef"},
	}
	for _, tt := range cases {
		got := ChoicesString(tt.candidates)
		t.Logf("%#v", got)
	}
}
