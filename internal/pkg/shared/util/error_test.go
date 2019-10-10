package util_test

import (
	"fmt"
	"go-brunel/internal/pkg/shared/util"
	"strings"
	"testing"

	"github.com/pkg/errors"
)

func TestErrorAppend(t *testing.T) {
	suites := []struct {
		a      error
		b      error
		expect error
	}{
		{
			a:      errors.New("a"),
			b:      errors.New("b"),
			expect: errors.New("a: b"),
		},
		{
			a:      errors.New("a"),
			b:      nil,
			expect: errors.New("a"),
		},
		{
			a:      nil,
			b:      errors.New("b"),
			expect: errors.New("b"),
		},
		{
			a:      nil,
			b:      nil,
			expect: nil,
		},
	}

	for i, suite := range suites {
		t.Run(
			fmt.Sprintf("suites[%d]", i),
			func(t *testing.T) {
				e := util.ErrorAppend(suite.a, suite.b)
				if e != nil && strings.Compare(e.Error(), suite.expect.Error()) != 0 {
					t.Errorf("'%s' does not match expected error '%s'", e.Error(), suite.expect.Error())
					t.Fail()
				} else if e == nil && e != suite.expect {
					t.Errorf("'%v' does not match expected error '%s'", e, suite.expect)
					t.Fail()
				}
			},
		)
	}
}
