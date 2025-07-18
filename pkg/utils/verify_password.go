package utils

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"golang.org/x/crypto/argon2"
	"strings"
)

func VerifyPassword(pass string, digest string) (bool, error) {
	parts := strings.Split(digest, ".")
	if len(parts) != 2 {
		return false, errors.New("invalid encoded hash format")
	}

	encodedSalt := parts[0]
	encodedHash := parts[1]

	salt, err := base64.StdEncoding.DecodeString(encodedSalt)
	if err != nil {
		return false, err
	}

	passHash, err := base64.StdEncoding.DecodeString(encodedHash)
	if err != nil {
		return false, err
	}

	hash := argon2.IDKey([]byte(pass), salt, 1, 64*1024, 4, 32)
	if len(hash) != len(passHash) {
		return false, nil
	}

	return subtle.ConstantTimeCompare(hash, passHash) == 1, nil
}
