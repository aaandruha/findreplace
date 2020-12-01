package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "findreplace"
	app.Usage = "Find and replace string in file(s) or stdin"
	app.Description = "Description"
	app.Authors = []*cli.Author{
		{Name: "Andrew Gly", Email: "glybin.av@gmail.com"},
	}

	app.Commands = []*cli.Command{
		findCommand(),
		replaceCommand(),
	}

	app.Action = mainAction

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func findCommand() *cli.Command {
	return &cli.Command{
		Name:    "find",
		Aliases: []string{"f"},
		Usage:   "find str in stream",
		Action: func(c *cli.Context) error {
			n := c.NArg()
			if n == 0 {
				return fmt.Errorf("No arguments for find command")
			}
			if n == 1 {
				return fmt.Errorf("No find string in arguments")
			}

			file := os.Args[3]
			if len(file) == 0 {
				//findLines(os.Stdin, os.Args[2])
			} else {
				//	for _, arg := range file {
				/*f, err := os.Open(file)
				if err != nil {
					log.Fatal(err)
				}
				findLines(f, os.Args[2])
				f.Close()*/
				walkDir(file)
				//	}
			}
			return nil
		},
	}
}

func replaceCommand() *cli.Command {
	return &cli.Command{
		Name:    "replace",
		Aliases: []string{"r"},
		Usage:   "replace str in stream",
		Action: func(c *cli.Context) error {
			n := c.NArg()
			if n == 0 {
				return fmt.Errorf("No arguments for replace command")
			}

			return nil
		},
	}
}

func mainAction(arg *cli.Context) error {
	arg.App.Command("help").Run(arg)
	return nil
}

func findLines(fileName, str string) {
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dul: %v\n", err)
	}
	input := bufio.NewScanner(f)
	i := 1
	for input.Scan() {
		//counts[input.Text()]++
		s := input.Text()
		if strings.Index(s, str) >= 0 {
			fmt.Printf("%s:%d - %s\n", fileName, i, input.Text())
		}
		i++
	}
	f.Close()

}

func walkDir(dir string) {
	// todo: добавить возможность подстановки последнего символа '/' в путь файла
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dul: %v\n", err)
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			findLines(dir+entry.Name(), os.Args[2])
		}
	}
}
