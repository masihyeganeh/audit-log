package auth

import (
	"github.com/brianvoe/sjwt"
	"github.com/pkg/errors"
	"time"
)

type Auth struct {
	secret          []byte
	N, r, p, keyLen int
}

type User struct {
	Username       string
	HasReadAccess  bool
	HasWriteAccess bool
}

func New(secret string) Auth {
	return Auth{
		secret: []byte(secret),
		// scrypt constants
		N:      16384,
		r:      8,
		p:      1,
		keyLen: 32,
	}
}

func (a *Auth) Login(user User) string {
	claims := sjwt.New()
	claims.Set("username", user.Username)
	claims.Set("has_read_access", user.HasReadAccess)
	claims.Set("has_write_access", user.HasWriteAccess)
	claims.SetTokenID()
	claims.SetIssuedAt(time.Now())
	claims.SetNotBeforeAt(time.Now())
	claims.SetExpiresAt(time.Now().Add(time.Hour * 24))
	return claims.Generate(a.secret)
}

func (a *Auth) Authenticate(jwtToken string) (*User, error) {
	// Verify that the secret signature is valid
	if !sjwt.Verify(jwtToken, a.secret) {
		return nil, errors.New("invalid jwt token")
	}

	// Parse jwt
	claims, err := sjwt.Parse(jwtToken)
	if err != nil {
		return nil, err
	}

	// Validate will check(if set) Expiration At and Not Before At dates
	err = claims.Validate()
	if err != nil {
		return nil, err
	}

	username, err := claims.Get("username")
	if err != nil {
		return nil, err
	}

	hasReadAccess, err := claims.Get("has_read_access")
	if err != nil {
		return nil, err
	}

	hasWriteAccess, err := claims.Get("has_write_access")
	if err != nil {
		return nil, err
	}

	return &User{
		Username:       username.(string),
		HasReadAccess:  hasReadAccess.(bool),
		HasWriteAccess: hasWriteAccess.(bool),
	}, nil
}
