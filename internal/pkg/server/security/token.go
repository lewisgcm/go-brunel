package security

import (
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

const (
	authHeader = "Authorization"
)

// TokenSerializer will encode/decode security tokens to and from http requests and responses
type TokenSerializer interface {
	// Encode should encode the Identity as a string
	Encode(identity Identity) (string, error)

	// Decode should decode the http request into an Identity
	Decode(r *http.Request) (*Identity, error)
}

// NewTokenSerializer will
func NewTokenSerializer(s []byte, role UserRole) TokenSerializer {
	return &jwtSerializer{
		hMACSecret:  s,
		defaultRole: role,
	}
}

type jwtSerializer struct {
	hMACSecret  []byte
	defaultRole UserRole
}

type roleClaims struct {
	jwt.StandardClaims
	Role UserRole `json:"role,omitempty"`
}

func (p *jwtSerializer) Encode(identity Identity) (string, error) {
	r := identity.Role
	if r == "" {
		r = p.defaultRole
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, roleClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   identity.Username,
			NotBefore: time.Now().Unix(),
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 25).Unix(),
		},
		Role: r,
	})

	t, err := token.SignedString(p.hMACSecret)
	if err != nil {
		return "", errors.New("error generating token")
	}
	return t, nil
}

func (p *jwtSerializer) getToken(r *http.Request) *string {
	queryToken := r.URL.Query().Get("token")
	if queryToken != "" {
		return &queryToken
	}

	parts := strings.Split(r.Header.Get(authHeader), "Bearer ")
	if len(parts) == 2 {
		return &parts[1]
	}

	return nil
}

func (p *jwtSerializer) Decode(r *http.Request) (*Identity, error) {
	token := p.getToken(r)
	if token != nil {
		claims := &roleClaims{}
		tkn, err := jwt.ParseWithClaims(*token, claims, func(token *jwt.Token) (interface{}, error) {
			return p.hMACSecret, nil
		})
		if err != nil {
			return nil, err
		}
		if !tkn.Valid {
			return nil, errors.New("invalid token")
		}
		return &Identity{
			Username: claims.Subject,
			Role:     claims.Role,
		}, nil
	}
	return nil, nil
}
