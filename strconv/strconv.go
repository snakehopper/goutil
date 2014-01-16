package strconv

import (
	"strconv"
)

func IncreaseInt(num string) (string, error) {
	inum, err := strconv.Atoi(num)
	if err != nil {
		return "", err
	}
	inum = inum + 1
	return strconv.Itoa(inum), nil
}

func DecreaseInt(num string) (string, error) {
	inum, err := strconv.Atoi(num)
	if err != nil {
		return "", err
	}
	inum = inum - 1
	return strconv.Itoa(inum), nil
}
