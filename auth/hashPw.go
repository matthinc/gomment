package auth

import (
    "github.com/matthewhartstonge/argon2"
)

func HashPw(pw string) (hash string, err error) {
    cfg := argon2.DefaultConfig()
    encoded, err := cfg.HashEncoded([]byte(pw))
    if err != nil {
        return "", err
    }
    return string(encoded), nil
}

func ValidatePw(pw string, hash string) bool {
    raw, err := argon2.Decode([]byte(hash))
    if err != nil {
        return false
    }
    ok, err := raw.Verify([]byte(pw))
    if err != nil {
        return false
    }
    return ok
}
