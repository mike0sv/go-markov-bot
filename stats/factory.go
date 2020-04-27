package stats

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

var (
	endOfSen      = ".!?"
	tokenRegex, _ = regexp.Compile("([0-9a-zA-Zа-яА-Яё@.,!?:_/\\-']+)")
)

func CreateFromFiles(paths ...string) (stats Stats) {
	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			log.Println("Error while opening file:", path, err)
			continue
		}

		sc := bufio.NewScanner(file)
		for sc.Scan() {
			stats.UpdateAll(parseLine(sc.Text()))
		}

		if err := sc.Err(); err != nil {
			log.Println("Error while reading file:", path, err)
		}
		file.Close()
	}
	return
}

func LoadFromFile(filename string) (stats Stats, err error) {
	payload, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	err = json.Unmarshal(payload, &stats)
	return
}

func parseLine(line string) Stats {
	var stats = NewStats()
	var prev, prevprev string
	var sos = true
	for _, token := range tokenize(line, false) {
		if sos {
			stats.UpdateStart(token)
		} else {
			stats.UpdateFollowing(token, prev, prevprev)
		}
		sos = strings.ContainsAny(token[len(token)-1:], endOfSen)
		prev, prevprev = token, prev
	}
	s, e := json.Marshal(stats)
	fmt.Println(e)
	if e != nil {
		log.Panic(e)
	}
	fmt.Println(string(s))
	return stats
}

func tokenize(line string, strip bool) []string {
	words := tokenRegex.FindAllString(strings.ToLower(line), -1)
	if strip {
		for i, word := range words {
			words[i] = strings.TrimSuffix(word, endOfSen)
		}
	}
	return words
}
