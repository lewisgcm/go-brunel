package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// ParseQueryTime will parse a query parameter as time.Time, returning an error on failure.
// Only unix timestamps are supported as the time format.
func ParseQueryTime(
	r *http.Request,
	name string,
	required bool,
	fallback time.Time,
) (time.Time, error) {
	i, e := ParseQueryInt(r, name, required, fallback.Unix())
	if e != nil {
		return fallback, errors.Wrap(e, "error parsing query parameter to unix timestamp")
	}
	return time.Unix(i, 0), nil
}

// ParseQueryInt will attempt to parse a int64 from a query parameter, returning an error on failure.
func ParseQueryInt(
	r *http.Request,
	name string,
	required bool,
	fallback int64,
) (int64, error) {
	q := r.URL.Query().Get(name)
	if q == "" {
		if !required {
			return fallback, nil
		}
		return fallback, fmt.Errorf("the query parameter '%s' should be specified", name)
	}
	i, err := strconv.ParseInt(q, 10, 64)
	if err != nil {
		return fallback, errors.Wrap(err, fmt.Sprintf("error parsing '%s' query parameter to integer", name))
	}
	return i, nil
}
