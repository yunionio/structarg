package structarg

import (
	"testing"
)

func newParser(d interface{}) (*ArgumentParser, error) {
	p, err := NewArgumentParser(
		d,
		"prog",
		"prog desc",
		"prog epilog",
	)
	return p, err
}

func TestPositional(t *testing.T) {
	p, err := newParser(
		&struct {
			POS            string
			NONPOSREQUIRED string `positional:false required:true`
			NONPOS         string `positional:false`
		}{},
	)
	if err != nil {
		t.Errorf("error parsing %s", err)
		return
	}
	if len(p.posArgs) != 1 {
		t.Errorf("wrong number of positionals, want 1, got %d", len(p.posArgs))
		return
	} else {
		a := p.posArgs[0]
		if !a.IsPositional() {
			t.Errorf("expecting positional argument, got %s", a)
			return
		}
		if !a.IsRequired() {
			t.Errorf("positional argument should be required, got %s", a)
			return
		}
		if a.String() != "<POS>" {
			t.Errorf("wrong usage %s", a)
			return
		}
	}
	if len(p.optArgs) != 2 {
		t.Errorf("wrong number of optionals, want 2, got %d", len(p.optArgs))
		return
	} else {
		var a Argument
		a = p.optArgs[0]
		if a.IsPositional() {
			t.Errorf("expecting positional argument, got %s", a)
			return
		}
		if a.IsRequired() {
			t.Errorf("expecting non-required optional argument, got %s", a)
			return
		}
		if a.String() != "[--nonpos NONPOS]" {
			t.Errorf("wrong usage %s", a)
			return
		}

		a = p.optArgs[1]
		if a.IsPositional() {
			t.Errorf("expecting positional argument, got %s", a)
			return
		}
		if !a.IsRequired() {
			t.Errorf("expecting required optional argument, got %s", a)
			return
		}
		if a.String() != "<--nonposrequired NONPOSREQUIRED>" {
			t.Errorf("wrong usage %s", a)
			return
		}
	}
}

func TestPositionalAppend(t *testing.T) {
	p, err := newParser(
		&struct {
			Opt0 string
			Opt1 string `required:true`
			Opt2 string
			Opt3 string `required:true`
			Opt4 string
		}{},
	)
	if err != nil {
		t.Errorf("err expected: %s", err)
		return
	}
	if len(p.optArgs) != 5 {
		t.Errorf("num optionals want 5, got %d", len(p.optArgs))
		return
	}
	// check order
	required := false
	for i, arg := range p.optArgs {
		if !arg.IsRequired() {
			if required {
				t.Errorf("bad order at %d", i)
				break
			}
		} else {
			if required == false {
				required = true
			}
		}
	}
}

func TestRequired(t *testing.T) {
	t.Run("optional positional", func(t *testing.T) {
		_, err := newParser(
			&struct {
				POS string `required:false`
			}{},
		)
		if err == nil {
			t.Errorf("no error for optional positional")
		}
	})
	t.Run("default positional", func(t *testing.T) {
		_, err := newParser(
			&struct {
				POS string `default:baddefault`
			}{},
		)
		if err == nil {
			t.Errorf("no error for positional with default")
		}
	})
	t.Run("required non-positional", func(t *testing.T) {
		p, err := newParser(
			&struct {
				RequiredOpt string `required:true`
				Opt         string
			}{},
		)
		if err != nil {
			t.Errorf("errored: %s", err)
			return
		}
		if len(p.optArgs) != 2 {
			t.Errorf("expecting 2 optArgs, got %d", len(p.optArgs))
			return
		}
		// the required ones should come at last
		a := p.optArgs[1]
		if !a.IsRequired() {
			t.Errorf("argument %s is optional, want required", a)
			return
		}
		if want := "<--required-opt REQUIRED_OPT>"; a.String() != want {
			t.Errorf("want %s, got %s", want, a.String())
		}
	})
}

func TestNonPositionalRequiredWithDefault(t *testing.T) {
	_, err := newParser(
		&struct {
			Opt int `default:100 required:true`
		}{},
	)
	if err == nil {
		t.Errorf("should error for non-positional argument with default value and required attribute")
	}
}

