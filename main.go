package main

import (
	"github.com/urfave/cli"

	"log"
	"os"
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
		Value:       "out.json",
		Usage:       "file to be processed",
		Destination: &outFile,
	},
}

func main() {
	app := cli.NewApp()
	app.Flags = flags
	app.Version = version

	app.Action = func(c *cli.Context) error {
		CounterInitialize()
		InitTTLCache()
		p, err := NewParser(inFile)
		p.Process()
		return err
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
