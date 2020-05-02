package markov

import (
	"errors"
	"math/rand"
	"strings"
)

type SPair struct {
	first, second string
}

func NewSPair(first, second string) SPair {
	return SPair{first: first, second: second}
}

func (sp SPair) MarshalText() (text []byte, err error) {
	return []byte(sp.first + " " + sp.second), nil
}

func (sp SPair) UnmarshalText(text []byte) error {
	words := strings.Split(string(text), " ")
	if len(words) != 2 {
		return errors.New("wrong pair key " + string(text))
	}
	sp.first = words[0]
	sp.second = words[1]
	return nil
}

type Unigrams map[string]int

func (un Unigrams) add(word string) {
	if _, ok := un[word]; !ok {
		un[word] = 1
	}
}

func (un Unigrams) update(other Unigrams) {
	//set if not exist
	for k, v := range other {
		if _, ok := un[k]; !ok {
			un[k] = v
		}
	}
}

func (un Unigrams) Choice() (string, error) {
	var totalWeight int
	for _, v := range un {
		totalWeight += v
	}

	r := rand.Intn(totalWeight)
	for word, weight := range un {
		r -= weight
		if r <= 0 {
			return word, nil
		}
	}
	return "", errors.New("No game selected")
}

type Bigrams map[string]Unigrams

func (bi Bigrams) update(other Bigrams) {
	for b, s := range other {
		bi.get(b).update(s)
	}
}

func (bi Bigrams) get(k string) *Unigrams {
	if r, ok := bi[k]; ok {
		return &r
	}
	//or create empty
	bi[k] = Unigrams{}
	r := bi[k]
	return &r
}

type Trigrams map[SPair]Unigrams

func (tri Trigrams) get(k SPair) *Unigrams {
	if r, ok := tri[k]; ok {
		return &r
	}
	//or create empty
	tri[k] = Unigrams{}
	r := tri[k]
	return &r
}

func (tri Trigrams) update(other Trigrams) {
	for b, s := range other {
		tri.get(b).update(s)
	}
}
