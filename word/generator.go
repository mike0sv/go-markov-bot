package word

import (
	"github.com/influxdata/platform/kit/errors"
	"github.com/mike0sv/go-markov-bot/stats"
	"log"
	"math/rand"
	"strings"
)

type Generator struct {
	defaultWordToGenerateCount int
}

func NewGenerator(defaultWordToGenerateCount int) *Generator {
	return &Generator{defaultWordToGenerateCount: defaultWordToGenerateCount}
}

func (g Generator) GenerateOne(stats *stats.Stats, context []string) string {
	words, err := g.GenerateN(stats, context, 1)
	if err != nil || len(words) == 0 {
		log.Println("Error generating word", err)
		return ""
	}
	return words[0]
}

func (g Generator) GenerateN(stats *stats.Stats, context []string, n int) (result []string, err error) {
	if n == 0 {
		n = g.defaultWordToGenerateCount
	}
	for i := 0; i < n; i++ {
		first, second, err := generateFirstAndSecond(stats, context)
		if err != nil {
			log.Panic(err)
			return
		}
		sen := make([]string, 2)
		sen[0] = first
		sen[1] = second
		var zed, third string
		for x := 0; x < 150; x++ {
			third, err = genWord(stats, first, second, context)
			if err != nil {
				continue
			}
			if third == "" {
				second, err = genWord(stats, zed, first, context)
				continue
			}
			sen = append(sen, third)
			// TODO endings
			zed, first, second = first, second, third
		}
		result = append(result, strings.Join(sen, " "))
	}
	return
}

func generateFirstAndSecond(stats *stats.Stats, context []string) (string, string, error) {
	var first, second string
	for i := 0; i < 10; i++ {
		first, err := makeContextChoice(stats.Start, context, .5)
		if err != nil {
			return "", "", err
		}
		if _, ok := stats.Bi[first]; !ok {
			continue
		}
		second, err = makeContextChoice(stats.Bi[first], context, 1)
		if err != nil {
			return "", "", err
		}

		break
	}
	if second == "" {
		return "", "", errors.New("cant start((((")
	}
	return first, second, nil
}

func makeContextChoice(u stats.Unigrams, context []string, mult float64) (string, error) {
	if len(context) > 0 {
		shuffle(&context)
		for _, word := range context {
			if _, ok := u[word]; ok { // TODO context probability and expiration
				return word, nil
			}
		}
	}
	return u.Choice()
}

func shuffle(list *[]string) {
	lst := *list
	rand.Shuffle(len(lst), func(i, j int) { lst[i], lst[j] = lst[j], lst[i] })
}

func genWord(statsInstance *stats.Stats, first string, second string, context []string) (string, error) {
	pair := stats.NewSPair(first, second)
	if val, ok := statsInstance.Tri[pair]; ok {
		return makeContextChoice(val, context, 1)
	}
	if val, ok := statsInstance.Bi[second]; ok {
		return makeContextChoice(val, context, 1)
	}
	return "", errors.New("no word generated")
}
