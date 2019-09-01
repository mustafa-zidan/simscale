package main

import (
	"github.com/mustafa-zidan/simscale/cache"
	"github.com/mustafa-zidan/simscale/parser"
	"github.com/mustafa-zidan/simscale/stats"
	"github.com/urfave/cli"
	"log"
	"os"
	"sync"
)

var inFile, outFile, version string
var flags = []cli.Flag{
	cli.StringFlag{
		Name:        "in-file",
		Value:       "",
		Usage:       "file to be processed",
		Destination: &inFile,
		Required:    true,
	},
	cli.StringFlag{
		Name:        "out-file",
		Value:       "resources/trace.txt",
		Usage:       "file to be processed",
		Destination: &outFile,
	},
}

func main() {
	app := cli.NewApp()
	app.Flags = flags
	app.Version = version

	app.Action = func(c *cli.Context) error {
		wg := &sync.WaitGroup{}
		stats.CounterInitialize()
		cache := cache.NewCache(outFile, wg)
		p, err := parser.NewParser(inFile, cache)
		p.Process()
		wg.Wait()
		log.Println(stats.CounterList())
		return err
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
