package markov

import (
	"bufio"
	"encoding/json"
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

func CreateStatsFromFiles(paths ...string) WordStats {
	stats := NewStats()
	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			log.Println("Error while opening file:", path, err)
			continue
		}

		sc := bufio.NewScanner(file)
		for sc.Scan() {
			stats.Merge(createStatsFromLine(sc.Text()))
		}

		if err := sc.Err(); err != nil {
			log.Println("Error while reading file:", path, err)
		}
		file.Close()
	}
	return stats
}

func LoadStatsFromFile(filename string) (stats WordStats, err error) {
	payload, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	err = json.Unmarshal(payload, &stats)
	return
}

func createStatsFromLine(line string) WordStats {
	stats := NewStats()
	var prev, prevprev string
	sos := true
	for _, token := range tokenRegex.FindAllString(strings.ToLower(line), -1) {
		if sos {
			stats.Add(token)
		} else {
			stats.Update(token, prev, prevprev)
		}
		sos = strings.ContainsAny(token[len(token)-1:], endOfSen)
		prev, prevprev = token, prev
	}
	return stats
}
