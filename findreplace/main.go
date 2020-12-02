package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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
				findLinesInFile(os.Stdin, os.Args[2])
			} else {
				walkDir(file)
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
		fmt.Fprintf(os.Stderr, "OpenError: %v\n", err)
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

func walkDir(path string) {
	// todo: добавить возможность подстановки последнего символа '/' в путь файла

	fi, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		entries, err := ioutil.ReadDir(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "wlkError: %v\n", err)
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				findLines(path+entry.Name(), os.Args[2])
			}
		}
	case mode.IsRegular():
		findLines(path, os.Args[2])
	}

}

func findLinesInFile(f *os.File, str string) {
	input := bufio.NewScanner(f)
	i := 1
	for input.Scan() {
		s := input.Text()
		if strings.Index(s, str) >= 0 {
			fmt.Printf("%s:%d - %s\n", filepath.Base(os.Args[3]), i, input.Text())
		}
		i++
	}

}
