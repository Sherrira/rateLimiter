package auth

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Token struct {
	Token string
	Limit int
}

type Authorizer interface {
	Authorize(tokenString string) *Token
}

type MockedAuthorizer struct {
	tokenLimits map[string]int
}

func NewAuthorizer() Authorizer {
	tokenLimitsEnv := os.Getenv("TOKEN_LIMITS")
	tokenLimits, err := parseTokenLimits(tokenLimitsEnv)
	if err != nil {
		log.Fatalf("Error parsing token limits: %v", err)
	}
	return &MockedAuthorizer{
		tokenLimits: tokenLimits,
	}
}

func (a *MockedAuthorizer) Authorize(tokenString string) *Token {
	if limit, exists := a.tokenLimits[tokenString]; exists {
		return &Token{
			Token: tokenString,
			Limit: limit,
		}
	}
	return nil
}

func parseTokenLimits(envVar string) (map[string]int, error) {
	tokenLimits := make(map[string]int)
	pairs := strings.Split(envVar, ",")
	for _, pair := range pairs {
		kv := strings.Split(pair, ":")
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid token limit pair: %s", pair)
		}
		key := kv[0]
		value, err := strconv.Atoi(kv[1])
		if err != nil {
			return nil, fmt.Errorf("invalid token limit value for key %s: %s", key, kv[1])
		}
		tokenLimits[key] = value
	}
	return tokenLimits, nil
}
