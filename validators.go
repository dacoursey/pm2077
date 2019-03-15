package main

import (
	"errors"
	"strings"
)

func validatePassword(pwd string) error {

	if len(pwd) < 6 {
		return errors.New("pwdShort")
	} else if !(strings.ContainsAny(pwd, "abcdefghijklmnopqrstuvwxyz") && strings.ContainsAny(pwd, "1234567890")) {
		return errors.New("pwdNotComplex")
	}

	return nil
}
