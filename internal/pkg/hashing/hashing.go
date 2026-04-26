package hashing

import (
	"crypto/rand"
	"encoding/base64"
	"errors"

	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidHash   = errors.New("invalid hash format")
	ErrHashMismatch = errors.New("password does not match hash")
)

type Params struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

type Hasher struct {
	params Params
}

func NewHasher(params Params) *Hasher {
	return &Hasher{params: params}
}

func (h *Hasher) HashPassword(password string) (string, error) {
	salt := make([]byte, h.params.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		h.params.Iterations,
		h.params.Memory,
		h.params.Parallelism,
		h.params.KeyLength,
	)

	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)

	return encodedSalt + ":" + encodedHash, nil
}

func (h *Hasher) VerifyPassword(hashStr, password string) error {
	parts := split(hashStr, ":")
	if len(parts) != 2 {
		return ErrInvalidHash
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return ErrInvalidHash
	}

	originalHash, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return ErrInvalidHash
	}

	verifyHash := argon2.IDKey(
		[]byte(password),
		salt,
		h.params.Iterations,
		h.params.Memory,
		h.params.Parallelism,
		h.params.KeyLength,
	)

	if string(verifyHash) != string(originalHash) {
		return ErrHashMismatch
	}

	return nil
}

func split(s, sep string) []string {
	result := []string{}
	current := ""
	for _, c := range s {
		if string(c) == sep {
			result = append(result, current)
			current = ""
		} else {
			current += string(c)
		}
	}
	result = append(result, current)
	return result
}

func DefaultParams() Params {
	return Params{
		Memory:      65536,
		Iterations:  3,
		Parallelism: 4,
		SaltLength: 16,
		KeyLength:  32,
	}
}