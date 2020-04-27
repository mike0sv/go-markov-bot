package main

import (
	"fmt"
	"github.com/mike0sv/go-markov-bot/handlers"
	"github.com/mike0sv/go-markov-bot/stats"
	"github.com/mike0sv/go-markov-bot/word"
	"github.com/urfave/cli"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

const defaultWordToGenerateCount = 10

func main() {
	rand.Seed(time.Now().UnixNano())
	app := cli.NewApp()
	app.Name = "markov"
	app.Commands = []cli.Command{
		{
			Name:   "parse",
			Action: Parse,
			Usage:  "parse [files] output",
		},
		{
			Name:   "generate",
			Action: Generate,
			Usage:  "generate file count context...",
		},
		{
			Name:   "run",
			Action: Run,
			Usage:  "run file (port?)",
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func Parse(ctx *cli.Context) error {
	args := ctx.Args()
	output := args[len(args)-1]
	files := args[:len(args)-1]
	statsInstance := stats.CreateFromFiles(files...)
	return statsInstance.DumpToFile(output)
}

func Generate(ctx *cli.Context) error {
	filename := ctx.Args()[0]
	count, err := strconv.Atoi(ctx.Args()[1])
	if err != nil {
		log.Panic(err)
		return err
	}
	wordGenerator := word.NewGenerator(defaultWordToGenerateCount)
	statsInstance, err := stats.LoadFromFile(filename)
	if err != nil {
		log.Panic(err)
	}

	words, err := wordGenerator.GenerateN(&statsInstance, ctx.Args()[2:], count)
	if err != nil {
		log.Panic(err)
	}

	for _, w := range words {
		fmt.Println(w)
	}
	return nil
}

func Run(ctx *cli.Context) error {
	filename := ctx.Args()[0]
	statsInstance, err := stats.LoadFromFile(filename)
	if err != nil {
		log.Panic(err)
		return err
	}
	wordGenerator := word.NewGenerator(defaultWordToGenerateCount)
	http.HandleFunc("/", handlers.CreateStats(&statsInstance, wordGenerator)) // установим роутер
	err = http.ListenAndServe("0.0.0.0:9000", nil)                            // задаем слушать порт
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
		return err
	}
	return nil
}
