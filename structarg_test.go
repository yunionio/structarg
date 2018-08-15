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
