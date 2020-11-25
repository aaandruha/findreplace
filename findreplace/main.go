package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Find and replace string in file(s) or stdin"
	app.Usage = "Simple tool for find and replace"

	myFlags := []cli.Flag{
		&cli.StringFlag{
			Name:  "str",
			Value: "The quick brown fox jumps over the lazy dog",
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:  "find",
			Usage: "find str in stream",
			Flags: myFlags,
			Action: func(c *cli.Context) error {
				fmt.Println(c.String("str"))
				return nil
			},
		},
		{
			Name:  "replace",
			Usage: "find str in stream",
			Flags: myFlags,
			Action: func(c *cli.Context) error {
				fmt.Println(c.String("str"))
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
