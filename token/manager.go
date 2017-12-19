package token

import (
	"time"
	"sync"
)

var (
	tokenManager *TokenManager
	once         sync.Once
)

type TokenManager struct {
	tokens []Token
}

func GetManager() *TokenManager {
	once.Do(func() {
		tokenManager = &TokenManager{}
	})
	return tokenManager
}

func (t *TokenManager) GetToken() Token {
	token := NewToken(1440 * time.Second)
	t.tokens = append(t.tokens, token)
	timer :=time.NewTimer(1440 * time.Second)
	go func() {
		<- timer.C
		t.DelToken(token)
	}()
	return token
}

func (t *TokenManager) DelToken(token Token) {
	tokens := make([]Token, len(t.tokens)-1)
	for i := 0; i < len(tokens); i++ {
		if t.tokens[i].Token != token.Token {
			tokens = append(tokens, token)
		}
	}
	t.tokens=tokens
}

func (t *TokenManager) IsExist(token string) bool {
	for i := 0; i < len(t.tokens); i++ {
		if t.tokens[i].Token == token {
			return true
		}
	}
	return false
}
