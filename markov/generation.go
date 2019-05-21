package main

import (
	"encoding/json"
	"fmt"
	"github.com/influxdata/platform/kit/errors"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"time"
)

func loadStats(filename string) (Stats, error) {
	payload, err := ioutil.ReadFile(filename)
	if err != nil {
		return Stats{}, err
	}
	stats := Stats{}
	err = json.Unmarshal(payload, &stats)
	return stats, err
}

func generateFirstAndSecond(stats *Stats, context []string) (string, string, error) {
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

/*
   def _context_choice(self, fcount, context, mult=1.):
       if context is not None:
           random.shuffle(context)
           for word in context:
               if word in fcount and random.uniform(0, 1) < self.context_probability * mult:
                   if self.context_expiring:
                       context.remove(word)
                   return word
       return fcount.choice()
*/

func makeContextChoice(u Unigrams, context []string, mult float64) (string, error) {
	if len(context) > 0 {
		shuffle(&context)
		for _, word := range context {
			if _, ok := u[word]; ok { // TODO context probability and expiration
				return word, nil
			}
		}
	}
	return u.choice()
}

func shuffle(list *[]string) {
	lst := *list
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(lst), func(i, j int) { lst[i], lst[j] = lst[j], lst[i] })
}

/*
   def gen_sentence(self, context=None, as_string=True, apx_length=20.):
       if context is None:
           context = []
       f, s = self._gen_first_second(context)
       sentence = [f, s]
       z = None
       for x in range(150):
           t = self._gen_word(f, s, context)
           if t is None:
               s = self._gen_word(z, f, context)
               continue
           sentence.append(t)
           if (f, s, t) in self.last_pair and random.uniform(0, 1.) < x / apx_length:
               # sentence.append('|')
               break
           z, f, s = f, s, t

       if as_string:
           return ' '.join(sentence)
       else:
           return sentence
*/

func GenerateOne(stats *Stats, context []string) string {
	first, second, err := generateFirstAndSecond(stats, context)
	if err != nil {
		log.Panic(err)
		return ""
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
	return strings.Join(sen, " ")
}

/*
   def _gen_word(self, first, second, context):
       pair = (first, second)
       if pair in self.trigrams:
           return self._context_choice(self.trigrams[(first, second)], context)
       if self.use_bigrams and second in self.bigrams:
           return self._context_choice(self.bigrams[second], context)
*/
func genWord(stats *Stats, first string, second string, context []string) (string, error) {
	pair := SPair{first, second}
	if val, ok := stats.Tri[pair]; ok {
		return makeContextChoice(val, context, 1)
	}
	if val, ok := stats.Bi[second]; ok {
		return makeContextChoice(val, context, 1)
	}
	return "", errors.New("no word generated")
}

func GenerateFromStats(stats *Stats, context []string, count int) {
	for i := 0; i < count; i++ {
		fmt.Println(GenerateOne(stats, context))
	}

}

func GenerateFromFile(filename string, context []string, count int) {
	if count == 0 {
		count = 10
	}
	stats, err := loadStats(filename)
	if err != nil {
		log.Panic(err)
		return
	}
	GenerateFromStats(&stats, context, count)
}
