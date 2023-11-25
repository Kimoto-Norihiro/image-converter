package myerror

import "fmt"

func Wrap(errp *error, info string) {
	if *errp != nil {
		*errp = fmt.Errorf("%s: %w", info, *errp)
	}
}
