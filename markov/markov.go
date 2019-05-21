package main

import (
	"github.com/urfave/cli"
	"log"
	"os"
	"strconv"
)

func Parse(ctx *cli.Context) error {
	args := ctx.Args()
	output := args[len(args)-1]
	files := args[:len(args)-1]
	return ParseFiles(files, output)
}

func Generate(ctx *cli.Context) error {
	filename := ctx.Args()[0]
	count, err := strconv.Atoi(ctx.Args()[1])
	if err != nil {
		log.Panic(err)
		return err
	}
	GenerateFromFile(filename, ctx.Args()[2:], count)
	return nil
}

func Run(ctx *cli.Context) error {
	filename := ctx.Args()[0]
	stats, err := loadStats(filename)
	if err != nil {
		log.Panic(err)
		return err
	}
	return RunServer(&stats)
}
func main() {
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
