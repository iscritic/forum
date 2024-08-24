package utils

import (
	"errors"
	"strconv"
)

// Etoi - An enchanted Atoi, accepting only positive numbers and rejecting those starting with zero.
func Etoi(s string) (int, error) {
	num, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	if num <= 0 {
		return 0, errors.New("value must be a positive integer")
	}

	if s[0] == '0' {
		return 0, errors.New("value cannot start with zero")
	}

	return num, nil
}
