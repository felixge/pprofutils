package utils

import "io"

type Anonymize struct {
	Input  io.Reader
	Output io.Writer
}
