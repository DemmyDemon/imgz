package do

import (
	"fmt"
	"os"
)

func Fuck(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "FUCK:  %s\n", err)
		os.Exit(1)
	}
}
