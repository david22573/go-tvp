package term

import (
	"os"

	"golang.org/x/term"
)

func Size() (width, height int, err error) {
	fd := int(os.Stdout.Fd())
	return term.GetSize(fd)
}
