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
		Name:        "in-file, i",
		Value:       "",
		Usage:       "file to be processed",
		Destination: &inFile,
		Required:    true,
	},
	cli.StringFlag{
		Name:        "out-file, o",
		Value:       "trace.txt",
		Usage:       "output trace file ",
		Destination: &outFile,
	},
}

func main() {
	app := cli.NewApp()
	app.Name = "Simscale"
	app.Usage = "Simscale coding challange"
	app.Flags = flags
	app.Version = version

	app.Action = func(c *cli.Context) error {
		wg := &sync.WaitGroup{}
		stats.CounterInitialize()
		cache := cache.NewCache(outFile, wg)
		p, err := parser.NewParser(inFile, cache)
		p.Process()
		wg.Wait()
		s := stats.CounterList()
		log.Printf("Total Number of Logs: \t\t %d\n", s["total"])
		log.Printf("Number of Records parsed: \t\t %d\n", s["success"])
		log.Printf("Number of Traces: \t\t\t %d\n", s["traces"])
		log.Printf("Number of Orphen Logs: \t\t %d\n", s["orphens"])
		log.Printf("Number of Malformed Logs: \t\t %d\n", s["malformed"])
		return err
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
