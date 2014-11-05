package user

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type ValidateDigit struct {
	Value string
}

func NewValidateDigit(digit string) *ValidateDigit {
	return &ValidateDigit{digit}
}
func RandomValidateDigit(length int) *ValidateDigit {
	var rangeDigit = 1
	for i := 0; i < length; i++ {
		rangeDigit = rangeDigit * 10
	}
	random := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(rangeDigit)
	digit := fmt.Sprintf("%0[2]*[1]d", random, length)
	return &ValidateDigit{digit}
}
func (v *ValidateDigit) HasMatch(reqValue string) bool {
	var v1, v2 int
	var err error
	if v1, err = strconv.Atoi(v.Value); err != nil {
		return false
	}
	if v2, err = strconv.Atoi(reqValue); err != nil {
		return false
	}
	return v1 == v2
}
