package main

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

var END_OF_SEN = ".!?"
var TOKEN_REGEX, _ = regexp.Compile("([0-9a-zA-Zа-яА-Яё@.,!?:_/\\-']+)")

func tokenize(line string, strip bool) []string {
	words := TOKEN_REGEX.FindAllString(strings.ToLower(line), -1)
	if strip {
		for i, word := range words {
			words[i] = strings.TrimSuffix(word, END_OF_SEN)
		}
	}
	return words
}

func updateStats(stats *Stats, word string, prev string, prevprev string, sos bool) {
	if sos {
		stats.Start.addOne(word)
		return
	}
	if prev != "" {
		stats.Bi.get(prev).addOne(word)
		if prevprev != "" {
			stats.Tri.get(SPair{prevprev, prev}).addOne(word)
		}
	}

}
func ParseLine(line string) Stats {
	var stats = NewStats()
	var prev, prevprev string
	var sos = true
	for _, token := range tokenize(line, false) {
		updateStats(&stats, token, prev, prevprev, sos)
		sos = strings.ContainsAny(token[len(token)-1:], END_OF_SEN)
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

func ParseFile(path string) (Stats, error) {
	file, err := os.Open(path)
	if err != nil {
		return Stats{}, err
	}
	defer file.Close()

	var stats = NewStats()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		stats.update(ParseLine(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		return stats, err
	}
	return stats, nil
}

func ParseFiles(files []string, output string) error {
	//sp := SPair{"a", "b"}
	//kek := map[SPair]int{sp: 1}
	//fmt.Println(json.Marshal(kek))
	//return nil
	var stats = NewStats()
	for _, f := range files {
		fmt.Println("parsing", f)
		fileStats, err := ParseFile(f)
		if err != nil {
			return err
		}
		stats.update(fileStats)
	}

	statsJson, _ := json.Marshal(stats)
	err := ioutil.WriteFile(output, statsJson, 0644)
	return err
}
