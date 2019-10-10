package main

import (
	"bytes"
	"fmt"
	"go-brunel/internal/pkg/shared/remote"
	"gopkg.in/yaml.v2"
)

func main() {
	c, err := remote.GenerateCredentials()
	if err != nil {
		println(err)
	}
	b := bytes.Buffer{}

	err = yaml.NewEncoder(&b).Encode(c)
	if err != nil {
		println(fmt.Errorf("error encoding yaml: %v", err))
	}
	println(b.String())
}
