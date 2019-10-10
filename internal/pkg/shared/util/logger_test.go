package util_test

import (
	"fmt"
	"go-brunel/internal/pkg/shared/util"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestLoggerWriter(t *testing.T) {
	suites := []struct {
		lines    []string
		expected []string
	}{
		{
			lines:    []string{"my line"},
			expected: []string{"my line"},
		},
		{
			lines:    []string{"my line\n"},
			expected: []string{"my line"},
		},
		{
			lines:    []string{"my\nline\n"},
			expected: []string{"my", "line"},
		},
		{
			lines:    []string{"my\n", "line\n"},
			expected: []string{"my", "line"},
		},
		{
			lines:    []string{"my\n", "line\n"},
			expected: []string{"my", "line"},
		},
	}

	for i, suite := range suites {
		t.Run(
			fmt.Sprintf("suites[%d]", i),
			func(t *testing.T) {
				var lines []string
				logWriter := util.LoggerWriter{
					Recorder: func(log string) error {
						lines = append(lines, log)
						return nil
					},
				}
				for _, line := range suite.lines {
					if _, err := logWriter.Write([]byte(line)); err != nil {
						t.Fatal(err)
					}
				}
				_ = logWriter.Close()

				if !gomock.Eq(lines).Matches(suite.expected) {
					t.Error(lines, "does not match", suite.expected)
					t.Fail()
				}
			},
		)
	}
}
