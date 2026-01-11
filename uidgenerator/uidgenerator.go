package uidgenerator

import (
	"crypto/rand"
	"strconv"
	"strings"
	"time"
)

const charset = "ABCDFGHIJKLMNOPQRSTUVXYZ"
const charsetNum = "0123456789"

func Generate() string {
	t := strconv.FormatInt(time.Now().UnixNano(), 10)
	t = t[:5]

	b := make([]byte, 5)
	_, _ = rand.Read(b)

	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}

	return t + string(b)
}

func Validate(uid string) bool {
	if len(uid) != 10 {
		return false
	}

	left := uid[:5]  // numeric part
	right := uid[5:] // alphabetic part

	for i := range left {
		if !strings.ContainsRune(charsetNum, rune(left[i])) {
			return false
		}
	}

	for i := range right {
		if !strings.ContainsRune(charset, rune(right[i])) {
			return false
		}
	}

	return true
}
