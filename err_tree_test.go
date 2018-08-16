package errtree

import (
	"errors"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("root-only tree", func(t *testing.T){
		t.Parallel()
		subject := New(errors.New("foo"))
		expect := "foo"
		if result := subject.Error(); result != expect {
			t.Errorf("expected %s but got %s", expect, result)
		}
	})

	t.Run("descendants in tree", func(t *testing.T) {
		t.Parallel()
		subject := New(errors.New("foo"), errors.New("bar"))
		expect := "foo:\n    bar"
		if result := subject.Error(); result != expect {
			t.Errorf("expected %s but got %s", expect, result)
		}
	})
}

func TestNewString(t *testing.T) {
	t.Parallel()

	t.Run("root-only tree", func(t *testing.T){
		t.Parallel()
		subject := NewString("foo")
		expect := "foo"
		if result := subject.Error(); result != expect {
			t.Errorf("expected %s but got %s", expect, result)
		}
	})

	t.Run("descendants in tree", func(t *testing.T) {
		t.Parallel()
		subject := NewString("foo", NewString("bar"))
		expect := "foo:\n    bar"
		if result := subject.Error(); result != expect {
			t.Errorf("expected %s but got %s", expect, result)
		}
	})
}

func Test_errorTree_Error(t *testing.T) {
	t.Parallel()

	expect := strings.Trim(`
a:
    b:
        c:
            D
    E
    f:
        g:
            H
    i
    J
    k
`, "\n")
	subject := &ErrorTree{
		err: errors.New("a"),
		descendants: []error{
			&ErrorTree{
				err: errors.New("b"),
				descendants: []error{
					&ErrorTree{
						err: errors.New("c"),
						descendants: []error{
							nil,
							errors.New("D"),
						},
					},
				},
			},
			errors.New("E"),
			&ErrorTree{
				errors.New("f"),
				[]error{
					&ErrorTree{
						err: errors.New("g"),
						descendants: []error{
							errors.New("H"),
							nil,
						},
					},
				},
			},
			&ErrorTree{
				err:         errors.New("i"),
				descendants: make([]error, 0),
			},
			nil,
			errors.New("J"),
			nil,
			&ErrorTree{
				err: errors.New("k"),
			},
		},
	}
	if result := subject.Error(); result != expect {
		t.Errorf("expected %q but got: %q", expect, result)
	}
}
