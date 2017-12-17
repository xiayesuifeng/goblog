package token

import (
	"crypto/sha1"
	"strconv"
	"fmt"
	"math/rand"
	"time"
)

var r *rand.Rand

type Token struct {
	Token string
	Ttl int64
}

func NewToken(ttl time.Duration) Token{
	if r==nil {
		r=rand.New(rand.NewSource(time.Now().Unix()))
	}
	i := r.Int63n(10000)
	h := sha1.New()
	h.Write([]byte(strconv.FormatInt(i, 10)))
	fmt.Println(i)
	token := fmt.Sprintf("%x", h.Sum([]byte("goblog")))
	return Token{Token:token,Ttl:int64(ttl.Seconds())}
}