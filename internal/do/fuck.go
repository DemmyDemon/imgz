package do

import (
	"fmt"
	"os"
)

func Fuck(context string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s:  %s\n", context, err)
		os.Exit(1)
	}
}
