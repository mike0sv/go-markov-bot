package handlers

import (
	"encoding/json"
	"github.com/mike0sv/go-markov-bot/stats"
	"github.com/mike0sv/go-markov-bot/word"
	"io/ioutil"
	"net/http"
	"strings"
)

type RaverRequest struct {
	query string `json:"q,omitempty"`
}

func CreateStats(stats *stats.Stats, wg *word.Generator) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var t RaverRequest
		err = json.Unmarshal(body, &t)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		context := strings.Split(t.query, " ") // TODO context preprocessing
		data := map[string]string{"a": wg.GenerateOne(stats, context)}
		payload, _ := json.Marshal(data)
		w.Write(payload)
	}
}
