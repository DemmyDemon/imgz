package do

import (
	"fmt"
)

var Verbosity = true

func Verbose(stuff ...any) (int, error) {
	if !Verbosity {
		return 0, nil
	}
	return fmt.Println(stuff...)
}
