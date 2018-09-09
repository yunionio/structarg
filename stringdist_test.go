package structarg

import "testing"

func identicalStringArray(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i += 1 {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestFindSimilar(t *testing.T) {
	cases := []struct{
		niddle string
		candidates []string
		want []string
	} {
		{"a", []string{"ab", "a"}, []string{"a", "ab"}},
		{"abc", []string{"ab", "abc", "abcd", "abcdef", "xyz"}, []string{"abc", "abcd", "ab", "abcdef"}},
		{"abc", []string{"abcd", "abc", "ab", "abcdef"}, []string{"abc", "abcd", "ab", "abcdef"}},
	}
	for _, tt := range cases {
		got := FindSimilar(tt.niddle, tt.candidates, -1, 0.5)
		t.Logf("%#v", got)
		if ! identicalStringArray(tt.want, got) {
			t.Errorf("want %#v got %#v", tt.want, got)
		}
	}
}

func TestChoicesString(t *testing.T) {
	cases := []struct{
		candidates []string
		want string
	} {
		{[]string{"ab", "a"}, "ab or a"},
		{[]string{"ab", "abc", "abcd", "abcdef"}, "ab, abc, abcd or abcdef"},
	}
	for _, tt := range cases {
		got := ChoicesString(tt.candidates)
		t.Logf("%#v", got)
		if got != tt.want {
			t.Errorf("want %#v got %#v", tt.want, got)
		}
	}
}
