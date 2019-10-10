/*
 * Author: Lewis Maitland
 *
 * Copyright (c) 2019 Lewis Maitland
 */

package parser

import (
	"go-brunel/internal/pkg/shared"
	"io"
)

type Parser interface {
	Parse(fileName string, reader io.Reader, progress io.Writer) (*shared.Spec, error)
}
