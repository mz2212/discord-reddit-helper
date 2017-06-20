// Based heavily upon "markov.go" found at: https://golang.org/doc/codewalk/markov.go

package markov

import (
	"math/rand"
	"strings"
	"time"
)

type Prefix []string

func (pre Prefix) String() string {
	return strings.Join(pre, " ")
}

func (pre Prefix) Shift(word string) {
	copy(pre, pre[1:])
	pre[len(pre)-1] = word
}

type Chain struct {
	chain     map[string][]string
	prefixLen int
}

func New(preLen int) *Chain {
	return &Chain{
		make(map[string][]string),
		preLen,
	}
}

func (cha *Chain) Build(text string) {
	words := strings.Split(text, " ")
	pre := make(Prefix, cha.prefixLen)
	for _, w := range words {
		key := pre.String()
		cha.chain[key] = append(cha.chain[key], w)
		pre.Shift(w)
	}
}

// The function that generates the sentence. Basically Magic.
// But not really.
func (cha *Chain) Generate(n int) string {
	rand.Seed(time.Now().Unix())
	pre := make(Prefix, cha.prefixLen)
	var words []string
	for i := 0; i < n; i++ {
		choices := cha.chain[pre.String()]
		if len(choices) == 0 {
			break
		}
		next := choices[rand.Intn(len(choices))]
		words = append(words, next)
		pre.Shift(next)
	}
	return strings.Join(words, " ")
}