func TestBoolField(t *testing.T) {
	t.Run("default (no flags)", func(t *testing.T) {
		s := &struct {
			Bool              bool
			BoolP             *bool
			BoolDefaultTrue   bool  `default:true`
			BoolPDefaultTrue  *bool `default:true`
			BoolDefaultFalse  bool  `default:false`
			BoolPDefaultFalse *bool `default:false`
		}{}
		p, err := newParser(s)
		if err != nil {
			t.Fatalf("newParser failed: %s", err)
		}
		args := []string{}
		err = p.ParseArgs(args, false)
		if err != nil {
			t.Fatalf("ParseArgs failed: %s", err)
		}
		if !(!s.Bool && s.BoolP == nil &&
			s.BoolDefaultTrue && s.BoolPDefaultTrue != nil && *s.BoolPDefaultTrue &&
			!s.BoolDefaultFalse && s.BoolPDefaultFalse != nil && !*s.BoolPDefaultFalse) {
			t.Errorf("wrong parse result: %#v", s)
		}
	})
	t.Run("--flags", func(t *testing.T) {
		s := &struct {
			Bool              bool
			BoolP             *bool
			BoolDefaultTrue   bool  `default:true`
			BoolPDefaultTrue  *bool `default:true`
			BoolDefaultFalse  bool  `default:false`
			BoolPDefaultFalse *bool `default:false`
		}{}
		p, err := newParser(s)
		if err != nil {
			t.Fatalf("newParser failed: %s", err)
		}
		args := []string{
			"--bool",
			"--bool-p",
			"--bool-default-true",
			"--bool-p-default-true",
			"--bool-default-false",
			"--bool-p-default-false",
		}
		err = p.ParseArgs(args, false)
		if err != nil {
			t.Fatalf("ParseArgs failed: %s", err)
		}
		if !(s.Bool && s.BoolP != nil && *s.BoolP &&
			!s.BoolDefaultTrue && s.BoolPDefaultTrue != nil && !*s.BoolPDefaultTrue &&
			s.BoolDefaultFalse && s.BoolPDefaultFalse != nil && *s.BoolPDefaultFalse) {
			t.Errorf("wrong parse result: %#v", s)
		}
	})
}

func TestChoices(t *testing.T) {
	s := &struct {
		String string `choices:tcp|udp|http|https`
	}{}
	p, err := newParser(s)
	if err != nil {
		t.Fatalf("newParser failed: %s", err)
	}
	args := []string{"--string", ""}
	t.Run("good choices", func(t *testing.T) {
		choices := []string{"tcp", "udp", "http", "https"}
		for _, choice := range choices {
			args[1] = choice
			err = p.ParseArgs(args, false)
			if err != nil {
				t.Fatalf("ParseArgs failed: %s", err)
			}
			if s.String != choice {
				t.Errorf("wrong parse result: want %q, got %q", choice, s.String)
			}
		}
	})
	t.Run("bad choices", func(t *testing.T) {
		choices := []string{"", "et"}
		for _, choice := range choices {
			args[1] = choice
			err = p.ParseArgs(args, false)
			if err == nil {
				t.Fatalf("ParseArgs should error")
			}
			if s.String != "" {
				t.Errorf("Struct member should not be set, got %s", s.String)
			}
		}
	})
}

func TestArgValue(t *testing.T) {
	s := &struct {
		String string
	}{}
	p, err := newParser(s)
	if err != nil {
		t.Fatalf("newParser failed: %s", err)
	}
	args := []string{"--string", ""}
	cases := []struct {
		name  string
		value string
	}{
		{
			name:  "with space",
			value: `Hello world`,
		},
		{
			name:  "with single quote",
			value: `'Hello 'world'`,
		},
		{
			name:  "with double quote",
			value: `"Hello "world"`,
		},
		{
			name:  "with newline and tab",
			value: `Hello\n\tworld\n`,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			args[1] = c.value
			err := p.ParseArgs(args, false)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if s.String != c.value {
				t.Errorf("want %s, got %s", c.value, s.String)
			}
		})
	}
}

func TestIgnoreUnexported(t *testing.T) {
	s := &struct {
		unexported string
	}{}
	p, err := newParser(s)
	if err != nil {
		t.Fatalf("newParser failed: %s", err)
	}
	args := []string{"--string", ""}
	p.ParseArgs(args, true)
}
