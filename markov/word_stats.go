package markov

import (
	"encoding/json"
	"io/ioutil"
)

type WordStats struct {
	Un  Unigrams
	Bi  Bigrams
	Tri Trigrams
}

func NewStats() WordStats {
	return WordStats{Unigrams{}, Bigrams{}, Trigrams{}}
}

func (ws *WordStats) Add(word string) {
	ws.Un.add(word)
}

func (ws *WordStats) Merge(other WordStats) {
	ws.Un.update(other.Un)
	ws.Bi.update(other.Bi)
	ws.Tri.update(other.Tri)
}

func (ws *WordStats) Update(word, prev, prevprev string) {
	if prev != "" {
		ws.Bi.get(prev).add(word)
		if prevprev != "" {
			ws.Tri.get(NewSPair(prevprev, prev)).add(word)
		}
	}
}

func (ws WordStats) DumpToFile(path string) error {
	statsJson, _ := json.Marshal(ws)
	return ioutil.WriteFile(path, statsJson, 0644)
}
