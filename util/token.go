package util

import (
	"context"
)

const (
	headerAuthorize string = "authorization"
	tokenPrefix     string = "bearer "
)

type Token struct {
	Value string
}

func (t *Token) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{headerAuthorize: t.Value}, nil
}

func (t *Token) RequireTransportSecurity() bool {
	return true
}

func GetToken(bearer string) *Token {
	return &Token{
		Value: tokenPrefix + bearer,
	}
}
