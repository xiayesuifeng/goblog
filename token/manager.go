package token

import (
	"time"
	"sync"
)

var tokenManager *TokenManager

type TokenManager struct {
	tokens []Token
}

func GetManager() *TokenManager {
	sync.Once{}.Do(func() {
		tokenManager=&TokenManager{}
	})
	return tokenManager
}

func (t *TokenManager)GetToken() Token{
	token:= NewToken(1440 * time.Second)
	t.tokens = append(t.tokens,token)
	timer :=time.NewTimer(1440 * time.Second)
	go func() {
		<- timer.C
		t.DelToken(token)
	}()
	return token
}

func (t *TokenManager)DelToken(token Token){
	tokens := make([]Token,len(t.tokens)-1)
	for i := 0; i < len(tokens); i++ {
		if  t.tokens[i].token!=token.token{
			tokens=append(tokens,token)
		}
	}
}