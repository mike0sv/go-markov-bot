package markov

import (
	"github.com/influxdata/platform/kit/errors"
	"github.com/mike0sv/go-markov-bot/sliceutils"
	"log"
	"strings"
)

var (
	errCantStart      = errors.New("cant start((((")
	errCantMakeChoice = errors.New("Unable to make choice")
)

type TextGenerator struct {
	defaultSentenceToGenerateCount int
}

func NewTextGenerator(defaultSentenceToGenerateCount int) *TextGenerator {
	return &TextGenerator{defaultSentenceToGenerateCount: defaultSentenceToGenerateCount}
}

func (g TextGenerator) Generate(stats *WordStats, context []string) (string, error) {
	first, second, err := generateFirstAndSecond(stats, context)
	if err != nil {
		return "", err
	}
	sen := make([]string, 2)
	sen[0] = first
	sen[1] = second
	var zed, third string
	for x := 0; x < 150; x++ {
		third, err = generateWord(stats, first, second, context)
		if err != nil {
			continue
		}
		if third == "" {
			second, err = generateWord(stats, zed, first, context)
			continue
		}
		sen = append(sen, third)
		// TODO endings
		zed, first, second = first, second, third
	}
	return strings.Join(sen, " "), nil
}

func (g TextGenerator) GenerateN(stats *WordStats, context []string, n int) (result []string, err error) {
	if n == 0 {
		n = g.defaultSentenceToGenerateCount
	}
	for i := 0; i < n; i++ {
		var sentence string
		sentence, err = g.Generate(stats, context)
		if err != nil {
			return
		}
		result = append(result, sentence)
	}
	return
}

func generateFirstAndSecond(stats *WordStats, context []string) (first, second string, err error) {
	//try 10 times
	for i := 0; i < 10; i++ {
		first, err := makeContextChoice(stats.Un, context, .5)
		if err != nil {
			log.Println("Error making context choice", err)
			return "", "", errCantMakeChoice
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
		return "", "", errCantStart
	}
	return first, second, nil
}

func makeContextChoice(u Unigrams, context []string, mult float64) (string, error) {
	if len(context) > 0 {
		sliceutils.Shuffle(&context)
		for _, word := range context {
			if _, ok := u[word]; ok { // TODO context probability and expiration
				return word, nil
			}
		}
	}
	return u.Choice()
}

func generateWord(stats *WordStats, first string, second string, context []string) (string, error) {
	pair := NewSPair(first, second)
	if val, ok := stats.Tri[pair]; ok {
		return makeContextChoice(val, context, 1)
	}
	if val, ok := stats.Bi[second]; ok {
		return makeContextChoice(val, context, 1)
	}
	return "", errors.New("no word generated")
}
