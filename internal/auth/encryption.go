package auth

import (
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/crypto/scrypt"
)

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	return b, err
}

func (a *Auth) EncryptPasswordWithSalt(password, salt string) (string, error) {
	s, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return "", err
	}

	dk, err := scrypt.Key([]byte(password), s, a.N, a.r, a.p, a.keyLen)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(dk), nil
}

func (a *Auth) EncryptPassword(password string) (string, string, error) {
	salt, err := generateRandomBytes(8)
	if err != nil {
		return "", "", err
	}

	s := base64.StdEncoding.EncodeToString(salt)

	hashedPassword, err := a.EncryptPasswordWithSalt(password, s)
	if err != nil {
		return "", "", err
	}

	return hashedPassword, s, nil
}
