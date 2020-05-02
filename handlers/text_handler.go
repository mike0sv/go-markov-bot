package handlers

import (
	"encoding/json"
	"github.com/mike0sv/go-markov-bot/markov"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type RaverRequest struct {
	query string `json:"q,omitempty"`
}

func CreateText(stats *markov.WordStats, tg *markov.TextGenerator) func(w http.ResponseWriter, r *http.Request) {
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
		text, err := tg.Generate(stats, context)
		if err != nil {
			log.Println("Error generating text", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		data := map[string]string{"a": text}
		payload, _ := json.Marshal(data)
		w.Write(payload)
	}
}
