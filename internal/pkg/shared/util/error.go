package util

import "fmt"

func ErrorAppend(e error, a error) error {
	if e != nil {
		if a != nil {
			return fmt.Errorf("%s: %s", e.Error(), a.Error())
		}
		return e
	}
	return a
}
