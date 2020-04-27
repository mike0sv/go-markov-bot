package stats

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"strings"
)

type SPair struct {
	first  string
	second string
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

func (h Unigrams) setDefault(k string, v int) (set bool, r int) {
	if r, set = h[k]; !set {
		h[k] = v
		r = v
		set = true
	}
	return
}

func (u Unigrams) update(other Unigrams) {
	for k, v := range other {
		u.setDefault(k, v)
	}
}

func (u Unigrams) AddOne(k string) {
	v, ok := u[k]
	if !ok {
		v = 0
	}
	u.setDefault(k, v+1)
}

func (u Unigrams) Choice() (string, error) {
	var totalWeight int
	for _, v := range u {
		totalWeight += v
	}

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
		u.Get(b).update(s)
	}
}

func (u Bigrams) Get(k string) *Unigrams {
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
		u.Get(b).update(s)
	}
}

func (u Trigrams) Get(k SPair) *Unigrams {
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

func NewStats() Stats {
	return Stats{Unigrams{}, Bigrams{}, Trigrams{}}
}

func (s *Stats) UpdateAll(other Stats) {
	s.Start.update(other.Start)
	s.Bi.update(other.Bi)
	s.Tri.update(other.Tri)
}

func (s *Stats) UpdateStart(word string) {
	s.Start.AddOne(word)
}

func (s *Stats) UpdateFollowing(word, prev, prevprev string) {
	if prev != "" {
		s.Bi.Get(prev).AddOne(word)
		if prevprev != "" {
			s.Tri.Get(NewSPair(prevprev, prev)).AddOne(word)
		}
	}
}

func (s Stats) DumpToFile(path string) error {
	statsJson, _ := json.Marshal(s)
	return ioutil.WriteFile(path, statsJson, 0644)
}
