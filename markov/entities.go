package main

import (
	"errors"
	"math/rand"
	"strings"
	"time"
)

type SPair struct {
	first  string
	second string
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

type Updatable interface {
	update(other *Updatable)
}

type Unigrams map[string]int

func (h Unigrams) setDefault(k string, v int) (set bool, r int) {
	if r, set = h[k]; !set {
		h[k] = v
		r = v
		set = true
	}
	return
}

//type Unigrams2 struct {
//	stat Counter
//}

func (u Unigrams) update(other Unigrams) {
	for k, v := range other {
		u.setDefault(k, v)
		//u.stat[k] += v
	}
}

func (u Unigrams) addOne(k string) {
	v, ok := u[k]
	if !ok {
		v = 0
	}
	u.setDefault(k, v+1)
}

func (u Unigrams) choice() (string, error) {
	var totalWeight int
	for _, v := range u {
		totalWeight += v
	}

	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(totalWeight)
	for word, weight := range u {
		r -= weight
		if r <= 0 {
			return word, nil
		}
	}
	return "", errors.New("No game selected")
}

type Bigrams map[string]Unigrams

func (u Bigrams) update(other Bigrams) {
	for b, s := range other {
		u.get(b).update(s)

	}
}

func (u Bigrams) get(k string) *Unigrams {

	r, ok := u[k]
	if ok {
		return &r
	} else {
		u[k] = Unigrams{}
		r = u[k]
		return &r
	}
}

type Trigrams map[SPair]Unigrams

func (u Trigrams) update(other Trigrams) {
	for b, s := range other {
		u.get(b).update(s)
	}
}

func (u Trigrams) get(k SPair) *Unigrams {
	r, ok := u[k]
	if ok {
		return &r
	} else {
		u[k] = Unigrams{}
		r = u[k]
		return &r
	}
}

type Stats struct {
	Start Unigrams
	Bi    Bigrams
	Tri   Trigrams
}

func (s Stats) update(other Stats) {
	s.Start.update(other.Start)
	s.Bi.update(other.Bi)
	s.Tri.update(other.Tri)
}

func NewStats() Stats {
	return Stats{Unigrams{}, Bigrams{}, Trigrams{}}
}
